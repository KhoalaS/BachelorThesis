package hypergraph

import "log"

func removeEdgeSlice(s []Edge, id int) []Edge {
    newEdges := make([]Edge, len(s))
	for i, e := range s {
		newEdges[i] = e
	}

    i := -1
    for j := 0; j < len(s); j++ {
        if s[j].id == id{
            i = j
            break
        }
    }
    if(i == -1){
        log.Fatalf("Edge with id %d does not exist", id)
    }
    newEdges[i] = s[len(s)-1]
    return newEdges[:len(s)-1]
}

func removeVertexSlice(s []Vertex, id int) []Vertex {
    newVertices := make([]Vertex, len(s))
	for i, e := range s {
		newVertices[i] = e
	}

    i := -1
    for j := 0; j < len(s); j++ {
        if s[j].id == id{
            i = j
            break
        }
    }
    if(i == -1){
        log.Fatalf("Vertex with id %d does not exist", id)
    }
    newVertices[i]= s[len(s)-1]
    return newVertices[:len(s)-1]
}