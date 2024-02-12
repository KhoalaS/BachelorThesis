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

func S_EdgeDominationRule(gf *HyperGraph, g *HyperGraph) (int, bool) {
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
		return 0, false
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

	remOuterLayer := false
	for msg := range channel {
		for eId := range msg {
			exec++
			for v := range gf.Edges[eId].V {
				if gf.VertexFrontier[v] {
					remOuterLayer = true
					break
				}
			}
			gf.F_RemoveEdge(eId, g)
		}
	}
	return exec, remOuterLayer
}

// Time Complexity: |E| * d

func S_RemoveEdgeRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool, t int) (int, bool) {
	if logging {
		defer LogTime(time.Now(), fmt.Sprintf("S_RemoveEdgeRule-%d", t))
	}

	rem := make(map[int32]bool)
	exec := 0
	remOuterLayer := false

	for eId, e := range gf.Edges {
		if len(e.V) == t {
			rem[eId] = true
		}
	}

	for e := range rem {
		exec++
		for v := range gf.Edges[e].V {
			c[v] = true
			for f := range gf.IncMap[v] {
				for v := range g.Edges[f].V {
					if gf.VertexFrontier[v] {
						remOuterLayer = true
						break
					}
				}
				delete(rem, f)
				gf.F_RemoveEdge(f, g)
			}
		}
	}
	return exec, remOuterLayer
}

