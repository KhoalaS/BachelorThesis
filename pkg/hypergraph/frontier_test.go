package hypergraph

import (
	"testing"
)

func TestGetFrontierGraph(t *testing.T) {
	g := NewHyperGraph()

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	gf := GetFrontierGraph(g, 2, 0)
	t.Log(gf)

}
