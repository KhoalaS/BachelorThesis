package hypergraph

import "fmt"

func GetFrontierGraph(g *HyperGraph, level int, remId int32) *HyperGraph {
	g2 := NewHyperGraph()
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

	fmt.Println("Init frontier:", frontier)

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
			for e := range g.IncMap[v] {
				if _, ex := g2.Edges[e]; !ex {
					g2.AddEdgeMapWithId(g.Edges[e].V, e)
					for w := range g.Edges[e].V {
						if _, ex2 := g2.Vertices[w]; !ex2 {
							g2.AddVertex(w, 0)
							nextFrontier[w] = true
						}
					}
				}
			}
		}
		if len(nextFrontier) == 0 {
			for v := range frontier {
				g2.VertexFrontier[v] = true
			}
			break
		}
		frontier = nextFrontier
		if i == level-1 {
			for v := range frontier {
				g2.VertexFrontier[v] = true
			}
		}
	}
	g2.IncMap = g.IncMap
	g2.AdjCount = g.AdjCount

	for v := range g2.VertexFrontier {
		cond := true
		for e := range g2.IncMap[v] {
			if _, ex := g2.Edges[e]; !ex {
				cond = false
				break
			}
		}
		if cond {
			delete(g2.VertexFrontier, v)
		}
	}

	return g2
}

func ExpandFrontier(gf *HyperGraph, g *HyperGraph, level int, expand map[int32]bool) {
	oldFrontier := make(map[int32]bool)

	for v := range expand {
		oldFrontier[v] = true
	}

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range expand {
			for e := range gf.IncMap[v] {
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
			for v := range expand {
				delete(gf.VertexFrontier, v)
			}
			break
		}
		expand = nextFrontier
		if i == level-1 {
			for v := range expand {
				gf.VertexFrontier[v] = true
			}
			for v := range oldFrontier {
				delete(gf.VertexFrontier, v)
			}
		}
	}
}

func F3_ExpandFrontier(gf *HyperGraph, g *HyperGraph, remId int32, level int) {
	frontier := make(map[int32]bool)
	remEdge := g.Edges[remId]

	for v := range remEdge.V {
		delete(gf.VertexFrontier, v)
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
		for e := range gf.IncMap[v] {
			gf.F_RemoveEdge(e, g)
		}
	}

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
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
		frontier = nextFrontier
		if i == level-1 {
			for v := range frontier {
				gf.VertexFrontier[v] = true
			}
		}
	}
}
