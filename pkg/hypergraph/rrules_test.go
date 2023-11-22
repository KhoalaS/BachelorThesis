package hypergraph

import (
	"testing"
)

func TestApproxVertexDominationRule(t *testing.T) {
	var vSize int32 = 8

	g := NewHyperGraph()

	var i int32 = 0
	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 2, 7)

	c := make(map[int32]bool)

	ApproxVertexDominationRule(g, c)

	// expected outcomes for c:
	// {2,7}, {1,2}, {0,7}, {0,1}
	if (c[2] && c[7]) || (c[1] && c[2]) || (c[0] && c[7]) || (c[0] && c[1]) {
		if len(c) != 2 {
			t.Fatalf("Partial solution is wrong.")
		}
	}

	if len(g.Edges) != 0 {
		t.Fatalf("Graph g has %d edges, the expected number is 0.", len(g.Edges))
	}

}

func TestEdgeDominationRule(t *testing.T) {
	var vSize int32 = 5
	g := NewHyperGraph()

	var i int32 = 0

	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 1)
	g.AddEdge(1, 4)

	c := make(map[int32]bool)

	EdgeDominationRule(g, c)

	// since edge 0 is a strict superset of edge 1, edge 0 will be removed by the rule.

	if len(c) != 0 {
		t.Fatalf("Partial Solution is not empty.")
	} else if len(g.Edges) != 2 {
		t.Fatalf("Graph g has %d edges, the expected number is 2.", len(g.Edges))
	}

	for _, edge := range g.Edges {
		if len(edge.v) != 2 {
			t.Fatalf("The wrong edge has been removed.")
		}
	}

}

func TestRemoveEdgeRule(t *testing.T) {
	var vSize int32 = 8
	g := NewHyperGraph()

	var i int32 = 0

	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 3, 2)
	g.AddEdge(1)
	g.AddEdge(2, 4)
	g.AddEdge(1, 6, 7)
	g.AddEdge(3, 6, 5)

	c := make(map[int32]bool)

	// this rule will remove edge (1) and will put vertex 1 into the partial solution
	// after putting vertex 1 into c1, edge (1,6,7) will be removed since vertex 1 is an element of it
	RemoveEdgeRule(g, c, TINY)

	// this rule will remove edge (2,4) and will put both vertex 2 and 4 into the partial solution c2
	// after putting vertex 2 and 4 into c2, edge (0,3,2) will be removed analogous to previous rule call
	RemoveEdgeRule(g, c, SMALL)

	psol := []int32{1, 2, 4}
	edgeSol := []int32{3, 6, 5}

	for _, v := range psol {
		if !c[v] {
			t.Fatalf("Vertex %d should have been part of solution.", v)
		}
	}

	if len(g.Edges) != 1 {
		t.Fatalf("Graph g has %d edges, the expected number is 1.", len(g.Edges))
	}

	for _, v := range edgeSol {
		for _, ep := range g.Edges {
			if !ep.v[v] {
				t.Fatalf("Vertex %d should have been part of solution.", v)
			}
		}
	}

}

func TestApproxVertexDominationRule3(t *testing.T) {
	var vSize int32 = 8
	g := NewHyperGraph()

	var i int32 = 0

	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(0, 2, 7)

	c := make(map[int32]bool)
	ApproxVertexDominationRule3(g, c)

	// possible solutions: [2,7], [0,2], [2,3]

	if len(c) != 2 {
		t.Fatalf("Partial solution is wrong.")
	}

	if !((c[2] && c[7]) || (c[0] && c[2]) || (c[2] && c[3])) {
		t.Fatalf("Partial solution is wrong.")
	}
}

func TestSmallTriangleRule(t *testing.T) {
	var vSize int32 = 7
	g := NewHyperGraph()

	var i int32 = 0

	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 4)
	g.AddEdge(4, 5, 6)

	c := make(map[int32]bool)

	SmallTriangleRule(g, c)

	if !(c[0] && c[1] && c[2]) {
		t.Fatalf("Partial solution is wrong.")
	}

	if len(g.Edges) != 1 {
		t.Fatalf("Graph g has %d edges, the expected number is 1.", len(g.Edges))
	}

	if _, ex := g.Edges[4]; !ex {
		t.Fatalf("The wrong edge has been removed.")
	}
}

func BenchmarkTinyEdgeRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)
	
	RemoveEdgeRule(g, c, TINY)
}

func BenchmarkSmallEdgeRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)
	
	RemoveEdgeRule(g, c, SMALL)
}

func BenchmarkSmallTriangleRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, false)
	c := make(map[int32]bool)

	SmallTriangleRule(g, c)
}

func BenchmarkApproxVertexDominationRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, false)
	c := make(map[int32]bool)

	ApproxVertexDominationRule3(g, c)
}