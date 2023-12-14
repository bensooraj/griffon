package parser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/bensooraj/griffon/blocks"
	"github.com/bensooraj/griffon/graph"
	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/oauth2"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

func ParseWithBodySchema(filename string, src []byte, evalCtx *hcl.EvalContext) (*blocks.Config, error) {
	config := blocks.Config{
		Griffon:   blocks.GriffonBlock{},
		Data:      make(map[blocks.BlockType]map[string]blocks.Block),
		Resources: make(map[blocks.BlockType]map[string]blocks.Block),
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	bodyContent, diags := file.Body.Content(schema.ConfigSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, errors.New("no blocks found")
	}

	var dGraphNodeCount int64 = 0

	hclBlocks := bodyContent.Blocks.ByType()
	for blockName, hclBlocks := range hclBlocks {
		switch blockName {
		case "griffon":
			if len(hclBlocks) != 1 {
				return nil, errors.New("only one griffon block allowed")
			}
			var griffon blocks.GriffonBlock
			if err := griffon.PreProcessHCLBlock(hclBlocks[0], evalCtx); err != nil {
				return nil, err
			}
			griffon.GraphID = 0
			config.Griffon = griffon
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey blocks.SSHKeyBlock
				sshKey.Name = hclBlock.Labels[0]

				if err := sshKey.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				sshKey.GraphID = dGraphNodeCount

				config.AddResource(&sshKey)
			}
		case "startup_script":
			for _, hclBlock := range hclBlocks {
				var startupScript blocks.StartupScriptBlock
				startupScript.Name = hclBlock.Labels[0]

				if err := startupScript.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				startupScript.GraphID = dGraphNodeCount

				config.AddResource(&startupScript)
			}
		case "instance":
			for _, hclBlock := range hclBlocks {
				var instance blocks.InstanceBlock
				instance.Name = hclBlock.Labels[0]

				if err := instance.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				instance.GraphID = dGraphNodeCount

				config.AddResource(&instance)
			}
		case "data":
			for _, hclBlock := range hclBlocks {
				dGraphNodeCount++
				dataBlock := blocks.DataBlock{
					Type:    blocks.BlockType(hclBlock.Labels[0]),
					Name:    hclBlock.Labels[1],
					GraphID: dGraphNodeCount,
				}

				switch dataBlock.Type {
				case "region":
					regionData := blocks.RegionDataBlock{DataBlock: dataBlock}
					if err := regionData.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
						return nil, err
					}

					config.AddData(&regionData)
				case "plan":
					planData := blocks.PlanDataBlock{DataBlock: dataBlock}
					if err := planData.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
						return nil, err
					}

					config.AddData(&planData)
				case "os":
					osData := blocks.OSDataBlock{DataBlock: dataBlock}
					if err := osData.PreProcessHCLBlock(hclBlock, evalCtx); err != nil {
						return nil, err
					}

					config.AddData(&osData)
				default:
					return nil, errors.New("unknown data type " + string(dataBlock.Type))
				}
			}
		default:
			slog.Debug("[PARSER] unknown block type", slog.String("block_type", blockName))
		}
	}

	return &config, nil
}

func CalculateEvaluationOrder(config *blocks.Config) (*graph.DependencyGraph, error) {
	dependencyGraph := graph.NewDependencyGraph()

	err := dependencyGraph.LoadGriffonConfig(config)
	if err != nil {
		return nil, err
	}

	if os.Getenv("GENERATE_DOT_FILE") == "true" {
		dotByteArr, err := dot.Marshal(dependencyGraph, "dependency_graph.dot", "", "")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile("dependency_graph.dot", dotByteArr, 0644)
		if err != nil {
			return nil, err
		}
	}

	sortedNodeIDs, err := dependencyGraph.GetSortedNodeIDs()
	if err != nil {
		return nil, err
	}
	config.EvaluationOrder = sortedNodeIDs

	return dependencyGraph, nil
}

