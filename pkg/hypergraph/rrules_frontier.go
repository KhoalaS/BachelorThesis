package hypergraph

import (
	"container/list"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Time Complexity: |E| * d^3

func S_EdgeDominationRule(gf *HyperGraph, g *HyperGraph, expand map[int32]bool) int {
	var wg sync.WaitGroup
	if logging {
		defer LogTime(time.Now(), "S_EdgeDomination")
	}

	subEdges := make(map[string]bool)
	domEdges := []int32{}
	exec := 0

	for eId, e := range gf.Edges {
		if len(e.V) == 2 {
			eHash := e.getHash()
			subEdges[eHash] = true
		} else {
			domEdges = append(domEdges, eId)
		}
	}

	if len(subEdges) == 0 {
		return 0
	}

	numCPU := runtime.NumCPU()
	lDom := len(domEdges)
	batchSize := lDom / numCPU

	if lDom < numCPU {
		numCPU = 1
		batchSize = lDom
	}

	channel := make(chan map[int32]bool, numCPU)

	wg.Add(numCPU)

	for i := 0; i < numCPU; i++ {
		start := i * batchSize
		end := start + batchSize
		if lDom-end < batchSize {
			end = lDom
		}
		go batchSubComp(&wg, gf, subEdges, domEdges[start:end], channel)
	}

	wg.Wait()
	close(channel)

	for msg := range channel {
		for eId := range msg {
			exec++
			for v := range g.Edges[eId].V {
				_, ex := gf.Vertices[v]
				if gf.VertexFrontier[v] || !ex {
					expand[v] = true
				}
			}
			gf.F_RemoveEdge(eId, g)
		}
	}
	return exec
}

// Time Complexity: |E| * d

func S_RemoveEdgeRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, t int, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), fmt.Sprintf("S_RemoveEdgeRule-%d", t))
	}

	rem := make(map[int32]bool)
	exec := 0

	for eId, e := range gf.Edges {
		if len(e.V) == t {
			rem[eId] = true
		}
	}

	for e := range rem {
		exec++
		for v := range g.Edges[e].V {
			c[v] = true
			for f := range g.IncMap[v] {
				for w := range g.Edges[f].V {
					_, ex := gf.Vertices[w]
					if gf.VertexFrontier[w] || !ex {
						expand[w] = true
					}
				}
				delete(rem, f)
				gf.F_RemoveEdge(f, g)
			}
		}
	}
	return exec
}

func S_ApproxVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "S_ApproxVertexDominationRule")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	// Time Complexity: |E| * d^2

	for v, inc := range g.IncMap {
		for e := range inc {
			for w := range g.Edges[e].V {
				if w == v {
					continue
				}
				if _, ex := adjCount[v]; !ex {
					adjCount[v] = make(map[int32]int32)
				}
				adjCount[v][w]++
			}
		}
	}

	// Time Complexity: |V| * (|V| + 4c)
	for solFound := true; solFound; {
		solFound = false

		for vId := range gf.Vertices {
			count := adjCount[vId]
			solution, ex := twoSum(count, int32(g.Deg(vId)+1))
			if !ex {
				continue
			}

			solFound = true
			exec++

			for _, w := range solution {
				c[w] = true
				for e := range g.IncMap[w] {
					for x := range g.Edges[e].V {
						_, ex := gf.Vertices[x]
						if gf.VertexFrontier[x] || !ex {
							expand[x] = true
						}
						if x == w {
							continue
						}
						subEdge, _ := SetMinus(g.Edges[e], x)
						for _, y := range subEdge {
							adjCount[x][y]--
							if adjCount[x][y] == 0 {
								delete(adjCount[x], y)
							}
						}
					}
					gf.F_RemoveEdge(e, g)
				}
				delete(adjCount, w)
			}
		}
	}

	return exec
}

