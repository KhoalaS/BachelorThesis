package hypergraph

// Iterate over all edges of G.

// Then compare the current edge e to all other edges of G, excluding the edge itself
// and already removed edges.

// If all vertices of the edge e are present in the compared edge comp
// and if comp has more vertices, than we conclude that comp is a strict superset of e.

// Time Complexity: |E|^2 * d

func EdgeDomination(g HyperGraph) HyperGraph {
	remEdges := make(map[int]bool)

	for eId, e := range g.edges {
		for cId, comp := range g.edges {
			if eId == cId || remEdges[cId]{
				continue
			}

			subset := true

			for id := range e.v {
				if !comp.v[id] {
					subset = false
					break
				}
			}

			if subset && len(comp.v) > len(e.v) {
				remEdges[cId] = true
			}
		}
	}

	newVertices := make([]Vertex, len(g.vertices))
	var newEdges []Edge

	newEdges = removeEdges(g.edges, remEdges)

	for i, v := range g.vertices {
		newVertices[i] = v
	}
	return NewHyperGraph(newVertices, newEdges)
}

// Time Complexity: |E| * d

func RemoveEdgeRule(g HyperGraph, c map[int]bool, t int) (HyperGraph, map[int]bool) {
	remEdges := make(map[int]bool)
	remVertices := make(map[int]bool)
	cCopy := make(map[int]bool)
	for k, v := range c {
        cCopy[k] = v
    }

	for id, e := range g.edges {
		if len(e.v) == t {
			remEdges[id] = true
			for v := range e.v {
				remVertices[v] = true
				cCopy[v] = true
			}
		}
	}

	if len(remEdges) == 0 {
		return g, c
	}

	
	newEdges := removeEdges(g.edges, remEdges)
	newVertices := removeVertices(g.vertices, remVertices)
	
	return NewHyperGraph(newVertices, newEdges), cCopy
}

func ApproxVertexDomination(g HyperGraph, V int, d int)  {
}