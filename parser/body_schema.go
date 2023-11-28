package parser

import (
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
	"gonum.org/v1/gonum/graph/topo"
)

func ParseHCLUsingBodySchema(filename string, src []byte, ctx *hcl.EvalContext, vc *govultr.Client) (*blocks.Config, error) {
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

	blockPathToGraphID := make(map[string]int64)

	dependencyGraph := graph.NewGraph()
	var dGraphNodeCount int64 = 0

	hclBlocks := bodyContent.Blocks.ByType()
	for blockName, hclBlocks := range hclBlocks {
		fmt.Println("blockName:", blockName)
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
			blockPathToGraphID[blockName] = griffon.GraphID

			dependencyGraph.AddNode(&griffon)
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey blocks.SSHKeyBlock
				sshKey.Name = hclBlock.Labels[0]

				if err := sshKey.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				sshKey.GraphID = dGraphNodeCount
				blockPathToGraphID[blocks.BuildBlockPath(blockName, sshKey.Name)] = sshKey.GraphID

				config.AddResource(&sshKey)
				dependencyGraph.AddNode(&sshKey)
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
				blockPathToGraphID[blocks.BuildBlockPath(blockName, startupScript.Name)] = startupScript.GraphID

				config.AddResource(&startupScript)
				dependencyGraph.AddNode(&startupScript)
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
				blockPathToGraphID[blocks.BuildBlockPath(blockName, instance.Name)] = instance.GraphID

				config.AddResource(&instance)
				dependencyGraph.AddNode(&instance)
			}
		case "data":
			for _, hclBlock := range hclBlocks {
				dGraphNodeCount++
				dataBlock := blocks.DataBlock{
					Type:    blocks.BlockType(hclBlock.Labels[0]),
					Name:    hclBlock.Labels[1],
					GraphID: dGraphNodeCount,
				}
				blockPathToGraphID[blocks.BuildBlockPath(blockName, string(dataBlock.Type), dataBlock.Name)] = dataBlock.GraphID

				switch dataBlock.Type {
				case "region":
					regionData := blocks.RegionDataBlock{DataBlock: dataBlock}
					if err := regionData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddResource(&regionData)
					dependencyGraph.AddNode(&regionData)
				case "plan":
					planData := blocks.PlanDataBlock{DataBlock: dataBlock}
					if err := planData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddResource(&planData)
					dependencyGraph.AddNode(&planData)
				case "os":
					osData := blocks.OSDataBlock{DataBlock: dataBlock}
					if err := osData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.AddResource(&osData)
					dependencyGraph.AddNode(&osData)
				default:
					return nil, errors.New("unknown data type " + string(dataBlock.Type))
				}
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}

	fmt.Println()

	fmt.Println("dependencyGraph.Nodes().Len():", dependencyGraph.Nodes().Len())
	fmt.Println("blockPathToGraphID:", blockPathToGraphID)

	nodes := dependencyGraph.Nodes()
	for nodes.Next() {
		node := nodes.Node().(blocks.Block)
		d := node.Dependencies()
		for _, dep := range d {
			if dNode := dependencyGraph.Node(blockPathToGraphID[dep]); dNode != nil {
				dependencyGraph.SetEdge(dependencyGraph.NewEdge(dNode, node))
			}
		}
		fmt.Println("node:", node.ID(), node)
	}

	dotByteArr, err := dot.Marshal(dependencyGraph, "dependency_graph.dot", "", "")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile("graph.dot", dotByteArr, 0644)
	if err != nil {
		return nil, err
	}

	fmt.Println()
	sDependencyGraph, err := topo.Sort(dependencyGraph)
	if err != nil {
		return nil, err
	}
	_ = sDependencyGraph

	// Test API calls
	nodes = dependencyGraph.Nodes()
	for nodes.Next() {
		node := nodes.Node().(blocks.Block)
		node.Create(ctx, vc)
		// fmt.Println("node:", node.ID(), node)
	}

	return &config, nil
}
