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
	IncMap      map[int32]map[int32]bool
	MaxLayer int
}

type Vertex struct {
	Id   int32
	Data any
}

type Edge struct {
	V map[int32]bool
	Layer int
}

func (g *HyperGraph) AddVertex(id int32, data any) {
	if _, ex := g.Vertices[id]; !ex {
		g.Vertices[id] = Vertex{id, data}
	}
}

func (g *HyperGraph) RemoveVertex(id int32) bool {
	if _, ex := g.Vertices[id]; !ex {
		return false
	}
	delete(g.Vertices, id)
	return true
}

func (g *HyperGraph) AddEdge(eps ...int32) {
	e := Edge{V: make(map[int32]bool)}

	for _, ep := range eps {
		e.V[ep] = true
	}
	for ep := range e.V {
		if _, ex := g.IncMap[ep]; !ex {
			g.IncMap[ep] = make(map[int32]bool)
		}
		g.IncMap[ep][g.edgeCounter] = true
	}
	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) AddEdgeMap(eps map[int32]bool) {
	e := Edge{V: make(map[int32]bool)}

	for ep := range eps {
		e.V[ep] = true
		if _, ex := g.IncMap[ep]; !ex {
			g.IncMap[ep] = make(map[int32]bool)
		}
		g.IncMap[ep][g.edgeCounter] = true
	}

	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) AddEdgeMapWLayer(eps map[int32]bool, l int) {
	e := Edge{V: make(map[int32]bool), Layer: l}

	for ep := range eps {
		e.V[ep] = true
		if _, ex := g.IncMap[ep]; !ex {
			g.IncMap[ep] = make(map[int32]bool)
		}
		g.IncMap[ep][g.edgeCounter] = true
	}

	g.Edges[g.edgeCounter] = e
	g.edgeCounter++
}

func (g *HyperGraph) RemoveEdge(e int32) bool {
	if _, ex := g.Edges[e]; !ex {
		return false
	}

	for v := range g.Edges[e].V {
		delete(g.IncMap[v], e)

		if len(g.IncMap[v]) == 0 {
			delete(g.IncMap, v)
			g.RemoveVertex(v)
		}
	}

	delete(g.Edges, e)
	return true
}

func (g *HyperGraph) RemoveElem(elem int32) bool {
	if _, ex := g.Vertices[elem]; !ex {
		return false
	}

	if _, ex := g.IncMap[elem]; !ex {
		return false
	}

	for e := range g.IncMap[elem] {
		delete(g.Edges[e].V, elem)
		if len(g.Edges[e].V) == 0 {
			g.RemoveEdge(e)
		}
	}

	g.RemoveVertex(elem)
	delete(g.IncMap, elem)
	return true
}

func (g *HyperGraph) Deg(v int32) int {
	return len(g.IncMap[v])
}

// not final
func (g *HyperGraph) Copy() *HyperGraph {
	edges := make(map[int32]Edge)
	vertices := make(map[int32]Vertex)
	IncMap := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		edges[eId] = Edge{V: make(map[int32]bool)}
		for v := range e.V {
			edges[eId].V[v] = true
		}
	}

	for vId, v := range g.Vertices {
		vertices[vId] = Vertex{Id: vId, Data: v.Data}
	}

	return &HyperGraph{edgeCounter: g.edgeCounter, Vertices: vertices, Edges: edges, IncMap: IncMap}
}

func (g HyperGraph) String() string {
	s := "Vertices: \n\t"
	for _, v := range g.Vertices {
		s += fmt.Sprintf("%d,", v.Id)
	}

	s += "\nEdges:\n"
	for eId, e := range g.Edges {
		ids := []int32{}
		for id, val := range e.V {
			if !val {
				continue
			}
			ids = append(ids, id)
		}
		s += fmt.Sprintf("\t%d:%d\n", eId, ids)
	}
	s += "--------------------------\n"
	return s
}

func NewHyperGraph() *HyperGraph {
	vertices := make(map[int32]Vertex)
	edges := make(map[int32]Edge)
	incMap := make(map[int32]map[int32]bool)
	return &HyperGraph{Vertices: vertices, Edges: edges, IncMap: incMap}
}

func (g *HyperGraph) IsSimple() bool {
	for _, inc := range g.IncMap {
		if len(inc) == 3 {
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
			g.RemoveEdge(eId)
		} else {
			hashes[hash] = true
		}
	}
}

func (g *HyperGraph) Draw() {

}

func NewEdge(eps map[int32]bool) *Edge {
	e := Edge{V: make(map[int32]bool)}

	for ep := range eps {
		e.V[ep] = true
	}

	return &e
}

// Time Complexity: 2d + d*log(d)
func (e Edge) getHash() string {

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

type IntTuple struct {
	A int
	B int
}
