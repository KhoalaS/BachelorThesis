package hypergraph

import "testing"

func TestApproxVertexDominationRule(t *testing.T){
	var vSize int32 = 8
	var eSize int32 = 2

	vertices := make([]Vertex, vSize)
	edges := make([]Edge, eSize)
	
	var i int32 = 0
	for ; i < vSize; i++ {
		vertices[i] = NewVertex(i, 0)
	}

	edges[0] = NewEdge(0,1,2)
	edges[1] = NewEdge(0,2,7)

	g := NewHyperGraph(vertices, edges)
	c := make(map[int32]bool)

	g1, c1 := ApproxVertexDominationRule(g, c)

	// expected outcomes for c1:
	// {2,7}, {1,2}, {0,7}, {0,1}
	if (c1[2] && c1[7]) || (c1[1] && c1[2]) || (c1[0] && c1[7]) || (c1[0] && c1[1]) {
		t.Log(c1)
		if len(c1) != 2 {
			t.Fatalf("Partial solution is wrong.")
		}
	}

	if len(g1.Edges) != 0 {
		t.Fatalf("Number of edges are incorrect, there are %d edges in the current graph.",len(g1.Edges))
	}

}