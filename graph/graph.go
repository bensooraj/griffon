package graph

import (
	"fmt"

	"github.com/bensooraj/griffon/blocks"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type DependencyGraph struct {
	blockPathToGraphID map[string]int64
	*simple.DirectedGraph
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		blockPathToGraphID: make(map[string]int64),
		DirectedGraph:      simple.NewDirectedGraph(),
	}
}

func (dg *DependencyGraph) LoadGriffonConfig(config *blocks.Config) error {

	// Add nodes to graph
	fmt.Println("Loading Griffon config into Dependency Graph")
	dg.AddNode(&config.Griffon)
	dg.blockPathToGraphID["griffon"] = config.Griffon.GraphID

	for blockType, blockMap := range config.Data {
		for blockName, block := range blockMap {
			fmt.Printf("...Loading Data %s::%s into Dependency Graph\n", blockType, blockName)
			dg.AddNode(block)
			dg.blockPathToGraphID[blocks.BuildBlockPath("data", string(blockType), blockName)] = block.ID()
		}
	}
	for blockType, blockMap := range config.Resources {
		for blockName, block := range blockMap {
			fmt.Printf("...Loading Resource %s::%s into Dependency Graph\n", blockType, blockName)
			dg.AddNode(block)
			dg.blockPathToGraphID[blocks.BuildBlockPath(string(blockType), blockName)] = block.ID()
		}
	}

	fmt.Println()
	// Add edges to graph
	nodes := dg.Nodes()
	for nodes.Next() {
		toNode := nodes.Node().(blocks.Block)
		fmt.Printf("...Processing node [%d] %s.%s\n", toNode.ID(), toNode.BlockType(), toNode.BlockName())
		deps := toNode.Dependencies()
		for _, dep := range deps {
			if fromNode := dg.Node(dg.blockPathToGraphID[dep]); fromNode != nil {
				fmt.Printf("... ...Adding edge from [%d] %s.%s to [%d] %s.%s\n", fromNode.ID(), fromNode.(blocks.Block).BlockType(), fromNode.(blocks.Block).BlockName(), toNode.ID(), toNode.BlockType(), toNode.BlockName())
				dg.SetEdge(dg.NewEdge(fromNode, toNode))
			}
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Printf("\nBlockPathToGraphID: %+v\n", dg.blockPathToGraphID)

	return nil
}

func (dg *DependencyGraph) BlockPathToGraphID(id int64) map[string]int64 {
	fmt.Printf("BlockPathToGraphID: %+v\n", dg.blockPathToGraphID)
	return dg.blockPathToGraphID
}

func (dg *DependencyGraph) GetSortedNodeIDs() ([]int64, error) {
	sortedDG, err := topo.Sort(dg)
	if err != nil {
		return nil, err
	}

	var sortedNodeIDs []int64
	// Add Griffon block to the beginning of the sorted list
	sortedNodeIDs = append(sortedNodeIDs, 0)
	for _, n := range sortedDG {
		if n.ID() != 0 {
			sortedNodeIDs = append(sortedNodeIDs, n.ID())
		}
	}
	return sortedNodeIDs, nil
}
