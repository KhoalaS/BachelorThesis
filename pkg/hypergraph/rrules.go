package hypergraph

// Iterate over all edges of G.

// Then compare the current edge e to all other edges of G, excluding the edge itself
// and already removed edges.

// If all vertices of the edge e are present in the compared edge comp
// and if comp has more vertices, than we conclude that comp is a strict superset of e.

// Time Complexity: |E|^2 * 3

func EdgeDomination(g HyperGraph) HyperGraph{
	remEdges := make(map[int]bool)

	for _, e := range g.edges {
		for _, comp := range g.edges {
			if e.id == comp.id || remEdges[comp.id]{
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
				remEdges[e.id] = true
			}
		}
	}

	for e := range remEdges {
		g.edges = removeEdgeSlice(g.edges, e)
	}

	return g
}

func RemoveEdges(g HyperGraph, c map[int]bool, t int) (HyperGraph, map[int]bool) {
	remEdges := []int{}
	remVertex := []int{}
	cCopy := make(map[int]bool)
	for k, v := range c {
        cCopy[k] = v
    }

	for _, e := range g.edges {
		if len(e.v) == t {
			remEdges = append(remEdges, e.id)
			for v := range e.v {
				remVertex = append(remVertex, v)
				cCopy[v] = true
			}
		}
	}

	if len(remEdges) == 0 {
		return g, c
	}

	if t == SMALL {
		for _, v := range remVertex {
			for i, w := range g.adjMatrix[v] {
				if w == 1 {
					g.edges = removeEdgeSlice(g.edges, i)
				}
			}
		}
	} else {
		for _, id := range remEdges {
			g.edges = removeEdgeSlice(g.edges, id)
		}
	}
	

	for _, id := range remVertex {
		g.vertices = removeVertexSlice(g.vertices, id)
	}

	return g, cCopy
}

func VertexDomination(g HyperGraph) HyperGraph {
	
}