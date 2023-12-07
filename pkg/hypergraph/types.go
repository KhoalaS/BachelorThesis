package hypergraph

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

type HyperGraph struct {
	Vertices map[int32]Vertex
	Edges map[int32]Edge
	edgeCounter int32
	Degree int
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

func NewHyperGraph() *HyperGraph {
	vertices := make(map[int32]Vertex)
	edges := make(map[int32]Edge)
	return &HyperGraph{Vertices: vertices, Edges: edges, Degree: 3}
}

func (g *HyperGraph) IsSimple() bool {
	// Time Complexity |E|*|d|
	degMap := make(map[int32]int32)
	simple := true

	for _, e := range g.Edges {
		for id := range e.v {
			degMap[id] = degMap[id]+1
			if degMap[id] == 3 {
				return false
			}
		}
	}
	return simple
}

func (g *HyperGraph) RemoveDuplicate() {
	hashes := make(map[uint32]bool)

	for eId, e := range g.Edges {
		hash := e.getHash()
		if hashes[hash] {
			delete(g.Edges, eId)
		}
	}
}

type Vertex struct {
	id int32
	data any
}

type Edge struct {
	v map[int32]bool
}

func NewEdge(eps map[int32]bool) *Edge {
	e := Edge{v: make(map[int32]bool)}
	
	for ep := range eps {
		e.v[ep] = true
	}

	return &e
}


// Time Complexity: 2d + d*log(d)
func (e *Edge) getHash() uint32 {
	h := xxhash.New32()

	arr := make([]int32, len(e.v))
	var i int32 = 0
	for ep := range e.v {
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
	r := strings.NewReader(in)
	io.Copy(h, r)

	return h.Sum32();
}	

const (
	TINY = 1
	SMALL = 2
)
