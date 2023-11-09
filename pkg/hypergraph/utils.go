package hypergraph

func removeEdges(e map[int]Edge, remIds map[int]bool) []Edge {
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
    vertices := []int{}
    for i, v := range e.v {
        if v {
            vertices = append(vertices, i)
        }
    }
    return NewEdge(vertices...)
}

func removeVertices(v map[int]Vertex, remIds map[int]bool) []Vertex {
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