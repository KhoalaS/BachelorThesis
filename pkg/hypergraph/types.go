package hypergraph

var edgeCounter int = -1

type HyperGraph struct {
	vertices []Vertex
	edges []Edge
	adjMatrix [][]int
}

func NewHyperGraph(vertices []Vertex, edges []Edge) HyperGraph {
	vSize := len(vertices)
	eSize := len(edges)
	adjMatrix := make([][]int, vSize)

	for i := range vertices {		
		adjMatrix[i] = make([]int, eSize)
	}
	
	for _, e := range edges {
		for v := range e.v {
			adjMatrix[v][e.id] = 1
		}
	}
	return HyperGraph{vertices, edges, adjMatrix}
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
	edgeCounter++
	return Edge{edgeCounter, s}
}

const (
	TINY = 1
	SMALL = 2
)
