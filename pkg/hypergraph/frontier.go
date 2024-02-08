package hypergraph

func GetFrontierGraph(g *HyperGraph, level int, remId int32) *HyperGraph {
	g2 := NewHyperGraph()
	frontier := make(map[int32]bool)
	remEdge := g.Edges[remId]
	hashes := make(map[string]bool)

	for v := range remEdge.V {
		for e := range g.IncMap[v] {
			for w := range g.Edges[e].V {
				if !remEdge.V[w] {
					frontier[w] = true
				}
			}
		}
	}

	// remove the edges adjacent to remEdge
	for v := range remEdge.V {
		for e := range g.IncMap[v] {
			g.RemoveEdge(e)
		}
	}


	fmt.Println("Init frontier:", frontier)

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
						g2.AddEdgeMapWLayer(g.Edges[e].V, i, e)
						for w := range g.Edges[e].V {
							g2.AddVertex(w, i+1)
							if !frontier[w] {
								g2.AddVertex(w, i)
								nextFrontier[w] = true
							}
						}
					}
				}
			}
		}
		g2.SetMaxLayer(i)
		if len(nextFrontier) == 0 {
			break
		}
		frontier = nextFrontier
	}
	g2.IncMap = nil
	g2.IncMap = g.IncMap

	return g2
}