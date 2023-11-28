package graph

import (
	"fmt"
	"testing"

	"github.com/bensooraj/griffon/blocks"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/topo"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Error("NewGraph() returned nil")
	}

	r1, r2 := &blocks.RegionDataBlock{}, &blocks.RegionDataBlock{}
	r1.GraphID = 1
	r2.GraphID = 2
	p3, p4 := &blocks.PlanDataBlock{}, &blocks.PlanDataBlock{}
	p3.GraphID = 3
	p4.GraphID = 4
	o5, o6 := &blocks.OSDataBlock{}, &blocks.OSDataBlock{}
	o5.GraphID = 5
	o6.GraphID = 6

	for _, n := range []graph.Node{r1, r2, p3, p4, o5, o6} {
		fmt.Println("adding node", n.ID())
		_, new := g.NodeWithID(n.ID())
		if new {
			g.AddNode(n)
		} else {
			t.Errorf("NodeWithID(%d) already exists", n.ID())
		}
	}

	g.SetEdge(g.NewEdge(r1, p3))
	g.SetEdge(g.NewEdge(r1, p4))
	g.SetEdge(g.NewEdge(r2, p3))
	g.SetEdge(g.NewEdge(r2, p4))
	g.SetEdge(g.NewEdge(p3, o5))
	g.SetEdge(g.NewEdge(p4, o6))

	sortedNodes, err := topo.Sort(g)
	if err != nil {
		t.Error(err)
	}

	for _, n := range sortedNodes {
		switch n := n.(type) {
		case *blocks.RegionDataBlock:
			t.Log("blocks.RegionDataBlock", n.GraphID)
		case *blocks.PlanDataBlock:
			t.Log("blocks.PlanDataBlock", n.GraphID)
		case *blocks.OSDataBlock:
			t.Log("blocks.OSDataBlock", n.GraphID)
		}
	}
}
