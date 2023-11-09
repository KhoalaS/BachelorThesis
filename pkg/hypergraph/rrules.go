package hypergraph

// Iterate over all edges of G.

// Then compare the current edge e to all other edges of G, excluding the edge itself
// and already removed edges.

// If all vertices of the edge e are present in the compared edge comp
// and if comp has more vertices, than we conclude that comp is a strict superset of e.

// Time Complexity: |E|^2 * d

func EdgeDomination(g HyperGraph) HyperGraph {
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
				remEdges[comp.id] = true
			}
		}
	}

	newVertices := make([]Vertex, len(g.vertices))
	var newEdges []Edge

	for e := range remEdges {
		newEdges = removeEdgeSlice(g.edges, e)
	}

	for i, v := range g.vertices {
		newVertices[i] = v
	}
	return NewHyperGraph(newVertices, newEdges)
}

// Time Complexity: |E| * d

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

	var newEdges []Edge
	var newVertices []Vertex


	if t == SMALL {
		for _, v := range remVertex {
			for i, w := range g.GetEntry(v) {
				if w == 1 {
					newEdges = removeEdgeSlice(g.edges, i)
				}
			}
		}
	} else {
		for _, id := range remEdges {
			newEdges = removeEdgeSlice(g.edges, id)
		}
	}
	

	for _, id := range remVertex {
		newVertices = removeVertexSlice(g.vertices, id)
	}

	return NewHyperGraph(newVertices, newEdges), cCopy
}

func ApproxVertexDomination(g HyperGraph, V int, d int)  {
	for _, e := range g.edges {
		if len(e.v) < d {
			continue
		}
		sumSet := make(map[int][]int)

		for _, f := range g.edges {
			if e.id == f.id {
				continue
			}
			sum := make([]int, V)

			for i := range e.v {
				if f.v[i] {
					sum[g.idIndexMap[i]] = 2
				} else {
					sum[g.idIndexMap[i]] = 1
				}
			}
			sumSet[f.id] = sum
		}

	}

}