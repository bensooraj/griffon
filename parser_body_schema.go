package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/topo"
)

func ParseHCLUsingBodySchema(filename string, src []byte, ctx *hcl.EvalContext) (*Config, error) {
	config := Config{
		SSHKeys:        make(map[string]SSHKeyBlock),
		StartupScripts: make(map[string]StartupScriptBlock),
		Instances:      make(map[string]InstanceBlock),
		DataBlocks:     make(map[string]map[string]Block),
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	bodyContent, diags := file.Body.Content(ConfigSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, errors.New("no blocks found")
	}

	blockPathToGraphID := make(map[string]int64)

	dependencyGraph := NewGraph()
	var dGraphNodeCount int64 = 0

	blocks := bodyContent.Blocks.ByType()
	for blockName, hclBlocks := range blocks {
		fmt.Println("blockName:", blockName)
		switch blockName {
		case "griffon":
			if len(hclBlocks) != 1 {
				return nil, errors.New("only one griffon block allowed")
			}
			var griffon GriffonBlock
			if err := griffon.PreProcessHCLBlock(hclBlocks[0], ctx); err != nil {
				return nil, err
			}
			griffon.GraphID = 0
			config.Griffon = griffon
			blockPathToGraphID[blockName] = griffon.GraphID

			dependencyGraph.AddNode(&griffon)
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey SSHKeyBlock
				sshKey.Name = hclBlock.Labels[0]

				if err := sshKey.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				sshKey.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, sshKey.Name)] = sshKey.GraphID

				config.SSHKeys[sshKey.Name] = sshKey
				dependencyGraph.AddNode(&sshKey)
			}
		case "startup_script":
			for _, hclBlock := range hclBlocks {
				var startupScript StartupScriptBlock
				startupScript.Name = hclBlock.Labels[0]

				if err := startupScript.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				startupScript.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, startupScript.Name)] = startupScript.GraphID

				config.StartupScripts[startupScript.Name] = startupScript
				dependencyGraph.AddNode(&startupScript)
			}
		case "instance":
			for _, hclBlock := range hclBlocks {
				var instance InstanceBlock
				instance.Name = hclBlock.Labels[0]

				if err := instance.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				instance.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, instance.Name)] = instance.GraphID

				config.Instances[instance.Name] = instance
				dependencyGraph.AddNode(&instance)
			}
		case "data":
			config.DataBlocks["region"] = make(map[string]Block)
			config.DataBlocks["plan"] = make(map[string]Block)
			config.DataBlocks["os"] = make(map[string]Block)

			for _, hclBlock := range hclBlocks {
				dGraphNodeCount++
				dataBlock := DataBlock{
					Type:    hclBlock.Labels[0],
					Name:    hclBlock.Labels[1],
					GraphID: dGraphNodeCount,
				}
				blockPathToGraphID[BuildBlockPath(blockName, dataBlock.Type, dataBlock.Name)] = dataBlock.GraphID

				switch dataBlock.Type {
				case "region":
					regionData := RegionDataBlock{DataBlock: dataBlock}
					if err := regionData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &regionData
					dependencyGraph.AddNode(&regionData)
				case "plan":
					planData := PlanDataBlock{DataBlock: dataBlock}
					if err := planData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &planData
					dependencyGraph.AddNode(&planData)
				case "os":
					osData := OSDataBlock{DataBlock: dataBlock}
					if err := osData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &osData
					dependencyGraph.AddNode(&osData)
				default:
					return nil, errors.New("unknown data type " + dataBlock.Type)
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
		node := nodes.Node().(Block)
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

	return &config, nil
}
