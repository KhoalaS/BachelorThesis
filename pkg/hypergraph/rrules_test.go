package hypergraph

import (
	"fmt"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgeDominationRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()

	for i := 0; i < 4; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 1)
	g.AddEdge(1, 3)

	EdgeDominationRule(g)

	// since edge 0 is a strict superset of edge 1, edge 0 will be removed by the rule.
	assert.Equal(3, len(g.Vertices), g.Vertices)
	assert.Equal(3, len(g.IncMap))
	assert.Equal(2, len(g.Edges))

	for _, edge := range g.Edges {
		assert.Equal(2, len(edge.V), "The wrong edge has been removed.")
	}

	// incidence
	assert.Equal(1, g.Deg(0))
	assert.Equal(2, g.Deg(1))
	assert.Equal(1, g.Deg(3))
}

func TestRemoveEdgeRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()

	for i := 0; i < 8; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(1)
	g.AddEdge(1, 6, 7)
	g.AddEdge(0, 3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(3, 6, 5)

	c := make(map[int32]bool)

	// this rule will remove edge (1) and will put vertex 1 into the partial solution
	// after putting vertex 1 into c1, edge (1,6,7) will be removed since vertex 1 is an element of it
	RemoveEdgeRule(g, c, TINY)
	assert.Equal(6, len(g.Vertices))
	assert.Equal(6, len(g.IncMap), g.IncMap)
	assert.Equal(3, len(g.Edges))

	assert.Equal(1, g.Deg(0))
	assert.Equal(2, g.Deg(2))
	assert.Equal(2, g.Deg(3))
	assert.Equal(1, g.Deg(4))
	assert.Equal(1, g.Deg(5))
	assert.Equal(1, g.Deg(6))

	// this rule will remove edge (2,4) and will put both vertex 2 and 4 into the partial solution c2
	// after putting vertex 2 and 4 into c2, edge (0,3,2) will be removed analogous to previous rule call
	// only edge (3,6,5) remains
	RemoveEdgeRule(g, c, SMALL)
	assert.Equal(3, len(g.Vertices))
	assert.Equal(3, len(g.IncMap), g.IncMap)
	assert.Equal(1, len(g.Edges))

	assert.Equal(1, g.Deg(3))
	assert.Equal(1, g.Deg(6))
	assert.Equal(1, g.Deg(5))
}

func TestApproxVertexDominationRule(t *testing.T) {
	assert := assert.New(t)
	g := NewHyperGraph()

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 5)

	c := make(map[int32]bool)
	ApproxVertexDominationRule(g, c)

	assert.Equal(1, len(g.Edges))
	assert.Equal(2, len(g.Vertices))
	assert.Equal(2, len(g.IncMap))
	assert.Equal(2, len(c))

	assert.Equal(1, g.Deg(3))
	assert.Equal(1, g.Deg(5))

}

func TestApproxDoubleVertexDominationRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()

	for i := 0; i < 4; i++ {
		g.AddVertex(int32(i), 0)
	}

	c := make(map[int32]bool)

	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)

	ApproxDoubleVertexDominationRule(g, c)
	assert.Equal(2, len(c))
	assert.Equal(0, len(g.Vertices))
	assert.Equal(0, len(g.IncMap))
	assert.Equal(0, len(g.Edges))
	assert.Equal(true, c[3])
	assert.Equal(true, c[2])
}

func TestSmallTriangleRule(t *testing.T) {
	assert := assert.New(t)
	g := NewHyperGraph()

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4, 5)

	c := make(map[int32]bool)

	SmallTriangleRule(g, c)

	assert.Equal(3, len(c))
	assert.Equal(3, len(g.Vertices))
	assert.Equal(1, len(g.Edges))
	assert.Equal(3, len(g.IncMap))

	assert.Equal(1, g.Deg(3))
	assert.Equal(1, g.Deg(4))
	assert.Equal(1, g.Deg(5))
}

func TestSmallEdgeDegreeTwoRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()
	c := make(map[int32]bool)

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2, 3)
	g.AddEdge(3, 4, 5)

	SmallEdgeDegreeTwoRule(g, c)

	assert.Equal(4, len(c))
	assert.Equal(0, len(g.Vertices))
	assert.Equal(0, len(g.Edges))
	assert.Equal(0, len(g.IncMap))

}

func TestF3TargetLowDegree(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()
	for i := 0; i < 7; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2, 3)
	g.AddEdge(4, 3, 5)
	g.AddEdge(3, 6)

	c := make(map[int32]bool)
	F3TargetLowDegree(g, c)

	assert.Equal(true, (c[3] && c[4] && c[5]), c)
	assert.Equal(true, !(c[0] && c[1] && c[2]))
	assert.Equal(1, len(g.Edges))
	assert.Equal(2, len(g.Vertices))
	assert.Equal(2, len(g.IncMap))

	assert.Equal(1, g.Deg(0))
	assert.Equal(1, g.Deg(1))
}

func TestVertexDominationRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()

	for i := 0; i < 5; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(2, 3, 4)
	g.AddEdge(2)

	c := make(map[int32]bool)
	VertexDominationRule(g, c)

	assert.Equal(0, len(c))
	assert.Equal(1, len(g.Vertices))
	assert.Equal(1, len(g.IncMap))
	assert.Equal(1, len(g.Edges))
	assert.Equal(1, g.Deg(2))
}

func TestExtendedTriangleRule(t *testing.T) {
	assert := assert.New(t)

	g := NewHyperGraph()

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2, 3)
	g.AddEdge(1, 2)
	g.AddEdge(3, 4, 5)
	g.AddEdge(4, 5)

	c := make(map[int32]bool)
	ExtendedTriangleRule(g, c)

	assert.Equal(2, len(g.Vertices))
	assert.Equal(2, len(g.IncMap))
	assert.Equal(1, len(g.Edges))
	assert.Equal(4, len(c))

	assert.Equal(1, g.Deg(4))
	assert.Equal(1, g.Deg(5))
}

func TestF3Rule(t *testing.T) {
	assert := assert.New(t)
	g := NewHyperGraph()

	for i := 0; i < 6; i++ {
		g.AddVertex(int32(i), 0)
	}

	g.AddEdge(0, 1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)

	c := make(map[int32]bool)

	F3Rule(g, c)

	assert.Equal(1, len(g.Edges))
	assert.Equal(2, len(g.Vertices))
	assert.Equal(2, len(g.IncMap))
	assert.Equal(3, len(c))

	assert.Equal(1, g.Deg(4))
	assert.Equal(1, g.Deg(5))

}

func BenchmarkF3TargetLowDegree(b *testing.B) {
	g := TestGraph(100000, 200000, false)
	c := make(map[int32]bool)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		F3TargetLowDegree(g, c)
	}
}

func BenchmarkSmallDegreeTwoRule(b *testing.B) {
	g := TestGraph(100000, 200000, false)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SmallEdgeDegreeTwoRule(g, c)
	}
}

func BenchmarkTinyEdgeRule(b *testing.B) {
	g := TestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveEdgeRule(g, c, SMALL)
	}
}

func BenchmarkSmallEdgeRule(b *testing.B) {
	g := TestGraph(1000000, 2000000, true)
	c := make(map[int32]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveEdgeRule(g, c, SMALL)
	}
}

func BenchmarkEdgeDominationRule(b *testing.B) {
	g := TestGraph(1000000, 2000000, false)

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
	g := TestGraph(100000, 2000000, false)
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
	g := TestGraph(100000, 200000, false)
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
		ApproxVertexDominationRule(g, c)
	}

}

func BenchmarkApproxDoubleVertexDominationRule(b *testing.B) {
	g := TestGraph(100000, 200000, false)
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
