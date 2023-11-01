package hypergraph

import "log"

func removeEdgeSlice(s []Edge, id int) []Edge {
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
    s[i]= s[len(s)-1]
    return s[:len(s)-1]
}

func removeVertexSlice(s []Vertex, id int) []Vertex {
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
    s[i]= s[len(s)-1]
    return s[:len(s)-1]
}