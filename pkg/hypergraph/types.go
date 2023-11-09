package hypergraph

import "fmt"


type HyperGraph struct {
	vertices []Vertex
	edges []Edge
	adjMatrix [][]int
	idIndexMap map[int]int
}

func (g HyperGraph) Print() {
	fmt.Println("Vertices:")
	for _, v := range g.vertices {
		fmt.Printf("\tID: %d\n", v.id)
	}

	fmt.Println("Edges:")
	for _, e := range g.edges {
		ids := []int{}
		for id := range e.v {
			ids = append(ids, id)
		}
		fmt.Printf("\t%d\n",ids)
	}
} 

func NewHyperGraph(vertices []Vertex, edges []Edge) HyperGraph {
	vSize := len(vertices)
	eSize := len(edges)
	adjMatrix := make([][]int, vSize)
	idIndexMap := make(map[int]int)

	// needs an object to map indices to vertex ids, 
	// in case a vertex deletion shifts the indices/ids
	// Example: deletion of vertex 0 shifts the indices when a new graph is created
	//		0	1
	//	0	1   1
	//	1   0   1

	// or keep the original size of the graph intact, which seems more complicated



	for i, v := range vertices {
		idIndexMap[v.id] = i
	}

	for i := range vertices {		
		adjMatrix[i] = make([]int, eSize)
	}

	for i, e := range edges {
		edges[i].id = i
		for v := range e.v {
			adjMatrix[idIndexMap[v]][i] = 1
		}
	}
	return HyperGraph{vertices, edges, adjMatrix, idIndexMap}
}

func (g HyperGraph) GetEntry(id int) []int {
	return g.adjMatrix[g.idIndexMap[id]]
}

func (g HyperGraph) IsSimple() bool {
	// Time Complexity |V|*|d|
	for i := range g.vertices {
		row := g.adjMatrix[i]
		degree := 0
		for _, val := range row {
			degree += val
		}
		if(degree > 2){
			return false
		}
	}
	return true
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
	id int
	v map[int]bool
}

func NewEdge(v... int) Edge {
	s := make(map[int]bool)
	for _, v := range v {
		s[v] = true
	}
	return Edge{0, s}
}

const (
	TINY = 1
	SMALL = 2
)
