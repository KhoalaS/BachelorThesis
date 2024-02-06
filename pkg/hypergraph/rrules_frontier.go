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

// Time Complexity: |E| * d^3

func S_EdgeDominationRule(g *HyperGraph) int {
	var wg sync.WaitGroup
	if logging {
		defer LogTime(time.Now(), "S_EdgeDomination")
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

func S_RemoveEdgeRule(g *HyperGraph, c map[int32]bool, t int) (int, bool) {
	if logging {
		defer LogTime(time.Now(), fmt.Sprintf("S_RemoveEdgeRule-%d", t))
	}

	rem := make(map[int32]bool)
	exec := 0
	remOuterLayer := false

	for eId, e := range g.Edges {
		if len(e.V) == t {
			rem[eId] = true
		}
	}

	for e := range rem {
		exec++
		for v := range g.Edges[e].V {
			c[v] = true
			for f := range g.IncMap[v] {
				if g.Edges[f].Layer == g.MaxLayer-1 {
					remOuterLayer = true
				}
				delete(rem, f)
				g.RemoveEdge(f)
			}
		}
	}
	return exec, remOuterLayer
}

func S_ApproxVertexDominationRule(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "ApproxVertexDominationRule")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
		if e.Layer == g.MaxLayer {
			continue
		}
		for v := range e.V {
			if _, ex := adjCount[v]; !ex {
				adjCount[v] = make(map[int32]int32)
			}

			for w := range e.V {
				if v != w {
					adjCount[v][w]++
				}
			}
		}
	}

	remOuterLayer := false

	// Time Complexity: |V| * (|V| + 4c)
	for solFound := true; solFound; {
		solFound = false

		for vId, count := range adjCount {
			solution, ex := twoSum(count, int32(g.Deg(vId)+1))
			if !ex {
				continue
			}

			solFound = true
			exec++

			for _, w := range solution {
				c[w] = true
				for e := range g.IncMap[w] {
					if g.Edges[e].Layer == g.MaxLayer-1 {
						remOuterLayer = true
					}
					for x := range g.Edges[e].V {
						if x == w {
							continue
						}
						//vDeg[x]--
						subEdge, _ := SetMinus(g.Edges[e], x)
						for _, y := range subEdge {
							adjCount[x][y]--
							if adjCount[x][y] == 0 {
								delete(adjCount[x], y)
							}
						}
					}
					g.RemoveEdge(e)
				}
				delete(adjCount, w)
			}
		}
	}

	return exec, remOuterLayer
}

func S_VertexDominationRule(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "VertexDominationRule")
	}
	exec := 0
	remOuterLayer := false

	for outer := true; outer; {
		outer = false
		for v, vStr := range g.Vertices {
			if vStr.Data == g.MaxLayer{
				continue
			}
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
				if vStr.Data == g.MaxLayer-1 {
					remOuterLayer = true
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
	return exec, remOuterLayer
}

// adjCount version
func S_ApproxDoubleVertexDominationRule(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "ApproxDoubleVertexDominationRule2")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0
	remOuterLayer := false

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
		if e.Layer == g.MaxLayer {
			continue
		}
		for v := range e.V {
			if _, ex := adjCount[v]; !ex {
				adjCount[v] = make(map[int32]int32)
			}

			for w := range e.V {
				if v != w {
					adjCount[v][w]++
				}
			}
		}
	}

	for outer := true; outer; {
		outer = false

		for _, e := range g.Edges {
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
				if e.Layer == g.MaxLayer -1 {
					remOuterLayer = true
				}
				solution := [2]int32{a, b}
				for _, w := range solution {
					c[w] = true
					for e := range g.IncMap[w] {
						for x := range g.Edges[e].V {
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
						g.RemoveEdge(e)
					}
				}
			}
		}
	}

	return exec, remOuterLayer
}

func S_SmallTriangleRule(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
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

	remOuterLayer := false
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
					if g.Vertices[u].Data == g.MaxLayer {
						remOuterLayer = true
					}
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

	return exec, remOuterLayer
}

func S_F3Rule(g *HyperGraph, c map[int32]bool) (int, bool) {
	s3Arr := make([]int32, len(g.Edges))

	i := 0
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			s3Arr[i] = eId
			i++
		}
	}

	remOuterLayer := false
	if i > 0 {
		r := rand.Intn(i)
		remOuterLayer = g.Edges[s3Arr[r]].Layer == g.MaxLayer
		for v := range g.Edges[s3Arr[r]].V {
			c[v] = true
			for e := range g.IncMap[v] {
				g.RemoveEdge(e)
			}
		}
	} else {
		return 0, false
	}

	return 1, remOuterLayer
}

func S_ExtendedTriangleRule(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "ExtendedTriangleRule")
	}
	exec := 0
	remOuterLayer := false

	for {
		outer := false
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
					exec++
					if f_0 == -1 {
						log.Panic("uhhh this should not happen")
					}
					if g.Edges[f_0].Layer == g.MaxLayer {
						remOuterLayer = true
					}
					for a := range g.Edges[f_0].V {
						c[a] = true
						for h := range g.IncMap[a] {
							g.RemoveEdge(h)
						}
					}

					if g.Vertices[z].Data == g.MaxLayer {
						remOuterLayer = true
					}
					c[z] = true
					for h := range g.IncMap[z] {
						g.RemoveEdge(h)
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

	return exec, remOuterLayer
}

func S_F3TargetLowDegree(g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "detectLowDegreeEdge")
	}
	closest := 1000000000
	var closestId int32 = -1
	var remEdge int32 = -1
	remOuterLayer := false

	for vId, vStr := range g.Vertices {
		if vStr.Data == g.MaxLayer {
			continue
		}
		deg := g.Deg(vId)
		if deg < closest && deg > 1 {
			closest = deg
			closestId = vId
		}
		if deg == 2 {
			found := false
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
				if g.Edges[remEdge].Layer == g.MaxLayer-1 {
					remOuterLayer = true
				}
				for v := range g.Edges[remEdge].V {
					c[v] = true
					for e := range g.IncMap[v] {
						g.RemoveEdge(e)
					}
				}
				return 1, remOuterLayer
			}
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
		return S_F3Rule(g, c)
	}

	if g.Edges[remEdge].Layer == g.MaxLayer-1 {
		remOuterLayer = true
	}
	for v := range g.Edges[remEdge].V {
		c[v] = true
		for e := range g.IncMap[v] {
			g.RemoveEdge(e)
		}
	}
	return 1, remOuterLayer
}

