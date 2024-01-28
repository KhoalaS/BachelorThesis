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

// General TODO:
// Build an interface ontop of the HyperGraph class and inplement the "crud" there

const logging = true

func batchSubComp(wg *sync.WaitGroup, g *HyperGraph, subEdges map[string]bool, domEdges []int32, done chan<- map[int32]bool) {
	runtime.LockOSThread()
	defer wg.Done()

	remEdges := make(map[int32]bool)

	epArr := []int32{}

	for _, eId := range domEdges {
		for ep := range g.Edges[eId].V {
			epArr = append(epArr, ep)
		}

		// compute all subsets of edge with id eId
		subsets := list.New()

		// TODO: only compute size 2 subsets
		for s := 2; s > 0; s-- {
			getSubsetsRec(epArr, s, subsets)
		}

		for item := subsets.Front(); item != nil; item = item.Next() {
			hash := GetHash(item.Value.([]int32))
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
	if logging {
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
	if logging {
		defer LogTime(time.Now(), "RemoveEdgeRule")
	}

	rem := make(map[int32]bool)
	exec := 0
	//defer LogTime(time.Now(), fmt.Sprintf("RemoveEdgeRule-%d", t))

	for eId, e := range g.Edges {
		if len(e.V) == t {
			rem[eId] = true
		}
	}

	for e := range rem {
		exec++
		for v := range g.Edges[e].V {
			c[v] = true
			g.RemoveVertex(v)
			for f := range g.IncMap[v] {
				delete(rem, f)
				g.RemoveEdge(f)
			}
		}
	}
	return exec
}

func ApproxVertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "ApproxVertexDominationRule")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
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
				g.RemoveVertex(w)
				delete(adjCount, w)
			}
		}
	}

	return exec
}

func ApproxVertexDominationRule2(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "ApproxVertexDominationRule")
	}

	exec := 0

	for outer := true; outer; {
		outer = false
		for _, edge := range g.Edges {
			if len(edge.V) != 3 {
				continue
			}

			found := true
			var yz []int32
			for x := range edge.V {
				yz, _ = SetMinus(edge, x)

				for f := range g.IncMap[x] {
					if !g.Edges[f].V[yz[0]] && !g.Edges[f].V[yz[1]] {
						found = false
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
				for _, v := range yz {
					c[v] = true
					for f := range g.IncMap[v] {
						g.RemoveEdge(f)
					}
				}
			}
		}
	}

	return exec
}

func VertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "VertexDominationRule")
	}
	exec := 0

	for outer := true; outer; {
		outer = false
		for v := range g.IncMap {
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
				outer = true
				g.RemoveElem(v)
				g.RemoveVertex(v)
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
	if logging {
		defer LogTime(time.Now(), "ApproxDoubleVertexDominationRule")
	}
	exec := 0
	s3Map := make(map[int32]bool)

	// |E| * d
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			s3Map[eId] = true
		}
	}

	eArr := make([]int32, 3)

	for {
		foundSol := false
		for eId, val := range s3Map {
			if !val {
				continue
			}
			foundLocalSol := false

			var i int32 = 0
			for v := range g.Edges[eId].V {
				eArr[i] = v
				i++
			}

			var a int32 = -1
			var b int32 = -1

			for i, v := range eArr {
				a = v
				vCount := make(map[int32]int32)
				var xyCount int32 = 0
				for j, w := range eArr {
					if i == j {
						continue
					}
					for eInc := range g.IncMap[w] {
						if g.IncMap[a][eInc] {
							continue
						}
						for x := range g.Edges[eInc].V {
							inBaseEdge := false
							for _, z := range eArr {
								if x == z {
									inBaseEdge = true
									break
								}
							}
							if !inBaseEdge {
								vCount[x]++
							}
						}
						xyCount++
					}
				}
				//log.Default().Println(len(vCount))
				for pb, occur := range vCount {
					if xyCount == occur {
						b = pb
						foundSol = true
						foundLocalSol = true
						break
					}
				}
				if foundLocalSol {
					break
				}
			}

			if foundLocalSol {
				foundLocalSol = false
				exec++

				for f := range g.IncMap[a] {
					g.RemoveEdge(f)
					delete(s3Map, f)
				}

				for f := range g.IncMap[b] {
					g.RemoveEdge(f)
					delete(s3Map, f)
				}

				c[a] = true
				c[b] = true
				g.RemoveVertex(a)
				g.RemoveVertex(b)
			}
		}
		if !foundSol {
			break
		}
	}
	return exec
}

