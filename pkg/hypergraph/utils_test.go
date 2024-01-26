package hypergraph

import (
	"log"
	"testing"
)

func TestBinomialCoefficient(t *testing.T) {
	c := binomialCoefficient(1000, 3)
	exp := 166167000

	if c != exp {
		log.Fatalf("Wrong solution, got %d, expected %d.", c, exp)
	}
}

func TestGetFrontierGraph(t *testing.T) {
	g := NewHyperGraph()
	g.AddEdge(0,1)
	g.AddEdge(0,2)
	g.AddEdge(2,4)
	g.AddEdge(1,3)
	g.AddEdge(4,3)
	g.AddEdge(4,5,6)

	incMap := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		for v := range e.V {
			if _, ex := incMap[v]; !ex {
				incMap[v] = make(map[int32]bool)
			}
			incMap[v][eId] = true
		}
	}

	h := GetFrontierGraph(g, incMap, 3, 0)
	t.Log(h)

}