func S_ApproxVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "S_ApproxVertexDominationRule")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	// Time Complexity: |E| * d^2
	
	for _, e := range g.Edges {
		for v := range e.V {
			for w := range e.V {
				if w == v {
					continue
				}
				if _,ex := adjCount[v]; !ex {
					adjCount[v] = make(map[int32]int32)
				}
				adjCount[v][w]++
			}
		}
	}

	remOuterLayer := false

	// Time Complexity: |V| * (|V| + 4c)
	for solFound := true; solFound; {
		solFound = false

		for vId := range gf.Vertices {
			count := adjCount[vId]
			solution, ex := twoSum(count, int32(gf.Deg(vId)+1))
			if !ex {
				continue
			}

			solFound = true
			exec++

			for _, w := range solution {
				c[w] = true
				if gf.VertexFrontier[w] {
					remOuterLayer = true
				}
				for e := range gf.IncMap[w] {
					for x := range g.Edges[e].V {
						if x == w {
							continue
						}
						//vDeg[x]--
						subEdge, _ := SetMinus(gf.Edges[e], x)
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

	return exec, remOuterLayer
}

func S_VertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "S_VertexDominationRule")
	}
	exec := 0
	remOuterLayer := false

	for outer := true; outer; {
		outer = false
		for v := range gf.Vertices {
			vCount := make(map[int32]int)
			for e := range gf.IncMap[v] {
				for w := range gf.Edges[e].V {
					vCount[w]++
				}
			}
			delete(vCount, v)

			dom := false
			//var vDom int32 = -1

			for _, value := range vCount {
				if value == gf.Deg(v) {
					dom = true
					//	vDom = key
					break
				}
			}

			if dom {
				if gf.VertexFrontier[v] {
					remOuterLayer = true
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
	return exec, remOuterLayer
}

// adjCount version
func S_ApproxDoubleVertexDominationRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "S_ApproxDoubleVertexDominationRule2")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0
	remOuterLayer := false

	for v := range g.Vertices {
		for e := range g.IncMap[v] {
			for w := range g.Edges[e].V {
				if w == v {
					continue
				}
				if _,ex := adjCount[v]; !ex {
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
					if adjCount[v][a] == int32(gf.Deg(v)) {
						vd = true
						break
					}

					for w, val := range adjCount[v] {
						if e.V[w] {
							continue
						}
						if adjCount[v][a]+val == int32(gf.Deg(v)) {
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
					if gf.VertexFrontier[w] {
						remOuterLayer = true
					}
					for e := range gf.IncMap[w] {
						for x := range g.Edges[e].V {
							if x == w {
								continue
							}
							subEdge, _ := SetMinus(gf.Edges[e], x)
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

	return exec, remOuterLayer
}

func S_SmallTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool) {
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
					if gf.VertexFrontier[u] {
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
		for e := range gf.IncMap[v] {
			gf.F_RemoveEdge(e, g)
		}
	}

	return exec, remOuterLayer
}

func S_F3Rule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool, int32) {
	s3Arr := make([]int32, len(gf.Edges))

	i := 0
	for eId, e := range gf.Edges {
		if len(e.V) == 3 {
			s3Arr[i] = eId
			i++
		}
	}

	remOuterLayer := false
	var remEdge int32 = -1
	if i > 0 {
		r := rand.Intn(i)
		remEdge = s3Arr[r]
		for v := range gf.Edges[s3Arr[r]].V {
			if gf.VertexFrontier[v] {
				remOuterLayer = true
				break
			}
		}
		for v := range gf.Edges[s3Arr[r]].V {
			c[v] = true
			for e := range gf.IncMap[v] {
				gf.F_RemoveEdge(e, g)
			}
		}
	} else {
		return 0, false, -1
	}

	return 1, remOuterLayer, remEdge
}

func S_ExtendedTriangleRule(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool) {
	if logging {
		defer LogTime(time.Now(), "ExtendedTriangleRule")
	}
	exec := 0
	remOuterLayer := false

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
				for f := range gf.IncMap[y] {
					// ensure f has size 3
					if len(g.Edges[f].V) != 3 {
						continue
					}

					// if z in f, then |e ∩ f| != 1
					if g.Edges[f].V[z] {
						continue
					}

					// iterate over edges _g incident to z
					for _g := range gf.IncMap[z] {
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
					for v := range g.Edges[f_0].V {
						if gf.VertexFrontier[v] {
							remOuterLayer = true
							break
						}
					}

					for a := range g.Edges[f_0].V {
						c[a] = true
						for h := range gf.IncMap[a] {
							gf.F_RemoveEdge(h, g)
						}
					}

					if gf.VertexFrontier[z] {
						remOuterLayer = true
					}
					c[z] = true
					for h := range gf.IncMap[z] {
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

	return exec, remOuterLayer
}

func S_F3TargetLowDegree(gf *HyperGraph, g *HyperGraph, c map[int32]bool) (int, bool, int32) {
	if logging {
		defer LogTime(time.Now(), "detectLowDegreeEdge")
	}
	closest := 1000000000
	var closestId int32 = -1
	var remEdge int32 = -1
	remOuterLayer := false

	for vId := range gf.Vertices {
		if gf.VertexFrontier[vId] {
			continue
		}
		deg := gf.Deg(vId)
		if deg < closest && deg > 1 {
			closest = deg
			closestId = vId
		}
		if deg == 2 {
			found := false
			for e := range gf.IncMap[closestId] {
				for v := range gf.Edges[e].V {
					if v == closestId {
						continue
					}
					for f := range gf.IncMap[v] {
						if f == e {
							continue
						}
						if !gf.Edges[f].V[closestId] && len(gf.Edges[f].V) == 3 {
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
				for v := range gf.Edges[remEdge].V {
					if gf.VertexFrontier[v] {
						remOuterLayer = true
						break
					}
				}

				for v := range gf.Edges[remEdge].V {
					c[v] = true
					for e := range gf.IncMap[v] {
						gf.F_RemoveEdge(e, g)
					}
				}
				return 1, remOuterLayer, remEdge
			}
		}
	}

	for e := range gf.IncMap[closestId] {
		found := false
		for v := range gf.Edges[e].V {
			if v == closestId {
				continue
			}
			for f := range gf.IncMap[v] {
				if !gf.Edges[f].V[closestId] && len(gf.Edges[f].V) == 3 {
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
		return S_F3Rule(gf, g, c)
	}

	for v := range gf.Edges[remEdge].V {
		if gf.VertexFrontier[v] {
			remOuterLayer = true
			break
		}
	}

	for v := range gf.Edges[remEdge].V {
		c[v] = true
		for e := range gf.IncMap[v] {
			gf.F_RemoveEdge(e, g)
		}
	}
	return 1, remOuterLayer, remEdge
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
				return remEdge
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