func ApproxDoubleVertexDominationRule2(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "ApproxDoubleVertexDominationRule")
	}

	adjCount := make(map[int32]map[int32]int32)
	exec := 0

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
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

			var sub []int32
			found := false
			var a int32 = -1
			var b int32 = -1

			for u := range e.V {
				a = u
				sub, _ = SetMinus(e, a)
				x := sub[0]
				y := sub[1]

				t_0 := int32(g.Deg(x) - 1)
				t_1 := int32(g.Deg(y) - 1)
				t := [2]int32{t_0, t_1}

				count := make(map[int32]int)
				need := 2

				for i := 0; i < 2; i++ {
					if adjCount[sub[i]][a] == t[i] {
						need--
					} else {
						for v := range adjCount[sub[i]] {
							if v == a || v == sub[(i+1)%2] {
								continue
							}
							if adjCount[sub[i]][a]+v == t[i] {
								count[v]++
							}
						}
					}
				}

				if need == 0 {
					//dom condition met
					maxDeg := 0
					for i := 0; i < 2; i++ {
						for v := range adjCount[sub[i]] {
							if v == a || v == sub[(i+1)%2] {
								continue
							}
							if g.Deg(v) > maxDeg {
								maxDeg = g.Deg(v)
								b = v
							}
						}
					}
					found = true
					break
				} else {
					for v, val := range count {
						if val == need {
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

	return exec
}

func SmallTriangleRule(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "SmallTriangleRule")
	}
	adjList := make(map[int32]map[int32]bool)
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)
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
	for x, val := range adjList {
		if len(val) < 2 {
			continue
		}
		arr := setToSlice(val)
		subsets := list.New()
		s := 2
		getSubsetsRec(arr, s, subsets)

		for item := subsets.Front(); item != nil; item = item.Next() {
			subset := item.Value.([]int32)
			//y := subset[0] and z := subset[1]
			// triangle condition
			if adjList[subset[0]][subset[1]] || adjList[subset[1]][subset[0]] {
				exec++
				remSet := map[int32]bool{subset[0]: true, subset[1]: true, x: true}
				for y := range remSet {
					c[y] = true
					remVertices[y] = true
					for z := range adjList[y] {
						for u := range remSet {
							delete(adjList[z], u)
						}
					}
					delete(adjList, y)
				}
				break
			}
		}
	}

	for id, e := range g.Edges {
		for v := range e.V {
			if remVertices[v] {
				remEdges[id] = true
				break
			}
		}
	}

	for eId := range remEdges {
		g.RemoveEdge(eId)
	}

	for vId := range remVertices {
		g.RemoveVertex(vId)
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

// TODO: redo s2 and s3 edge selection
func SmallEdgeDegreeTwoRule(g *HyperGraph, c map[int32]bool) int {
	if logging {
		LogTime(time.Now(), "SmallEdgeDegreeTwoRule")
	}

	exec := 0

	for {
		outer := false
		for v := range g.IncMap {
			deg := g.Deg(v)
			if deg != 2 {
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

			if s3Edge+s2Edge < 0 {
				continue
			}

			found := false

			found = smallDegreeTwoSub(g, c, v, s2Edge, s3Edge)

			if found {
				outer = true
				exec++
			}
		}

		if !outer {
			break
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
		for h := range g.IncMap[x] {
			g.RemoveEdge(h)
		}

		for w := range g.Edges[remEdge].V {
			for h := range g.IncMap[w] {
				g.RemoveEdge(h)
			}
			c[w] = true
			delete(g.IncMap, w)
			g.RemoveVertex(w)
		}

		c[x] = true
		delete(g.IncMap, x)
		g.RemoveVertex(x)
	}
	return found
}

func ExtendedTriangleRule(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "ExtendedTriangleRule")
	}
	exec := 0

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

func F3TargetLowDegree(g *HyperGraph, c map[int32]bool) int {
	if logging {
		defer LogTime(time.Now(), "detectLowDegreeEdge")
	}
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
				for v := range g.Edges[remEdge].V {
					c[v] = true
					for e := range g.IncMap[v] {
						g.RemoveEdge(e)
					}
				}
				return 1
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

func F3TargetLowDegree2(g *HyperGraph, c map[int32]bool) (int, int) {
	if logging {
		defer LogTime(time.Now(), "detectLowDegreeEdge")
	}
	closest := 1000000
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
			if found {
				h := GetFrontierGraph(g, 2, remEdge)
				c_h := make(map[int32]bool)
				bestRatio := 3.0
				rule := 0
				for i := 1; i < 10; i++ {
					execs := make(map[string]int)
					applyRules(h, c_h, execs, i)
					r := getRatio(execs)
					if r < bestRatio {
						bestRatio = r
						rule = i
						fmt.Println("deg2", i, r, execs)
					}
				}

				for v := range g.Edges[remEdge].V {
					c[v] = true
					g.RemoveVertex(v)
					for e := range g.IncMap[v] {
						g.RemoveEdge(e)
					}
				}
				return 1, rule
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
		return F3Rule(g, c), 0
	}

	h := GetFrontierGraph(g, 2, remEdge)
	c_h := make(map[int32]bool)
	bestRatio := 3.0
	rule := 0
	for i := 1; i < 10; i++ {
		execs := make(map[string]int)
		applyRules(h, c_h, execs, i)
		r := getRatio(execs)
		if r < bestRatio {
			bestRatio = r
			rule = i
			fmt.Println(fmt.Sprintf("deg%d", closest), i, r, execs)
		}
	}

	for v := range g.Edges[remEdge].V {
		c[v] = true
		for e := range g.IncMap[v] {
			g.RemoveEdge(e)
		}
	}
	return 1, rule
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

type IdValueHolder struct {
	Id    int32
	Value int32
}

func applyRules(g *HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

	switch prio {
	case 1:
		exec := RemoveEdgeRule(g, c, TINY)
		execs["kTiny"] += exec
	case 2:
		exec := EdgeDominationRule(g)
		execs["kEdgeDom"] += exec
	case 3:
		exec := VertexDominationRule(g, c)
		execs["kVertDom"] += exec
	case 4:
		exec := ApproxVertexDominationRule(g, c)
		execs["kApVertDom"] += exec
	case 5:
		exec := ApproxDoubleVertexDominationRule(g, c)
		execs["kApDoubleVertDom"] += exec
	case 6:
		exec := SmallEdgeDegreeTwoRule(g, c)
		execs["kSmallEdgeDegTwo"] += exec
	case 7:
		exec := SmallTriangleRule(g, c)
		execs["kTri"] += exec
	case 8:
		exec := ExtendedTriangleRule(g, c)
		execs["kExtTri"] += exec
	case 9:
		exec := RemoveEdgeRule(g, c, SMALL)
		execs["kSmall"] += exec
	}

	kTiny := RemoveEdgeRule(g, c, TINY)
	kEdgeDom := EdgeDominationRule(g)
	kVertDom := VertexDominationRule(g, c)
	kTiny += RemoveEdgeRule(g, c, TINY)
	kApVertDom := ApproxVertexDominationRule(g, c)
	//kApDoubleVertDom := ApproxDoubleVertexDominationRule(g, c)
	kApDoubleVertDom := 0
	kSmallEdgeDegTwo := SmallEdgeDegreeTwoRule(g, c)
	kTri := SmallTriangleRule(g, c)
	kExtTri := ExtendedTriangleRule(g, c)
	kSmall := RemoveEdgeRule(g, c, SMALL)

	execs["kTiny"] += kTiny
	execs["kVertDom"] += kVertDom
	execs["kEdgeDom"] += kEdgeDom
	execs["kTri"] += kTri
	execs["kExtTri"] += kExtTri
	execs["kSmall"] += kSmall
	execs["kApVertDom"] += kApVertDom
	execs["kApDoubleVertDom"] += kApDoubleVertDom
	execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

	return execs
}
