package hypergraph

// Currently the rules will output a new hypergraph struct.
// A implementation that only manipulates the partial solution C and
// computes the 'current graph' derived from C is possibly better.

// Iterate over all edges of G.

// Then compare the current edge e to all other edges of G, excluding the edge itself
// and already removed edges.

// If all vertices of the edge e are present in the compared edge comp
// and if comp has more vertices, than we conclude that comp is a strict superset of e.

// Time Complexity: |E|^2 * d

func EdgeDominationRule(g HyperGraph) HyperGraph {
	remEdges := make(map[int32]bool)

	for eId, e := range g.Edges {
		for cId, comp := range g.Edges {
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

	newVertices := make([]Vertex, len(g.Vertices))
	var newEdges []Edge

	newEdges = removeEdges(g.Edges, remEdges)

	for i, v := range g.Vertices {
		newVertices[i] = v
	}
	return NewHyperGraph(newVertices, newEdges)
}

// Time Complexity: |E| * d

func RemoveEdgeRule(g HyperGraph, c map[int32]bool, t int) (HyperGraph, map[int32]bool) {
	remEdges := make(map[int32]bool)
	remVertices := make(map[int32]bool)
	cCopy := make(map[int32]bool)
	for k, v := range c {
        cCopy[k] = v
    }

	for id, e := range g.Edges {
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

	
	newEdges := removeEdges(g.Edges, remEdges)
	newVertices := removeVertices(g.Vertices, remVertices)
	
	return NewHyperGraph(newVertices, newEdges), cCopy
}

// Complexity: (|E| * d)^2 
// Currently a lot of overlap. 

func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool) (HyperGraph, map[int32]bool) {
	cCopy := make(map[int32]bool)
	for k, v := range c {
        cCopy[k] = v
    }

	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	var yz Edge
	var xDom int32 = -1
	
	for id, edge := range g.Edges {
		if len(edge.v) < 3 {
			continue
		}

		for x := range edge.v {

			cond := true

			for idComp, edgeComp := range g.Edges {
				if id == idComp {
					continue
				}
				if edgeComp.v[x] {
					sum := 0
					for vertex := range edge.v {
						if edgeComp.v[vertex] {
							sum += 1
						}
					}
					if sum < 2 {
						cond = false
						break
					}
				}
			}
			if cond {
				xDom = x	
				yz = edge
				break
			}
		}
		if xDom != -1 {
			break
		}
	}

	if xDom != -1 {
		for vertex := range yz.v {
			if vertex != xDom {
				remVertices[vertex] = true
				cCopy[vertex] = true
				for eId, edge := range g.Edges {
					if edge.v[vertex] {
						remEdges[eId] = true
					}
				}
			}
		}
		newVertices := removeVertices(g.Vertices, remVertices)
		newEdges := removeEdges(g.Edges, remEdges)
		return NewHyperGraph(newVertices, newEdges), cCopy
	}

	return g,c
}

