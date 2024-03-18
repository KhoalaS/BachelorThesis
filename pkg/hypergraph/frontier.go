package hypergraph

func ExpandFrontier(g *HyperGraph, level int, expand map[int32]bool) *HyperGraph {
	gNew := NewHyperGraph()

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range expand {
			for e := range g.IncMap[v] {
				if _, ex := gNew.Edges[e]; !ex {
					gNew.AddEdgeMapWithId(g.Edges[e].V, e)
					for w := range g.Edges[e].V {
						if _, ex2 := gNew.Vertices[w]; !ex2 {
							gNew.AddVertex(w, 0)
							nextFrontier[w] = true
						}
					}
				}
			}
		}
		if len(nextFrontier) == 0 {
			break
		}
		expand = nextFrontier
	}
	gNew.IncMap = g.IncMap
	gNew.AdjCount = g.AdjCount
	return gNew
}

func F3_ExpandFrontier(g *HyperGraph, remId int32, level int) *HyperGraph {
	frontier := make(map[int32]bool)
	remEdge := g.Edges[remId]

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
	return ExpandFrontier(g, level, frontier)
}

func ExtendFrontier(gf *HyperGraph, g *HyperGraph, level int, expand map[int32]bool) {

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range expand {
			for e := range g.IncMap[v] {
				if _, ex := gf.Edges[e]; !ex {
					gf.AddEdgeMapWithId(g.Edges[e].V, e)
					for w := range g.Edges[e].V {
						if _, ex2 := gf.Vertices[w]; !ex2 {
							gf.AddVertex(w, 0)
							nextFrontier[w] = true
						}
					}
				}
			}
		}
		if len(nextFrontier) == 0 {
			break
		}
		expand = nextFrontier
	}
}