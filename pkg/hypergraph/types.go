package hypergraph

import (
	"fmt"
)

// Simple Hypergraph structure
// We use maps for vertices and edges.
// An earlier approach utilized slices for both edges and vertices and an additional
// incidence matrix. While some operations are faster on the incidence matrix,
// the algorithms for the reduction rules often rely on iterating through all pairs of edges.
// Or checking if a vertex is an element of an edge.

// Since the use cases of the incidence matrix was so low, we decided to remove it from the 
// data structure and only use maps for both vertices and edges.

type HyperGraph struct {
	Vertices map[int32]Vertex
	Edges map[int32]Edge
	edgeCounter int32
	Degree int32
}

func (g *HyperGraph) AddVertex(id int32, data any) {
	g.Vertices[id] = Vertex{id, data}
}

func (g *HyperGraph) AddEdge(eps... int32) {
	e := Edge{v: make(map[int32]bool)}
	
	for _, ep := range eps {
		e.v[ep] = true
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) AddEdgeMap(eps map[int32]bool) {
	e := Edge{v: make(map[int32]bool)}
	
	for ep := range eps {
		e.v[ep] = true
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g HyperGraph) Print() {
	fmt.Print("Vertices: \n\t")
	for _, v := range g.Vertices {
		fmt.Printf("%d,", v.id)
	}

	fmt.Println("\nEdges:")
	for eId, e := range g.Edges {
		ids := []int32{}
		for id := range e.v {
			ids = append(ids, id)
		}
		fmt.Printf("\t%d:%d\n",eId, ids)
	}
	fmt.Println("--------------------------")
} 

func NewHyperGraph() HyperGraph {
	vertices := make(map[int32]Vertex)
	edges := make(map[int32]Edge)
	return HyperGraph{Vertices: vertices, Edges: edges, Degree: 3}
}

func (g HyperGraph) IsSimple() bool {
	// Time Complexity |E|*|d|
	degMap := make(map[int32]int32)
	simple := true
	outerBreak := false

	for _, e := range g.Edges {
		for id := range e.v {
			degMap[id] = degMap[id]+1
			if degMap[id] == 3 {
				simple = false
				outerBreak = true
				break
			}
		}
		if outerBreak {
			break
		}
	}
	
	return simple
}

// TODO: Dont use New[Vertex/Edge] functions but methods on the graph 
// itself like addVertex(). 

type Vertex struct {
	id int32
	data any
}

type Edge struct {
	v map[int32]bool
}

const (
	TINY = 1
	SMALL = 2
)
