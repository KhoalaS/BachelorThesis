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
	Vertices map[int]Vertex
	Edges map[int]Edge
}

func (g HyperGraph) Print() {
	fmt.Print("Vertices: \n\t")
	for _, v := range g.Vertices {
		fmt.Printf("%d,", v.id)
	}

	fmt.Println("\nEdges:")
	for _, e := range g.Edges {
		ids := []int{}
		for id := range e.v {
			ids = append(ids, id)
		}
		fmt.Printf("\t%d\n",ids)
	}
	fmt.Println("--------------------------")
} 

// Hypergraph Constructor
// Arguments: Vertex slice v, Edge slice e
// We map the vertex id to the vertex itself and the edge ids are numbered and mapped
// to ids 0 to |E|-1.
// We explicity do not ensure that the resulting hypergraph is a decoupled
// from the inputs. That should be done before calling the constructor.
func NewHyperGraph(v []Vertex, e []Edge) HyperGraph {
	
	vertices := make(map[int]Vertex)
	edges := make(map[int]Edge)
	
	for _, vertex := range v {
		vertices[vertex.id] = vertex		
	}

	for i, edge := range e {
		edges[i] = edge
	}

	return HyperGraph{vertices, edges}
}

func (g HyperGraph) IsSimple() bool {
	// Time Complexity |E|*|d|
	degMap := make(map[int]int)
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
	id int
	data any
}

func NewVertex(id int, data any) Vertex {
	return Vertex{id, data}
}

type Edge struct {
	v map[int]bool
}

func NewEdge(v... int) Edge {
	s := make(map[int]bool)
	for _, v := range v {
		s[v] = true
	}
	return Edge{s}
}

const (
	TINY = 1
	SMALL = 2
)
