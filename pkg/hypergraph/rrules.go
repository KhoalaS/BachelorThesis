package hypergraph

import (
	"container/list"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// General TODO:
// Build an interface ontop of the HyperGraph class and inplement the "crud" there

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
			hash := getHash(item.Value.([]int32))
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
	//defer LogTime(time.Now(), "EdgeDomination")

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
	defer LogTime(time.Now(), "RemoveEdgeRule")
	rem := make(map[int32]bool)
	inc := make(map[int32]map[int32]bool)
	exec := 0
	//defer LogTime(time.Now(), fmt.Sprintf("RemoveEdgeRule-%d", t))

	for eId, e := range g.Edges {
		if len(e.V) == t {
			rem[eId] = true
		}
		for v := range e.V {
			if _, ex := inc[v]; !ex {
				inc[v] = make(map[int32]bool)
			}
			inc[v][eId] = true
		}
	}

	for e := range rem {
		exec++
		for v := range g.Edges[e].V {
			c[v] = true
			g.RemoveVertex(v)
			for f := range inc[v] {
				delete(rem, f)
				g.RemoveEdge(f)
			}
		}
	}
	return exec
}

func ApproxVertexDominationRule(g *HyperGraph, c map[int32]bool, lock bool) int {
	defer LogTime(time.Now(), "ApproxVertexDominationRule")

	vDeg := make(map[int32]int)
	adjCount := make(map[int32]map[int32]int32)
	inc := make(map[int32]map[int32]bool)

	exec := 0

	// Time Complexity: |E| * d^2
	for eId, e := range g.Edges {
		for v := range e.V {
			if _, ex := inc[v]; !ex {
				inc[v] = make(map[int32]bool)
			}
			inc[v][eId] = true

			if _, ex := adjCount[v]; !ex {
				adjCount[v] = make(map[int32]int32)
			}

			for w := range e.V {
				if v != w {
					adjCount[v][w]++
				}
			}
			vDeg[v]++
		}
	}

	// Time Complexity: |V| * (|V| + 4c)
	for solFound:=true; solFound; {
		solFound = false
		
		for vId, count := range adjCount {
			if c[vId] {
				// TODO: check if this is just a remnant of an earlier version
				// would be concerning if still needed
				fmt.Println("Uhh this should not happen")
				continue
			}

			// probably not relevant anymore
			// used to be a lock mechanism to not trigger this rule on deg 1 vertices 
			// and reserve these edges for the vertex dom rule
			if lock && vDeg[vId] == 1 {
				continue
			}

			solution, ex := twoSum(count, int32(vDeg[vId]+1))
			if !ex {
				continue
			}

			solFound = true
			exec++

			for _, w := range solution {
				c[w] = true
				for e := range inc[w] {
					for x := range g.Edges[e].V {
						if x == w {
							continue
						}
						vDeg[x]--
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
				delete(inc, w)
				delete(vDeg, w)
				delete(adjCount, w)
			}
		}
	}

	return exec
}

func VertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	//defer LogTime(time.Now(), "VertexDominationRule")

	vDeg := make(map[int32]int32)
	incList := make(map[int32]map[int32]bool)
	exec := 0

	for eId, e := range g.Edges {
		for vId := range e.V {
			vDeg[vId]++
			if _, ex := incList[vId]; !ex {
				incList[vId] = make(map[int32]bool)
			}
			incList[vId][eId] = true
		}
	}

	for outer := true; outer;{
		outer = false
		for v := range g.Vertices {
			vCount := make(map[int32]int32)
			for e := range incList[v] {
				for w := range g.Edges[e].V {
					vCount[w]++
				}
			}
			delete(vCount, v)

			dom := false
			//var vDom int32 = -1

			for _, value := range vCount {
				if value == vDeg[v] {
					dom = true
					//	vDom = key
					break
				}
			}

			if dom {
				outer = true
				//fmt.Println("Vertex ", v, " is dominated by ", vDom)
				for e := range incList[v] {
					g.RemoveElem(e, v)
				}
				g.RemoveVertex(v)
				delete(incList, v)
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
	//defer LogTime(time.Now(), "ApproxDoubleVertexDominationRule")

	incList := make(map[int32]map[int32]bool)
	exec := 0
	s3Arr := make([]int8, g.edgeCounter)

	// |E| * d
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			s3Arr[eId] = 1
		}
		for v := range e.V {
			if _, ex := incList[v]; !ex {
				incList[v] = make(map[int32]bool)
			}
			incList[v][eId] = true
		}
	}

	eArr := make([]int32, 3)

	for {
		foundSol := false
		for eId, val := range s3Arr {
			if val != 1 {
				continue
			}
			foundLocalSol := false

			var i int32 = 0
			for v := range g.Edges[int32(eId)].V {
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
					for eInc := range incList[w] {
						if incList[a][eInc] {
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

				for f := range incList[a] {
					for v := range g.Edges[f].V {
						delete(incList[v], f)
					}
					g.RemoveEdge(f)
					s3Arr[f] = 0
				}

				for f := range incList[b] {
					for v := range g.Edges[f].V {
						delete(incList[v], f)
					}
					g.RemoveEdge(f)
					s3Arr[f] = 0
				}

				c[a] = true
				c[b] = true
				delete(incList, a)
				delete(incList, b)
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
	n := len(g.Vertices)
	m := len(g.Edges)
	incMatrix := make([][]int32, n)

	for i := range incMatrix {
		incMatrix[i] = make([]int32, m)
	}

	exec := 0

	// |E| * d
	for eId, e := range g.Edges {
		for v := range e.V {
			incMatrix[v][eId] = 1
		}
	}

	//for id, row := range incMatrix {
	//	fmt.Println(id,":", row)
	//}
	//fmt.Println("--------------before")

	for {
		foundSol := false
		for _, e := range g.Edges {
			foundLocalSol := false

			if len(e.V) != 3 {
				continue
			}

			var a int32 = -1
			var b int32 = -1

			//fmt.Println("curr Edge:", e)

			for v := range e.V {
				a = v
				vCount := make(map[int32]int32)
				var xyCount int32 = 0
				for w := range e.V {
					if a == w {
						continue
					}
					for eId, eInc := range incMatrix[w] {
						if eInc == 0 {
							continue
						}

						if incMatrix[a][eId] == 1 {
							continue
						}

						for x := range g.Edges[int32(eId)].V {
							//fmt.Println("(w,a,x,eId)",w, a, x, eId)
							if !e.V[x] {
								vCount[x]++
							}
						}
						//fmt.Println("after")
						xyCount++
					}
				}
				//log.Default().Println(len(vCount))
				//fmt.Println(vCount)
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

				exec++
				c[a] = true
				c[b] = true

				//fmt.Println("removing (a,b):", a, b)

				remEdges := make(map[int]bool)
				for eId, f := range incMatrix[a] {
					if f == 1 {
						remEdges[eId] = true
						g.RemoveEdge(int32(eId))
					}
				}

				for eId, f := range incMatrix[b] {
					if f == 1 {
						remEdges[eId] = true
						g.RemoveEdge(int32(eId))
					}
				}

				for rem := range remEdges {
					for vId := range incMatrix {
						incMatrix[vId][rem] = 0
					}
				}

				for i := range incMatrix[a] {
					incMatrix[a][i] = 0
					incMatrix[b][i] = 0

				}
				foundLocalSol = false
			}
		}
		if !foundSol {
			break
		}
	}

	//for id, row := range incMatrix {
	//	fmt.Println(id,":", row)
	//}
	//fmt.Println("--------------after")

	return exec
}

func SmallTriangleRule(g *HyperGraph, c map[int32]bool) int {
	//defer LogTime(time.Now(), "SmallTriangleRule")

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
			if adjList[subset[0]][subset[1]] {
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

func F3Prepocess(g *HyperGraph, c map[int32]bool, n int) int {
	remVertices := make(map[int32]bool)

	i := 0
	for _, e := range g.Edges {
		if i == n {
			break
		}
		if len(e.V) == 3 {
			add := true
			for v := range e.V {
				if remVertices[v] {
					add = false
					break
				}
			}

			if add {
				i++
				for v := range e.V {
					remVertices[v] = true
					g.RemoveVertex(v)
					c[v] = true
				}
			}

		}
	}

	for eId, e := range g.Edges {
		for v := range e.V {
			if remVertices[v] {
				g.RemoveEdge(eId)
				break
			}
		}
	}
	return i
}

func SmallEdgeDegreeTwoRule(g *HyperGraph, c map[int32]bool) int {
	//defer LogTime(time.Now(), "SmallEdgeDegreeTwoRule")

	exec := 0
	vDeg := make(map[int32]int)
	incMap := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		for v := range e.V {
			vDeg[v]++
			if _, ex := incMap[v]; !ex {
				incMap[v] = make(map[int32]bool)
			}
			incMap[v][eId] = true
		}
	}

	for {
		outer := false
		for vId, deg := range vDeg {
			if deg != 2 {
				continue
			}

			var s2Edge int32
			var s3Edge int32
			small := 0
			found := false

			for eId := range incMap[vId] {
				l := len(g.Edges[eId].V)
				if l == 3 {
					s3Edge = eId
				} else if l == 2 {
					if small == 1 {
						s3Edge = eId
					} else {
						s2Edge = eId
					}
					small++
				}
			}

			if small == 2 {
				found = smallDegreeTwoSub(g, c, vId, s2Edge, s3Edge, incMap, vDeg)
				if found {
					exec++
					outer = true
					continue
				}
				found = smallDegreeTwoSub(g, c, vId, s3Edge, s2Edge, incMap, vDeg)
			} else if small == 1 {
				found = smallDegreeTwoSub(g, c, vId, s2Edge, s3Edge, incMap, vDeg)

			}

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

func smallDegreeTwoSub(g *HyperGraph, c map[int32]bool, vId int32, s2Edge int32, s3Edge int32, incMap map[int32]map[int32]bool, vDeg map[int32]int) bool {
	var x int32 = -1
	var remEdge int32 = -1
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

		for f := range incMap[w] {
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
		remEdges := make(map[int32]bool)
		for h := range incMap[x] {
			for b := range g.Edges[h].V {
				if vDeg[b] > 0 {
					vDeg[b]--
				}
			}
			remEdges[h] = true
		}
		c[x] = true
		delete(vDeg, x)
		delete(incMap, x)
		g.RemoveVertex(x)

		for a := range g.Edges[remEdge].V {
			for h := range incMap[a] {
				for b := range g.Edges[h].V {
					if vDeg[b] > 0 {
						vDeg[b]--
					}
				}
				remEdges[h] = true
			}
			delete(incMap, a)
			delete(vDeg, a)
			g.RemoveVertex(a)
			c[a] = true
		}

		for h := range remEdges {
			for a := range g.Edges[h].V {
				delete(incMap[a], h)
			}
			g.RemoveEdge(h)
		}

	}
	return found
}

func ExtendedTriangleRule(g *HyperGraph, c map[int32]bool) int {
	//defer LogTime(time.Now(), "ExtendedTriangleRule")

	exec := 0
	incMap := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		for v := range e.V {
			if _, ex := incMap[v]; !ex {
				incMap[v] = make(map[int32]bool)
			}
			incMap[v][eId] = true
		}
	}

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
				for f := range incMap[y] {
					// ensure f has size 3
					if len(g.Edges[f].V) != 3 {
						continue
					}

					// if z in f, then |e ∩ f| != 1
					if g.Edges[f].V[z] {
						continue
					}

					// iterate over edges _g incident to z
					for _g := range incMap[z] {
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
					remEdges := make(map[int32]bool)
					if f_0 == -1 {
						log.Panic("uhhh this should not happen")
					}
					for a := range g.Edges[f_0].V {
						for h := range incMap[a] {
							remEdges[h] = true
						}
						delete(incMap, a)
						g.RemoveVertex(a)
						c[a] = true
					}

					for h := range incMap[z] {
						remEdges[h] = true
					}
					delete(incMap, z)
					g.RemoveVertex(z)
					c[z] = true

					for h := range remEdges {
						for a := range g.Edges[h].V {
							delete(incMap[a], h)
						}
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
	//defer LogTime(time.Now(), "detectLowDegreeEdge")

	vDeg := make(map[int32]int32)
	incMap := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		for v := range e.V {
			vDeg[v]++
			if _, ex := incMap[v]; !ex {
				incMap[v] = make(map[int32]bool)
			}
			incMap[v][eId] = true
		}
	}
	var closest int32 = 1000000
	var closestId int32 = -1

	for vId, val := range vDeg {
		if val < closest {
			closest = val
			closestId = vId
		}
		if val == 2 {
			break
		}
	}

	var remEdge int32 = -1

	for e := range incMap[closestId] {
		found := false
		for v := range g.Edges[e].V {
			for f := range incMap[v] {
				if !g.Edges[f].V[closestId] {
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
		return F3Prepocess(g, c, 1)
	}

	for v := range g.Edges[remEdge].V {
		c[v] = true
		g.RemoveVertex(v)
		for e := range incMap[v] {
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

type IdValueHolder struct {
	Id    int32
	Value int32
}
