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
		t.Fatalf("Graph g has %d edges, the expected number is 0.", len(g.Edges))
	}

}

func TestEdgeDominationRule(t *testing.T){
	var vSize int32 = 5
	var eSize int32 = 3

	vertices := make([]Vertex, vSize)
	edges := make([]Edge, eSize)
	
	var i int32 = 0
	for ; i < vSize; i++ {
		vertices[i] = NewVertex(i, 0)
	}

	edges[0] = NewEdge(0,1,2)
	edges[1] = NewEdge(0,1)
	edges[2] = NewEdge(1,4)

	g := NewHyperGraph(vertices, edges)
	c := make(map[int32]bool)
	
	g1, c1 := EdgeDominationRule(g, c)

	// since edge 0 is a strict superset of edge 1, edge 0 will be removed by the rule.

	if len(c1) != 0 {
		t.Fatalf("Partial Solution is not empty.")
	}else if len(g1.Edges) != 2 {
		t.Fatalf("Graph g has %d edges, the expected number is 2.", len(g.Edges))
	}

	for _, edge := range g1.Edges {
		if len(edge.v) != 2 {
			t.Fatalf("The wrong edge has been removed.")
		}
	}

}

func TestRemoveEdgeRule(t *testing.T) {
	var vSize int32 = 8
	var eSize int32 = 5

	vertices := make([]Vertex, vSize)
	edges := make([]Edge, eSize)
	
	var i int32 = 0
	for ; i < vSize; i++ {
		vertices[i] = NewVertex(i, 0)
	}

	edges[0] = NewEdge(0,3,2)
	edges[1] = NewEdge(1)
	edges[2] = NewEdge(2,4)
	edges[3] = NewEdge(1,6,7)
	edges[4] = NewEdge(3,6,5)


	g := NewHyperGraph(vertices, edges)
	c := make(map[int32]bool)
	

	// this rule will remove edge (1) and will put vertex 1 into the partial solution
	// after putting vertex 1 into c1, edge (1,6,7) will be removed since vertex 1 is an element of it
	g1, c1 := RemoveEdgeRule(g, c, TINY)

	// this rule will remove edge (2,4) and will put both vertex 2 and 4 into the partial solution c2
	// after putting vertex 2 and 4 into c2, edge (0,3,2) will be removed analogous to previous rule call
	g2, c2 := RemoveEdgeRule(g1, c1, SMALL)

	psol := []int32{1,2,4}
	edgeSol := []int32{3,6,5}

	for _, v := range psol {
		if !c2[v] {
			t.Fatalf("Vertex %d should have been part of solution.", v)
		}
	}

	if len(g2.Edges) != 1 {
		t.Fatalf("Graph g has %d edges, the expected number is 1.", len(g.Edges))
	}

	for _, v := range edgeSol {
		if !g2.Edges[0].v[v] {
			t.Fatalf("Vertex %d should have been part of solution.", v)
		}
	}

	
}