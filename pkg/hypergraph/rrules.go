package hypergraph

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var Logging = false

func batchSubComp(wg *sync.WaitGroup, g *HyperGraph, subEdges map[string]bool, domEdges []int32, done chan<- map[int32]bool) {
	runtime.LockOSThread()
	defer wg.Done()

	remEdges := make(map[int32]bool)

	epArr := []int32{}

	for _, eId := range domEdges {
		for ep := range g.Edges[eId].V {
			epArr = append(epArr, ep)
		}

		// compute all size 2 subsets of edge with id eId
		subsets := list.New()
		getSubsetsRec(epArr, 2, subsets)

		for item := subsets.Front(); item != nil; item = item.Next() {
			hash := GetHash(item.Value.([]int32)...)
			if subEdges[hash] {
				remEdges[eId] = true
				break
			}
		}
		epArr = nil
	}

	done <- remEdges

	runtime.UnlockOSThread()
}

// Time Complexity: |E| * d^3

func EdgeDominationRule(g *HyperGraph) int {
	var wg sync.WaitGroup
	if Logging {
		defer LogTime(time.Now(), "EdgeDomination")
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
			g.RemoveEdge(eId)
		}
	}
	return exec
}

// Time Complexity: |E| * d

func RemoveEdgeRule(g *HyperGraph, c map[int32]bool, t int) int {
	if Logging {
		defer LogTime(time.Now(), fmt.Sprintf("RemoveEdgeRule-%d", t))
	}

	exec := 0

	for _, e := range g.Edges {
		if len(e.V) == t {
			exec++
			for v := range e.V {
				c[v] = true
				for f := range g.IncMap[v] {
					g.RemoveEdge(f)
				}
			}
		}
	}

	return exec
}

func ApproxVertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "ApproxVertexDominationRule")
	}

	exec := 0

	// Time Complexity: |V| * (|V| + 4c)
	for outer := true; outer; {
		outer = false
		for vId, count := range g.AdjCount {
			solution, ex := twoSum(count, int32(g.Deg(vId)+1))
			if !ex {
				continue
			}

			outer = true
			exec++

			for _, w := range solution {
				c[w] = true
				for e := range g.IncMap[w] {
					g.RemoveEdge(e)
				}
			}
		}
	}

	return exec
}

func VertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "VertexDominationRule")
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
				outer = true
				g.RemoveElem(v)
				exec++
			}
		}
	}

	if exec > 0 {
		g.RemoveDuplicate()
	}
	return exec
}

func ApproxDoubleVertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "S_ApproxDoubleVertexDominationRule_New")
	}

	exec := 0

	for outer := true; outer; {
		outer = false

		tsHashes := make(map[string]map[int32]bool)

		for x := range g.Vertices {
			skip := false
			twoSumAll(g.AdjCount[x], int32(g.Deg(x)), func(x0, x1 int32) {
				if skip {
					return
				}

				hash := GetHash(x0, x1)
				if _, ex := tsHashes[hash]; ex {
					for y := range tsHashes[hash] {
						if y == x {
							continue
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
							outer = true
							skip = true

							sol := [2]int32{x0, x1}

							for _, a := range sol {
								c[a] = true
								for e := range g.IncMap[a] {
									g.RemoveEdge(e)
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
		}
	}

	return exec
}

func SmallTriangleRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "SmallTriangleRule")
	}
	adjList := make(map[int32]map[int32]bool)
	rem := make(map[int32]bool)
	exec := 0

	// Time Compelxity: |E|
	for _, e := range g.Edges {
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
			g.RemoveEdge(e)
		}
	}

	return exec
}

func F3Rule(g *HyperGraph, c map[int32]bool) int {
	s3Arr := make([]int32, len(g.Edges))

	i := 0
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			s3Arr[i] = eId
			i++
		}
	}

	if i > 0 {
		r := rand.Intn(i)
		for v := range g.Edges[s3Arr[r]].V {
			c[v] = true
			for e := range g.IncMap[v] {
				g.RemoveEdge(e)
			}
		}
	} else {
		return 0
	}

	return 1
}

func SmallEdgeDegreeTwoRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		LogTime(time.Now(), "SmallEdgeDegreeTwoRule")
	}

	exec := 0

	for outer := true; outer; {
		outer = false
		for v := range g.IncMap {
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

			found := false

			found = smallDegreeTwoSub(g, c, v, s2Edge, s3Edge)

			if found {
				outer = true
				exec++
			}
		}
	}
	return exec
}

func smallDegreeTwoSub(g *HyperGraph, c map[int32]bool, vId int32, s2Edge int32, s3Edge int32) bool {
	var x int32 = -1
	var remEdge int32 = -1

	// extract x from s2Edge
	for w := range g.Edges[s2Edge].V {
		if w == vId {
			continue
		}
		x = w
	}

	found := false

	for w := range g.Edges[s3Edge].V {
		if w == vId {
			continue
		}

		for f := range g.IncMap[w] {
			if g.Edges[f].V[x] || s3Edge == f {
				continue
			} else {
				remEdge = f
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if found {
		// should be possible to delete immidietly
		c[x] = true
		for h := range g.IncMap[x] {
			g.RemoveEdge(h)
		}

		for w := range g.Edges[remEdge].V {
			c[w] = true
			for h := range g.IncMap[w] {
				g.RemoveEdge(h)
			}
		}

	}
	return found
}

func ExtendedTriangleRule(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "ExtendedTriangleRule")
	}
	exec := 0

	for outer := true; outer; {
		outer = false
		for _, e := range g.Edges {
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
					outer = true
					exec++
					if f_0 == -1 {
						log.Panic("uhhh this should not happen")
					}
					for a := range g.Edges[f_0].V {
						c[a] = true
						for h := range g.IncMap[a] {
							g.RemoveEdge(h)
						}
					}

					c[z] = true
					for h := range g.IncMap[z] {
						g.RemoveEdge(h)
					}
					break
				}
			}
		}
	}

	return exec
}

func F3TargetLowDegree(g *HyperGraph, c map[int32]bool) int {
	if Logging {
		defer LogTime(time.Now(), "detectLowDegreeEdge")
	}
	closest := 1000000000
	var closestId int32 = -1
	var remEdge int32 = -1

	for x := range g.IncMap {
		deg := g.Deg(x)
		if deg == 2 {
			found := false
			for e := range g.IncMap[x] {
				for v := range g.Edges[e].V {
					if v == x {
						continue
					}
					for f := range g.IncMap[v] {
						if f == e {
							continue
						}
						if !g.Edges[f].V[x] && len(g.Edges[f].V) == 3 {
							found = true
							remEdge = f
							break
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				for v := range g.Edges[remEdge].V {
					c[v] = true
					for e := range g.IncMap[v] {
						g.RemoveEdge(e)
					}
				}
				return 1
			}
		} else if deg < closest && deg > 1 {
			// deg > 1 not needed in alg, but will break for isolated usage
			closest = deg
			closestId = x
		}
	}

	for e := range g.IncMap[closestId] {
		found := false
		for v := range g.Edges[e].V {
			if v == closestId {
				continue
			}
			for f := range g.IncMap[v] {
				if !g.Edges[f].V[closestId] && len(g.Edges[f].V) == 3 {
					found = true
					remEdge = f
					break
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if remEdge < 0 {
		return F3Rule(g, c)
	}

	for v := range g.Edges[remEdge].V {
		c[v] = true
		for e := range g.IncMap[v] {
			g.RemoveEdge(e)
		}
	}
	return 1
}

func setToSlice[K comparable, V any](m map[K]V) []K {
	arr := make([]K, len(m))

	i := 0
	for val := range m {
		arr[i] = val
		i++
	}

	return arr
}