func S_VertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "S_VertexDominationRule")
	}
	exec := 0

	for outer := true; outer; {
		outer = false
		for v := range gf.Vertices {
			vCount := make(map[int32]int)
			for e := range g.IncMap[v] {
				for w := range g.Edges[e].V {
					vCount[w]++
				}
			}
			delete(vCount, v)

			dom := false
			//var vDom int32 = -1

			for _, value := range vCount {
				if value == g.Deg(v) {
					dom = true
					//	vDom = key
					break
				}
			}

			if dom {
				_, ex := gf.Vertices[v]
				if gf.VertexFrontier[v] || !ex {
					expand[v] = true
				}
				outer = true
				gf.F_RemoveElem(v, g)
				exec++
			}
		}
	}

	if exec > 0 {
		gf.RemoveDuplicate()
	}
	return exec
}

// adjCount version
func S_ApproxDoubleVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "S_ApproxDoubleVertexDominationRule2")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	for v, inc := range g.IncMap {
		for e := range inc {
			for w := range g.Edges[e].V {
				if w == v {
					continue
				}
				if _, ex := adjCount[v]; !ex {
					adjCount[v] = make(map[int32]int32)
				}
				adjCount[v][w]++
			}
		}
	}

	for outer := true; outer; {
		outer = false

		for _, e := range gf.Edges {
			if len(e.V) != 3 {
				continue
			}

			found := false
			var a int32 = -1
			var b int32 = -1
			for u := range e.V {
				a = u

				count := make(map[int32]int)
				vd := false

				for v := range e.V {
					if v == a {
						continue
					}
					if adjCount[v][a] == int32(g.Deg(v)) {
						vd = true
						break
					}

					for w, val := range adjCount[v] {
						if e.V[w] {
							continue
						}
						if adjCount[v][a]+val == int32(g.Deg(v)) {
							count[w]++
						}
					}

				}

				if !vd {
					for v, val := range count {
						if val == 2 {
							found = true
							b = v
							break
						}
					}
				}
				if found {
					break
				}

			}

			if found {
				exec++
				solution := [2]int32{a, b}
				for _, w := range solution {
					c[w] = true
					for e := range g.IncMap[w] {
						for x := range g.Edges[e].V {
							_, ex := gf.Vertices[x]
							if gf.VertexFrontier[x] || !ex {
								expand[x] = true
							}
							if x == w {
								continue
							}
							subEdge, _ := SetMinus(g.Edges[e], x)
							for _, y := range subEdge {
								adjCount[x][y]--
								if adjCount[x][y] == 0 {
									delete(adjCount[x], y)
								}
							}
						}
						gf.F_RemoveEdge(e, g)
					}
				}
			}
		}

	}

	return exec
}

func S_SmallTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "SmallTriangleRule")
	}
	adjList := make(map[int32]map[int32]bool)
	rem := make(map[int32]bool)
	exec := 0

	// Time Compelxity: |E|
	for _, e := range gf.Edges {
		if len(e.V) != 2 {
			continue
		}
		arr := setToSlice(e.V)

		if _, ex := adjList[arr[0]]; !ex {
			adjList[arr[0]] = make(map[int32]bool)
		}
		if _, ex := adjList[arr[1]]; !ex {
			adjList[arr[1]] = make(map[int32]bool)
		}

		adjList[arr[0]][arr[1]] = true
		adjList[arr[1]][arr[0]] = true
	}

	// Time Compelxity: |V|^2

	for z, val := range adjList {
		if len(val) < 2 {
			continue
		}
		arr := setToSlice(val)
		subsets := list.New()
		getSubsetsRec(arr, 2, subsets)

		for item := subsets.Front(); item != nil; item = item.Next() {
			s := item.Value.([]int32)
			//y := subset[0] and z := subset[1]
			// triangle condition
			if adjList[s[0]][s[1]] || adjList[s[1]][s[0]] {
				exec++
				remSet := map[int32]bool{s[0]: true, s[1]: true, z: true}
				for u := range remSet {
					c[u] = true
					rem[u] = true
					for v := range adjList[u] {
						delete(adjList[v], u)
					}
					delete(adjList, u)
				}
				break
			}
		}
	}

	for v := range rem {
		for e := range g.IncMap[v] {
			for w := range g.Edges[e].V {
				_, ex := gf.Vertices[w]
				if gf.VertexFrontier[w] || !ex {
					expand[w] = true
				}
			}
			gf.F_RemoveEdge(e, g)
		}
	}

	return exec
}

