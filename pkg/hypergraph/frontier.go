package hypergraph

import "fmt"

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

func ExpandFrontier(gf *HyperGraph, g *HyperGraph, level int) {
	frontier := make(map[int32]bool)
	oldMax := gf.MaxLayer

	for vId, v := range gf.Vertices {
		if v.Data == gf.MaxLayer+1 {
			frontier[vId] = true
		}
	}

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
			for e := range gf.IncMap[v] {
				if _, ex := gf.Edges[e]; !ex {
					gf.AddEdgeMapWLayer(g.Edges[e].V, oldMax+i+1, e)
					for w := range g.Edges[e].V {
						if _, ex2 := gf.Vertices[w]; !ex2 {
							gf.AddVertex(w, oldMax+i+2)
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
		gf.SetMaxLayer(oldMax + i + 1)
	}
}

func F3_ExpandFrontier(gf *HyperGraph, g *HyperGraph, remId int32, level int){
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

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
			for e := range gf.IncMap[v] {
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
						gf.AddEdgeMapWLayer(g.Edges[e].V, i, e)
						for w := range g.Edges[e].V {
							gf.AddVertex(w, i+1)
							if !frontier[w] {
								nextFrontier[w] = true
							}
						}
					}
				}
			}
		}
		//g2.SetMaxLayer(i)
		if len(nextFrontier) == 0 {
			break
		}
		frontier = nextFrontier
	}


}