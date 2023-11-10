package hypergraph

func removeEdges(e map[int32]Edge, remIds map[int32]bool) []Edge {
    eCopy := []Edge{}
    for id, edge := range e {
        if remIds[id] {
            continue
        }
        eCopy = append(eCopy, copyEdge(edge))
    }
    return eCopy
}

func copyEdge(e Edge) Edge {
    vertices := []int32{}
    for i, v := range e.v {
        if v {
            vertices = append(vertices, i)
        }
    }
    return NewEdge(vertices...)
}

func removeVertices(v map[int32]Vertex, remIds map[int32]bool) []Vertex {
    vCopy := []Vertex{}
    for id, vertex := range v {
        if remIds[id] {
            continue
        }
        // mind that data can be non primitive
        vCopy = append(vCopy, Vertex{id: vertex.id, data: vertex.data})
    }
    return vCopy
}