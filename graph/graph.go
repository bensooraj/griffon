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
			dg.blockPathToGraphID[blocks.BuildBlockPath(string(blockType), blockName)] = block.ID()
		}
	}
	for blockType, blockMap := range config.Resources {
		for blockName, block := range blockMap {
			fmt.Printf("...Loading Resource %s::%s into Dependency Graph\n", blockType, blockName)
			dg.AddNode(block)
			dg.blockPathToGraphID[blocks.BuildBlockPath(string(blockType), blockName)] = block.ID()
		}
	}

	// Add edges to graph
	nodes := dg.Nodes()
	for nodes.Next() {
		node := nodes.Node().(blocks.Block)
		d := node.Dependencies()
		for _, dep := range d {
			if dNode := dg.Node(dg.blockPathToGraphID[dep]); dNode != nil {
				dg.SetEdge(dg.NewEdge(dNode, node))
			}
		}
	}

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
	for _, n := range sortedDG {
		sortedNodeIDs = append(sortedNodeIDs, n.(blocks.Block).ID())
	}
	return sortedNodeIDs, nil
}