func EvaluateConfig(evalCtx *hcl.EvalContext, config *blocks.Config, vc *govultr.Client) error {
	dependencyGraph, err := CalculateEvaluationOrder(config)
	if err != nil {
		return err
	}

	for _, nodeID := range config.EvaluationOrder {
		node := dependencyGraph.Node(nodeID).(blocks.Block)

		// Process the block using the evaluation context
		err := node.ProcessConfiguration(evalCtx)
		if err != nil {
			return err
		}

		switch b := node.(type) {
		case *blocks.GriffonBlock:
			slog.Debug("[Evaluate] GriffonBlock", slog.String("block_type", string(b.BlockType())), slog.String("block_name", string(b.BlockName())), slog.Any("griffon", b))
			if b.VultrAPIKey == "" {
				return fmt.Errorf("vultr_api_key is required")
			}
			oauthConfig := &oauth2.Config{}
			tokenSource := oauthConfig.TokenSource(context.Background(), &oauth2.Token{AccessToken: b.VultrAPIKey})
			vc = govultr.NewClient(oauth2.NewClient(context.Background(), tokenSource))

			// Optional changes
			_ = vc.SetBaseURL("https://api.vultr.com")
			vc.SetUserAgent("mycool-app")
			vc.SetRateLimit(500)

		case *blocks.RegionDataBlock,
			*blocks.OSDataBlock,
			*blocks.PlanDataBlock:
			err := b.Get(context.Background(), evalCtx, vc)
			if err != nil {
				return err
			}
			err = AddBlockToEvalContext(evalCtx, b)
			if err != nil {
				return err
			}

		case *blocks.SSHKeyBlock,
			*blocks.StartupScriptBlock,
			*blocks.InstanceBlock:
			err := b.Create(context.Background(), evalCtx, vc)
			if err != nil {
				return err
			}
			err = AddBlockToEvalContext(evalCtx, b)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown block type %T", node)
		}
		slog.Debug("[Evaluate] Evaluation context", slog.String("block_type", string(node.BlockType())), slog.String("block_name", string(node.BlockName())), slog.String("eval_ctx", evalCtx.Variables[string(node.BlockType())].GoString()))
	}

	return nil
}

func AddBlockToEvalContext(evalCtx *hcl.EvalContext, block blocks.Block) error {
	data, ok := evalCtx.Variables["data"]
	if !ok {
		return fmt.Errorf("data block not found in evaluation context")
	}

	switch block.BlockType() {
	case blocks.RegionBlockType,
		blocks.OSBlockType,
		blocks.PlanBlockType:
		dataMap := data.AsValueMap()
		// Check if the blockTypeVal key exists in dataVars
		if blockCtyVal := data.GetAttr(string(block.BlockType())); !blockCtyVal.IsNull() {
			// Get the map of blocks
			var blockMap map[string]cty.Value
			// Convert the block to a cty.Value
			blockVal, err := block.ToCtyValue()
			if err != nil {
				return err
			}

			if blockCtyVal.Equals(cty.EmptyObjectVal).True() {
				blockMap = make(map[string]cty.Value)
			} else {
				blockMap = blockCtyVal.AsValueMap()
			}
			// Add the block to the region map
			blockMap[block.BlockName()] = blockVal
			// Set the data key to the updated data
			dataMap[string(block.BlockType())] = cty.ObjectVal(blockMap)
			evalCtx.Variables["data"] = cty.ObjectVal(dataMap)
		}
	case blocks.InstanceBlockType,
		blocks.SSHKeyBlockType,
		blocks.StartupScriptBlockType:
		// Get the map of blocks
		var blockMap map[string]cty.Value
		// Convert the block to a cty.Value
		blockVal, err := block.ToCtyValue()
		if err != nil {
			return err
		}

		blockCtyVal := evalCtx.Variables[string(block.BlockType())]
		if blockCtyVal.Equals(cty.EmptyObjectVal).True() {
			blockMap = make(map[string]cty.Value)
		} else {
			blockMap = blockCtyVal.AsValueMap()
		}
		// Add the block to the region map
		blockMap[block.BlockName()] = blockVal
		evalCtx.Variables[string(block.BlockType())] = cty.ObjectVal(blockMap)

	default:
		return fmt.Errorf("block type %s not supported", block.BlockType())
	}

	return nil
}
