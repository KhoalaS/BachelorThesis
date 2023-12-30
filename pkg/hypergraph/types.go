package hypergraph

import (
	"fmt"
	"sort"
	"strconv"
)

type HyperGraph struct {
	Vertices    map[int32]Vertex
	Edges       map[int32]Edge
	edgeCounter int32
	Degree      int
	VDeg 		map[int32]int32
}

func (g *HyperGraph) AddVertex(id int32, data any) {
	g.Vertices[id] = Vertex{id, data}
}

func (g *HyperGraph) AddEdge(eps ...int32) {
	e := Edge{V: make(map[int32]bool)}

	for _, ep := range eps {
		e.V[ep] = true
	}
	for ep := range e.V {
		g.VDeg[ep]++
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) RemoveEdge(id int32) bool {
	if _, ex := g.Edges[id]; !ex {
		return false
	}

	for v := range g.Edges[id].V {
		g.VDeg[v]--
	}

	delete(g.Edges, id)

	return true
}

func (g *HyperGraph) AddEdgeMap(eps map[int32]bool) {
	e := Edge{V: make(map[int32]bool)}

	for ep := range eps {
		e.V[ep] = true
	}
	for ep := range e.V {
		g.VDeg[ep]++
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) AddEdgeArr(eps []int32) {
	e := Edge{V: make(map[int32]bool)}

	for _, ep := range eps {
		e.V[ep] = true
	}
	for ep := range e.V {
		g.VDeg[ep]++
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) Copy() *HyperGraph {
	edges := make(map[int32]Edge)
	vertices := make(map[int32]Vertex)
	VDeg := make(map[int32]int32)

	for eId, e := range g.Edges {
		edges[eId] = Edge{V: make(map[int32]bool)}
		for v := range e.V {
			edges[eId].V[v] = true
		}
	}

	for vId, v := range g.Vertices {
		vertices[vId] = Vertex{id: vId, data: v.data}
	}

	for vId, degree := range g.VDeg {
		VDeg[vId] = degree
	}

	return &HyperGraph{edgeCounter: g.edgeCounter, Vertices: vertices, Edges: edges, Degree: g.Degree, VDeg: VDeg}
}

func (g HyperGraph) Print() {
	fmt.Print("Vertices: \n\t")
	for _, v := range g.Vertices {
		fmt.Printf("%d,", v.id)
	}

	fmt.Println("\nEdges:")
	for eId, e := range g.Edges {
		ids := []int32{}
		for id, val := range e.V {
			if !val {
				continue
			}
			ids = append(ids, id)
		}
		fmt.Printf("\t%d:%d\n", eId, ids)
	}
	fmt.Println("--------------------------")
}

func NewHyperGraph() *HyperGraph {
	vertices := make(map[int32]Vertex)
	edges := make(map[int32]Edge)
	vdeg := make(map[int32]int32)
	return &HyperGraph{Vertices: vertices, Edges: edges, Degree: 3, VDeg: vdeg}
}

func (g *HyperGraph) IsSimple() bool {
	for _, degree := range g.VDeg {
		if degree == 3 {
			return false
		}
	}
	return true
}

func (g *HyperGraph) RemoveDuplicate() {
	hashes := make(map[string]bool)

	for eId, e := range g.Edges {
		hash := e.getHash()
		if hashes[hash] {
			for v := range g.Edges[eId].V {
				g.VDeg[v]--
			}
			delete(g.Edges, eId)
		} else {
			hashes[hash] = true
		}
	}
}

func (g* HyperGraph) Draw() {

}

type Vertex struct {
	id   int32
	data any
}

type Edge struct {
	V map[int32]bool
}

func NewEdge(eps map[int32]bool) *Edge {
	e := Edge{V: make(map[int32]bool)}

	for ep := range eps {
		e.V[ep] = true
	}

	return &e
}

// Time Complexity: 2d + d*log(d)
func (e *Edge) getHash() string {

	arr := make([]int32, len(e.V))
	var i int32 = 0
	for ep := range e.V {
		arr[i] = ep
		i++
	}
	sort.Slice(arr, func(i2, j int) bool {
		return arr[i2] < arr[j]
	})

	in := "|"

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}

	return in
}

const (
	TINY  = 1
	SMALL = 2
)

