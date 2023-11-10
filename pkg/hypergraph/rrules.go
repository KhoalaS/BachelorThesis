package hypergraph

import "fmt"

// Currently the rules will output a new hypergraph struct.
// A implementation that only manipulates the partial solution C and
// computes the 'current graph' derived from C is possibly better.

// Iterate over all edges of G.

// Then compare the current edge e to all other edges of G, excluding the edge itself
// and already removed edges.

// If all vertices of the edge e are present in the compared edge comp
// and if comp has more vertices, than we conclude that comp is a strict superset of e.

// Time Complexity: |E|^2 * d

func EdgeDominationRule(g HyperGraph, c map[int32]bool) {
	remEdges := make(map[int32]bool)
	l := len(g.Edges)
	counter := 0

	for _, e := range g.Edges {
		if len(e.v) == int(g.Degree) {
			continue
		}
		counter++
		fmt.Printf("%d/%d\r", counter, l)
		for cId, comp := range g.Edges {
			if remEdges[cId] || len(comp.v) <= len(e.v){
				continue
			}

			subset := true

			for id := range e.v {
				if !comp.v[id] {
					subset = false
					break
				}
			}

			if subset {
				remEdges[cId] = true
			}
		}
	}

	for eId := range remEdges {
		delete(g.Edges, eId)
	}
}

// Time Complexity: |E| * d

func RemoveEdgeRule(g HyperGraph, c map[int32]bool, t int) {
	remEdges := make(map[int32]bool)
	remVertices := make(map[int32]bool)

	for _, e := range g.Edges {
		if len(e.v) == t {
			for v := range e.v {
				remVertices[v] = true
				c[v] = true
			}
		}
	}

	fmt.Println(len(remVertices))

	for id, e := range g.Edges {
		for v := range e.v {
			if remVertices[v] {
				remEdges[id] = true
				break
			}		
		}
	}

	for eId := range remEdges {
		delete(g.Edges, eId)
	}	

	for vId := range remVertices {
		delete(g.Vertices, vId)
	}
	fmt.Println("delete finished")
}

// Complexity: (|E| * d)^2 
// Currently a lot of overlap. 

func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool) {
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
				c[vertex] = true
				for eId, edge := range g.Edges {
					if edge.v[vertex] {
						remEdges[eId] = true
					}
				}
			}
		}
		for eId := range remEdges {
			delete(g.Edges, eId)
		}	
	
		for vId := range remVertices {
			delete(g.Vertices, vId)
		}
	}
}

