package parser

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bensooraj/griffon/blocks"
	"github.com/bensooraj/griffon/graph"
	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/vultr/govultr/v3"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

func ParseWithBodySchema(filename string, src []byte, ctx *hcl.EvalContext, vc *govultr.Client) (*blocks.Config, error) {
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
		// fmt.Println("blockName:", blockName)
		switch blockName {
		case "griffon":
			if len(hclBlocks) != 1 {
				return nil, errors.New("only one griffon block allowed")
			}
			var griffon blocks.GriffonBlock
			if err := griffon.PreProcessHCLBlock(hclBlocks[0], ctx); err != nil {
				return nil, err
			}
			griffon.GraphID = 0
			config.Griffon = griffon
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey blocks.SSHKeyBlock
				sshKey.Name = hclBlock.Labels[0]

				if err := sshKey.PreProcessHCLBlock(hclBlock, ctx); err != nil {
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

				if err := startupScript.PreProcessHCLBlock(hclBlock, ctx); err != nil {
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

				if err := instance.PreProcessHCLBlock(hclBlock, ctx); err != nil {
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
					if err := regionData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddData(&regionData)
				case "plan":
					planData := blocks.PlanDataBlock{DataBlock: dataBlock}
					if err := planData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddData(&planData)
				case "os":
					osData := blocks.OSDataBlock{DataBlock: dataBlock}
					if err := osData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddData(&osData)
				default:
					return nil, errors.New("unknown data type " + string(dataBlock.Type))
				}
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}
	fmt.Println()

	err := EvaluateConfig(&config, vc)
	if err != nil {
		fmt.Println("Error evaluating config:", err)
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
	//
	// fmt.Println("\nEvluation order:")
	// for i, nID := range sortedNodeIDs {
	// 	b := dependencyGraph.Node(nID).(blocks.Block)
	// 	fmt.Println(i, nID, b.BlockType(), b.BlockName())
	// }
	fmt.Println()

	return dependencyGraph, nil
}

func EvaluateConfig(config *blocks.Config, vc *govultr.Client) error {
	dependencyGraph, err := CalculateEvaluationOrder(config)
	if err != nil {
		return err
	}

	evalCtx := GetEvalContext()

	for _, nodeID := range config.EvaluationOrder {
		node := dependencyGraph.Node(nodeID).(blocks.Block)

		// Process the block using the evaluation context
		err := node.ProcessConfiguration(evalCtx)
		if err != nil {
			return err
		}

		switch b := node.(type) {
		case *blocks.GriffonBlock:
			fmt.Println("GriffonBlock:", b.BlockType(), b.BlockName())
		case *blocks.RegionDataBlock,
			*blocks.OSDataBlock,
			*blocks.PlanDataBlock:
			fmt.Println("Data:", b.BlockType(), b.BlockName())
			_, err := b.Get(context.Background(), evalCtx, vc)
			if err != nil {
				return err
			}
			fmt.Println("Evaluation context:", evalCtx.Variables["data"].GoString())

		case *blocks.SSHKeyBlock,
			*blocks.StartupScriptBlock,
			*blocks.InstanceBlock:
			fmt.Println("Resource:", b.BlockType(), b.BlockName())
			_, err := b.Create(context.Background(), evalCtx, vc)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown block type %T", node)
		}
	}

	return nil
}
