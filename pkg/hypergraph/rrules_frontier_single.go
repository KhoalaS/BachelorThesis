package hypergraph

import (
	"container/list"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var logHistory = true

func FS_EdgeDominationRule(g *HyperGraph, expand map[int32]bool) int {
	var wg sync.WaitGroup
	if Logging {
		defer LogTime(time.Now(), "FS_EdgeDomination")
	}

	subEdges := make(map[string]bool)
	domEdges := []int32{}
	exec := 0

	for eId, e := range g.Edges {
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
		go batchSubComp(&wg, g, subEdges, domEdges[start:end], channel)
	}

	wg.Wait()
	close(channel)

	for msg := range channel {
		for eId := range msg {
			exec++
			for v := range g.Edges[eId].V {
				expand[v] = true
			}
			g.RemoveEdge(eId)
		}
	}

	if logHistory && exec > 0 {
		log.Default().Println("EDom")
	}

	return exec
}

// Time Complexity: |E| * d

func FS_TinyEdgeRule(g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_TinyEdgeRule")
	}

	exec := 0

	for _, e := range g.Edges {
		if len(e.V) == TINY {
			exec++
			for v := range e.V {
				c[v] = true
				for f := range g.IncMap[v] {
					for w := range g.Edges[f].V {
						expand[w] = true
					}
					g.RemoveEdge(f)
				}
			}
		}
	}
	
	if logHistory && exec > 0 {
		log.Default().Println("Tiny")
	}

	return exec
}

func FS_VertexDominationRule(g *HyperGraph, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_VertexDominationRule")
	}
	exec := 0

	for outer := true; outer; {
		outer = false
		for v := range g.Vertices {
			dom := false
			for _, value := range g.AdjCount[v] {
				if int(value) == g.Deg(v) {
					dom = true
					break
				}
			}

			if dom {
				for f := range g.IncMap[v] {
					for w := range g.Edges[f].V {
						expand[w] = true
					}
				}
				outer = true
				g.RemoveElem(v)
				exec++
			}
		}
	}

	if exec > 0 {
		g.RemoveDuplicate()
	}

	if logHistory && exec > 0 {
		log.Default().Println("VDom")
	}

	return exec
}

// adjusted
func FS_RemoveEdgeRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, t int, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), fmt.Sprintf("FS_RemoveEdgeRule-%d", t))
	}

	exec := 0

	for _, e := range gf.Edges {
		if len(e.V) == t {
			exec++
			for v := range e.V {
				c[v] = true
				for f := range g.IncMap[v] {
					for w := range g.Edges[f].V {
						expand[w] = true
					}
					gf.F_RemoveEdge(f, g)
				}
			}
			break
		}
	}

	if logHistory && exec > 0 {
		log.Default().Println("RemEdge-", t)
	}

	return exec
}

// adjusted
func FS_ApproxVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_ApproxVertexDominationRule")
	}

	exec := 0

	for vId := range gf.Vertices {
		count := g.AdjCount[vId]
		solution, ex := twoSum(count, int32(g.Deg(vId)+1))

		if !ex {
			continue
		}

		exec++

		for _, w := range solution {
			c[w] = true
			for e := range g.IncMap[w] {
				for x := range g.Edges[e].V {
					expand[x] = true
				}
				gf.F_RemoveEdge(e, g)
			}
		}
		break
	}

	if logHistory && exec > 0 {
		log.Default().Println("APVD")
	}

	return exec
}

// adjusted
func FS_ApproxDoubleVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_ApproxDoubleVertexDominationRule")
	}

	exec := 0

	tsHashes := make(map[string]map[int32]bool)

	for x := range gf.Vertices {
		skip := false
		twoSumAll(g.AdjCount[x], int32(g.Deg(x)), func(x0, x1 int32) {
			if skip {
				return
			}

			hash := GetHash(x0, x1)
			if _, ex := tsHashes[hash]; ex {
				for y := range tsHashes[hash] {
					if y == x {
						return
					}
					h1 := GetHash(x, y, x0)
					h2 := GetHash(x, y, x1)

					found := false

					for e := range g.IncMap[y] {
						he := g.Edges[e].getHash()
						if he == h1 || he == h2 {
							found = true
							break
						}
					}

					if found {
						exec++
						skip = true

						sol := [2]int32{x0, x1}

						for _, a := range sol {
							c[a] = true
							for e := range g.IncMap[a] {
								for w := range g.Edges[e].V {
									expand[w] = true
								}
								gf.F_RemoveEdge(e, g)
							}
						}
						return
					}
					tsHashes[hash][x] = true
				}
			} else {
				tsHashes[hash] = make(map[int32]bool)
				tsHashes[hash][x] = true
			}
		})
		if skip {
			break
		}
	}

	if logHistory && exec > 0 {
		log.Default().Println("APDVD")
	}

	return exec
}

// adjusted
func FS_SmallTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_SmallTriangleRule")
	}
	adjList := make(map[int32]map[int32]bool)
	rem := make(map[int32]bool)
	exec := 0

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

		if exec > 0 {
			break
		}
	}

	for v := range rem {
		for e := range g.IncMap[v] {
			for w := range g.Edges[e].V {
				expand[w] = true
			}
			gf.F_RemoveEdge(e, g)
		}
	}

	if logHistory && exec > 0 {
		log.Default().Println("Tri")
	}

	return exec
}

// adjusted
func FS_ExtendedTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "FS_ExtendedTriangleRule")
	}
	exec := 0

Outer:
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
							expand[w] = true
						}
						gf.F_RemoveEdge(h, g)
					}
				}

				c[z] = true
				for h := range g.IncMap[z] {
					for w := range g.Edges[h].V {
						expand[w] = true
					}
					gf.F_RemoveEdge(h, g)
				}
				break Outer
			}
		}
	}

	if logHistory && exec > 0 {
		log.Default().Println("ETri")
	}

	return exec
}

// adjusted
func FS_SmallEdgeDegreeTwoRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, expand map[int32]bool) (int, int) {
	if Logging {
		LogTime(time.Now(), "FS_SmallEdgeDegreeTwoRule")
	}

	exec0 := 0
	exec1 := 0

	for v := range gf.Vertices {
		if g.Deg(v) != 2 {
			continue
		}
		// assert that deg(v) = 2

		var s2Edge int32 = -1
		var s3Edge int32 = -1

		for e := range g.IncMap[v] {
			l := len(g.Edges[e].V)
			if l == 3 {
				s3Edge = e
			} else if l == 2 {
				s2Edge = e
			}
		}

		if s3Edge == -1 || s2Edge == -1 {
			continue
		}

		found, remSize := S_smallDegreeTwoSub(gf, g, c, v, s2Edge, s3Edge, expand)

		if found {
			if remSize == 3 {
				exec0++
			} else {
				exec1++
			}
			break
		}
	}

	if logHistory && exec0 + exec1 > 0 {
		log.Default().Println("SED2")
	}

	return exec0, exec1
}