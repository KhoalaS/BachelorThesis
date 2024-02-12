package hypergraph

import (
	"fmt"
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


	// This will start the frontier at the vertices of edges incident to 0 and 1 (0 and 1 excluded).
	// This results in a initial frontier with the vertex 2.
	// Since we want to go 2 levels deep, only the edges (2,3) and (3,4) should be in
	// in the new graph gf.
	gf := GetFrontierGraph(g, 2, 0)
	ExpandFrontier(gf, g, 2)
	t.Log(gf)
	for _, v := range gf.Vertices {
		fmt.Println(v.Id, v.Data)
	}
	t.Log(gf.VertexFrontier)
}
