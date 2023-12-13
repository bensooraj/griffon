package graph

import (
	"log/slog"

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
	dg.AddNode(&config.Griffon)
	dg.blockPathToGraphID["griffon"] = config.Griffon.GraphID

	for blockType, blockMap := range config.Data {
		for blockName, block := range blockMap {
			slog.Info("[Graph] Loading Data into Dependency Graph", slog.Int64("graph_id", block.ID()), slog.String("block_type", string(blockType)), slog.String("block_name", blockName))
			dg.AddNode(block)
			dg.blockPathToGraphID[blocks.BuildBlockPath("data", string(blockType), blockName)] = block.ID()
		}
	}
	for blockType, blockMap := range config.Resources {
		for blockName, block := range blockMap {
			slog.Info("[Graph] Loading Resource into Dependency Graph", slog.Int64("graph_id", block.ID()), slog.String("block_type", string(blockType)), slog.String("block_name", blockName))
			dg.AddNode(block)
			dg.blockPathToGraphID[blocks.BuildBlockPath(string(blockType), blockName)] = block.ID()
		}
	}

	// Add edges to graph
	nodes := dg.Nodes()
	for nodes.Next() {
		toNode := nodes.Node().(blocks.Block)
		slog.Debug("[Graph] Adding edges for node", slog.Int64("graph_id", toNode.ID()), slog.String("block_type", string(toNode.BlockType())), slog.String("block_name", toNode.BlockName()))
		deps := toNode.Dependencies()
		for _, dep := range deps {
			if fromNode := dg.Node(dg.blockPathToGraphID[dep]); fromNode != nil {
				slog.Debug("[Graph] Adding edge/dependency", slog.Int64("from_graph_id", fromNode.ID()), slog.String("from_block_type", string(fromNode.(blocks.Block).BlockType())), slog.String("from_block_name", fromNode.(blocks.Block).BlockName()), slog.Int64("to_graph_id", toNode.ID()), slog.String("to_block_type", string(toNode.BlockType())), slog.String("to_block_name", toNode.BlockName()))
				dg.SetEdge(dg.NewEdge(fromNode, toNode))
			}
		}
	}
	return nil
}

func (dg *DependencyGraph) BlockPathToGraphID(id int64) map[string]int64 {
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
