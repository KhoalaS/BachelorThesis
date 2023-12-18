package alg

import (
	"testing"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func TestPotentialTriangleSituation(t *testing.T) {

	g := hypergraph.NewHyperGraph()

	for i:=0; i<5; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(1,2,3)
	g.AddEdge(1,3,4)
	g.AddEdge(1,4,2)

	x, ex := PotentialTriangle(g)
	
	if !ex {
		t.Fatal("Expected to find a potential triangle situation")
	} 

	if x != 1 {
		t.Fatalf("Vertex %d was removed, expected vertex 1.", x)
	}
}

func BenchmarkParallelPotentialTriangleSituation(b *testing.B) {
	g := hypergraph.GenerateTestGraph(1000000, 2000000, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelPotentialTriangle(g)
	}
}

func BenchmarkPotentialTriangleSituation(b *testing.B) {
	g := hypergraph.GenerateTestGraph(1000000, 2000000, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PotentialTriangle(g)
	}
}