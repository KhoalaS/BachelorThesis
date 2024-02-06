package hypergraph

func GetFrontierGraph(g *HyperGraph, level int, remId int32) *HyperGraph {
	g2 := NewHyperGraph()
	frontier := make(map[int32]bool)
	remEdge := g.Edges[remId]
	hashes := make(map[string]bool)

	for v := range g.Edges[remId].V {
		for e := range g.IncMap[v] {
			for w := range g.Edges[e].V {
				if !g.Edges[remId].V[w] {
					frontier[w] = true
					g2.AddVertex(w, 0)
				}
			}
		}
	}

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
			for e := range g.IncMap[v] {
				found := true
				for w := range g.Edges[e].V {
					if remEdge.V[w] {
						found = false
						break
					}
				}
				if found {
					hash := g.Edges[e].getHash()
					if !hashes[hash] {
						hashes[hash] = true
						g2.AddEdgeMapWLayer(g.Edges[e].V, i)
						for w := range g.Edges[e].V {
							if !frontier[w] {
								g2.AddVertex(w, i)
								nextFrontier[w] = true
							}
						}
					}
				}
			}
		}
		frontier = nextFrontier
	}

	return g2
}