func S_ExtendedTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "ExtendedTriangleRule")
	}
	exec := 0

	for {
		outer := false
		for _, e := range gf.Edges {
			found := false
			if len(e.V) != 2 {
				continue
			}
			// e has size 2
			eArr := make([]int32, 2)

			k := 0
			for ep := range e.V {
				eArr[k] = ep
				k++
			}

			for i, vert := range eArr {
				// fix y and z
				y := vert
				z := eArr[(i+1)%2]

				var f_0 int32 = -1

				// iterate over edges f incident to y
				for f := range g.IncMap[y] {
					// ensure f has size 3
					if len(g.Edges[f].V) != 3 {
						continue
					}

					// if z in f, then |e ∩ f| != 1
					if g.Edges[f].V[z] {
						continue
					}

					// iterate over edges _g incident to z
					for _g := range g.IncMap[z] {
						cond := true
						for ep := range g.Edges[_g].V {
							if ep == z {
								continue
							}
							// check if the other vertices of g are in f
							if !g.Edges[f].V[ep] {
								cond = false
								break
							}
						}

						if cond {
							f_0 = f
							found = true
							break
						}
					}

					if found {
						break
					}
				}

				if found {
					exec++

					for a := range g.Edges[f_0].V {
						c[a] = true
						for h := range g.IncMap[a] {
							for w := range g.Edges[h].V {
								_, ex := gf.Vertices[w]
								if gf.VertexFrontier[w] || !ex {
									expand[w] = true
								}
							}
							gf.F_RemoveEdge(h, g)
						}
					}

					if gf.VertexFrontier[z] {
						expand[z] = true
					}
					c[z] = true
					for h := range g.IncMap[z] {
						for w := range g.Edges[h].V {
							_, ex := gf.Vertices[w]
							if gf.VertexFrontier[w] || !ex {
								expand[w] = true
							}
						}
						gf.F_RemoveEdge(h, g)
					}
					break
				}
			}
			if found {
				outer = true
			}
		}

		if !outer {
			break
		}
	}

	return exec
}

func F3TargetLowDegreeDetect(g *HyperGraph) int32 {
	closest := 1000000000
	var closestId int32 = -1
	var remEdge int32 = -1

	for vId := range g.IncMap {
		deg := g.Deg(vId)
		if deg < closest && deg > 1 {
			closest = deg
			closestId = vId
		}
		if deg == 2 {
			for e := range g.IncMap[closestId] {
				for v := range g.Edges[e].V {
					if v == closestId {
						continue
					}
					for f := range g.IncMap[v] {
						if f == e {
							continue
						}
						if !g.Edges[f].V[closestId] && len(g.Edges[f].V) == 3 {
							return f
						}
					}

				}
			}
		}
	}

	for e := range g.IncMap[closestId] {
		for v := range g.Edges[e].V {
			if v == closestId {
				continue
			}
			for f := range g.IncMap[v] {
				if !g.Edges[f].V[closestId] && len(g.Edges[f].V) == 3 {
					return f
				}
			}
		}
	}

	if remEdge < 0 {
		return F3RuleDetect(g)
	}

	return remEdge
}

func F3RuleDetect(g *HyperGraph) int32 {
	s3Arr := make([]int32, len(g.Edges))

	i := 0
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			s3Arr[i] = eId
			i++
		}
	}

	var remEdge int32 = -1
	if i > 0 {
		r := rand.Intn(i)
		remEdge = s3Arr[r]
		return remEdge
	} else {
		return -1
	}
}
