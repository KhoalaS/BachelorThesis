package hypergraph

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

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

	EdgeDominationRule(g)

	// since edge 0 is a strict superset of edge 1, edge 0 will be removed by the rule.

	if len(g.Edges) != 2 {
		t.Fatalf("Graph g has %d edges, the expected number is 2.", len(g.Edges))
	}

	for _, edge := range g.Edges {
		if len(edge.V) != 2 {
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
			if !ep.V[v] {
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
	g.AddEdge(1, 4)

	c := make(map[int32]bool)
	ApproxVertexDominationRule(g, c, false)

	// possible solutions: [2,7], [0,2], [2,3]

	if len(c) != 2 {
		t.Fatalf("Partial solution is wrong.")
	}

	if !((c[2] && c[7]) || (c[0] && c[2]) || (c[2] && c[3])) {
		t.Fatalf("Partial solution is wrong.")
	}
	log.Println(g)
}

func TestApproxDoubleVertexDominationRule(t *testing.T) {
	g := NewHyperGraph()
	g.AddEdge(1, 2, 3)
	g.AddEdge(2, 3, 4)
	g.AddEdge(2, 5, 6)
	g.AddEdge(1, 6)

	// possible solutions [1,2], [2,6], [3,5], [3,6]

	c := make(map[int32]bool)

	ApproxDoubleVertexDominationRule(g, c)
	if len(c) != 2 {
		t.Fatalf("Partial solution is wrong.")
	}

	if !((c[1] && c[2]) || (c[2] && c[6]) || (c[3] && (c[5] || c[6]))) {
		t.Fatalf("Partial solution is wrong.")
	}
}

func TestApproxDoubleVertexDominationRule2(t *testing.T) {
	g := NewHyperGraph()
	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(1, 2, 3)
	g.AddEdge(1, 4, 5)
	g.AddEdge(0, 5)

	// possible solutions [0,1], [1,5], [2,4], [2,5]

	c := make(map[int32]bool)

	ApproxDoubleVertexDominationRule2(g, c)
	if len(c) != 2 {
		t.Fatalf("Partial solution is wrong.")
	}

	if !((c[0] && c[1]) || (c[1] && c[5]) || (c[2] && (c[4] || c[5]))) {
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

func TestSmallEdgeDegreeTwoRule(t *testing.T) {
	var i int32 = 0
	g := NewHyperGraph()
	c := make(map[int32]bool)

	for ; i < 6; i++ {
		g.AddVertex(i, 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2, 3)
	g.AddEdge(3, 4, 5)

	exec := SmallEdgeDegreeTwoRule(g, c)
	log.Println(g)
	if exec != 1 {
		log.Fatalf("Number of rule executions is wrong. Expected %d, got %d.", 1, exec)
	}

	if len(g.Vertices) != 2 || !c[0] || !c[3] || !c[4] || !c[5] {
		t.Fatalf("Partial solution is wrong.")
	}
}

func BenchmarkSmallDegreeTwoRule(b *testing.B) {
	g := GenerateTestGraph(100000, 200000, false)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SmallEdgeDegreeTwoRule(g, c)
	}
}

func BenchmarkTinyEdgeRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveEdgeRule(g, c, SMALL)
	}
}

func BenchmarkSmallEdgeRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveEdgeRule(g, c, SMALL)
	}
}

func BenchmarkEdgeDominationRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, false)

	f, err := makeProfile("edgeDom")
	if err != nil {
		b.Fatal("Could not create cpu profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EdgeDominationRule(g)
	}
}

func BenchmarkSmallTriangleRule(b *testing.B) {
	g := GenerateTestGraph(100000, 2000000, false)
	c := make(map[int32]bool)

	f, err := makeProfile("smallTriangleRule")
	if err != nil {
		b.Fatal("Could not create cpu profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SmallTriangleRule(g, c)
	}
}

func BenchmarkApproxVertexDominationRule(b *testing.B) {
	g := GenerateTestGraph(1000000, 2000000, false)
	g.RemoveDuplicate()

	c := make(map[int32]bool)
	name := "approxVertexDom"
	f, err := makeProfile(name)
	if err != nil {
		b.Fatal("Could not create cpu profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproxVertexDominationRule(g, c, false)
	}

}

func BenchmarkApproxDoubleVertexDominationRule(b *testing.B) {
	g := GenerateTestGraph(100000, 200000, false)
	g.RemoveDuplicate()
	c := make(map[int32]bool)

	name := "approxDoubleVertexDom"

	f, err := makeProfile(name)
	if err != nil {
		b.Fatal("Could not create cpu profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproxDoubleVertexDominationRule(g, c)
	}

	m, err := os.Create(fmt.Sprintf("../../profiles/mem_%s.prof", name))
	if err != nil {
		b.Fatal("could not create memory profile: ", err)
	}
	defer m.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(m); err != nil {
		b.Fatal("could not write memory profile: ", err)
	}

}

func BenchmarkApproxDoubleVertexDominationRule2(b *testing.B) {
	g := GenerateTestGraph(10000, 100000, false)
	g.RemoveDuplicate()
	c := make(map[int32]bool)

	name := "approxDoubleVertexDom2"

	f, err := makeProfile(name)
	if err != nil {
		b.Fatal("Could not create cpu profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproxDoubleVertexDominationRule2(g, c)
	}

	m, err := os.Create(fmt.Sprintf("../../profiles/mem_%s.prof", name))
	if err != nil {
		b.Fatal("could not create memory profile: ", err)
	}
	defer m.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(m); err != nil {
		b.Fatal("could not write memory profile: ", err)
	}

}

func makeProfile(name string) (*os.File, error) {
	f, err := os.Create(fmt.Sprintf("../../profiles/benchmark_%s.prof", name))
	if err != nil {
		return nil, err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return f, err
	}

	return f, nil
}

func stopProfiling(f *os.File) {
	pprof.StopCPUProfile()
	f.Close()
}
