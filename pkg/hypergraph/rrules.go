package hypergraph

import (
	"container/list"
	"log"
	"runtime"
	"sync"
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

		for s := 2; s > 0; s-- {
			data := make([]int32, s)
			getSubsetsRec(&epArr, 0, len(epArr), s, &data, 0, subsets)
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
			delete(g.Edges, eId)
		}
	}
	return exec
}

// Time Complexity: |E| * d

func RemoveEdgeRule(g *HyperGraph, c map[int32]bool, t int) int {
	remEdges := make(map[int32]bool)
	adjList := make(map[int32]map[int32]bool)
	exec := 0

	for eId, e := range g.Edges {
		if len(e.V) == t {
			remEdges[eId] = true
		}
		for v := range e.V {
			if _, ex := adjList[v]; !ex {
				adjList[v] = make(map[int32]bool)
			}
			adjList[v][eId] = true
		}
	}

	for e := range remEdges {
		exec++
		for v := range g.Edges[e].V {
			c[v] = true
			delete(g.Vertices, v)
			for f := range adjList[v] {
				if e != f {
					delete(remEdges, f)
					delete(g.Edges, f)
				}
			}
		}
		delete(g.Edges, e)
	}
	return exec
}

//Deprecated:
func ApproxVertexDominationRule(g *HyperGraph, c map[int32]bool) bool {
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	var yz Edge
	var xDom int32 = -1

	for id, edge := range g.Edges {
		if len(edge.V) < 3 {
			continue
		}

		for x := range edge.V {

			cond := true

			for idComp, edgeComp := range g.Edges {
				if id == idComp {
					continue
				}
				if edgeComp.V[x] {
					sum := 0
					for vertex := range edge.V {
						if edgeComp.V[vertex] {
							sum += 1
						}
					}
					if sum < 2 {
						cond = false
						break
					}
				}
			}
			if cond {
				xDom = x
				yz = edge
				break
			}
		}
		if xDom != -1 {
			break
		}
	}

	if xDom != -1 {
		for vertex := range yz.V {
			if vertex != xDom {
				remVertices[vertex] = true
				c[vertex] = true
				for eId, edge := range g.Edges {
					if edge.V[vertex] {
						remEdges[eId] = true
					}
				}
			}
		}
		for eId := range remEdges {
			delete(g.Edges, eId)
		}

		for vId := range remVertices {
			delete(g.Vertices, vId)
		}
		return true
	}
	return false
}

//Deprecated:
func ApproxVertexDominationRule2(g *HyperGraph, c map[int32]bool) bool {
	vSub := make(map[int32]map[string]bool)
	vSubCount := make(map[int32]map[int32]int32)
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
		for vId0 := range e.V {
			sub := []int32{}

			if _, ex := vSubCount[vId0]; !ex {
				vSubCount[vId0] = make(map[int32]int32)
				vSub[vId0] = make(map[string]bool)
			}

			for vId1 := range e.V {
				if vId0 != vId1 {
					sub = append(sub, vId1)
					vSubCount[vId0][vId1]++
				}
			}

			subHash := getHash(sub)
			vSub[vId0][subHash] = true
		}
	}

	// Time Complexity: |V| * (|V| + 4 * c)
	for vId, count := range vSubCount {
		if c[vId] {
			continue
		}
		arr := make([]IdValueHolder, len(count))
		i := 0
		for id, val := range count {
			arr[i] = IdValueHolder{Id: id, Value: val}
			i++
		}
		solutions := twoSumOld(arr, len(vSub[vId])+1)
		solFound := false

		for _, sol := range solutions {
			hash := getHash(sol)
			if vSub[vId][hash] {

				isNew := true

				for _, v := range sol {
					if c[v] {
						isNew = false
						break
					}
				}

				if !isNew {
					continue
				}

				for _, v := range sol {
					remVertices[v] = true
					c[v] = true
				}
				// delete edge here
				solFound = true
				break
			}
		}
		if solFound {
			break
		}
	}

	if len(remVertices) == 0 {
		return false
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
		delete(g.Edges, eId)
	}

	for vId := range remVertices {
		delete(g.Vertices, vId)
	}

	return true
}

func ApproxVertexDominationRule3(g *HyperGraph, c map[int32]bool, lock bool) int {
	vDeg := make(map[int32]int)
	vSubCount := make(map[int32]map[int32]int32)
	adjList := make(map[int32]map[int32]bool)

	exec := 0

	// Time Complexity: |E| * d^2
	for eId, e := range g.Edges {
		for vId0 := range e.V {
			if _, ex := adjList[vId0]; !ex {
				adjList[vId0] = make(map[int32]bool)
			}
			adjList[vId0][eId] = true

			if _, ex := vSubCount[vId0]; !ex {
				vSubCount[vId0] = make(map[int32]int32)
			}

			for vId1 := range e.V {
				if vId0 != vId1 {
					vSubCount[vId0][vId1]++
				}
			}
			vDeg[vId0]++
		}
	}

	// Time Complexity: |V| * (|V| + 4c)
	for {
		solFound := false
		for vId, count := range vSubCount {
			if c[vId] {
				continue
			}

			if lock && vDeg[vId] == 1 {
				continue
			}

			solution, ex := twoSum(count, int32(vDeg[vId]+1))
			if !ex {
				continue
			}

			solFound = true
			exec++

			for _, v := range solution {
				c[v] = true
				for remEdge := range adjList[v] {
					for w := range g.Edges[remEdge].V {
						if w == v {
							continue
						}
						subEdge, succ := SetMinus(g.Edges[remEdge], w)
						for _, u := range subEdge {
							vSubCount[w][u]--
						}
						if succ {
							vDeg[w]--
						}
					}
					delete(g.Edges, remEdge)
				}
				delete(g.Vertices, v)
				delete(adjList, v)
				delete(vDeg, v)
				delete(vSubCount, v)
			}
		}
		if !solFound {
			break
		}
	}

	return exec
}

func VertexDominationRule(g *HyperGraph, c map[int32]bool) int {
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
			//fmt.Println("Vertex ", v, " is dominated by ", vDom)
			for e := range incList[v] {
				delete(g.Edges[e].V, v)
			}
			delete(g.Vertices, v)
			delete(incList, v)
			exec++
		} else {
			//fmt.Println("Vertex ", v, " is NOT dominated")
		}
	}

	if exec > 0 {
		g.RemoveDuplicate()
	}
	return exec
}

func ApproxDoubleVertexDominationRule(g *HyperGraph, c map[int32]bool) int {
	incList := make(map[int32]map[int32]bool)
	exec := 0

	// |E| * d
	for eId, e := range g.Edges {
		for v := range e.V {
			if _, ex := incList[v]; !ex {
				incList[v] = make(map[int32]bool)
			}
			incList[v][eId] = true
		}
	}

	for {
		foundSol := false
		for _, e := range g.Edges {
			foundLocalSol := false
			if len(e.V) != 3 {
				continue
			}

			var a int32 = -1
			var b int32 = -1

			for v := range e.V {
				a = v
				vCount := make(map[int32]int32)
				var xyCount int32 = 0
				for w := range e.V {
					if a == w {
						continue
					}
					for eInc := range incList[w] {
						if incList[a][eInc] {
							continue
						}
						for x := range g.Edges[eInc].V {
							if !e.V[x] {
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
				}

				for f := range incList[b] {
					for v := range g.Edges[f].V {
						delete(incList[v], f)
					}
					g.RemoveEdge(f)
				}

				c[a] = true
				c[b] = true
				delete(incList, a)
				delete(incList, b)
				delete(g.Vertices, a)
				delete(g.Vertices, b)
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
						delete(g.Edges, int32(eId))
					}
				}

				for eId, f := range incMatrix[b] {
					if f == 1 {
						remEdges[eId] = true
						delete(g.Edges, int32(eId))
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
	adjList := make(map[int32]map[int32]bool)
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)
	exec := 0
	degThree := []int32{}

	// Time Compelxity: |E|
	for eId, e := range g.Edges {
		if len(e.V) == 3 {
			// potential setup for worsening step, but also possible outside the rule
			degThree = append(degThree, eId)
			continue
		}
		if len(e.V) != 2 {
			continue
		}
		arr := mapToSlice(e.V)

		if _, ex := adjList[arr[0]]; !ex {
			adjList[arr[0]] = make(map[int32]bool)
		}
		adjList[arr[0]][arr[1]] = true

		if _, ex := adjList[arr[1]]; !ex {
			adjList[arr[1]] = make(map[int32]bool)
		}
		adjList[arr[1]][arr[0]] = true
	}

	// Time Compelxity: |V|^2
	for x, val := range adjList {
		if len(val) < 2 {
			continue
		}
		arr := mapToSlice(val)
		subsets := list.New()
		s := 2
		data := make([]int32, s)
		getSubsetsRec(&arr, 0, len(arr), s, &data, 0, subsets)

		for item := subsets.Front(); item != nil; item = item.Next() {
			subset := item.Value.([]int32)
			//y := subset[0] and z := subset[1]
			// triangle condition
			if adjList[subset[0]][subset[1]] {
				exec++
				remSet := map[int32]bool{subset[0]: true, subset[1]: true, x: true}
				//if len(degThree) > 0 {
				//	wEdge := degThree[len(degThree)-1]
				//	removed := false
				//	for v := range g.Edges[wEdge].V {
				//		if remSet[v] {
				//			removed = true
				//			break
				//		}
				//	}
				//	if !removed {
				//		for v := range g.Edges[wEdge].V {
				//			remSet[v] = true
				//		}
				//	}
				//	degThree = degThree[:len(degThree)-1]
				//}
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
		delete(g.Edges, eId)
	}

	for vId := range remVertices {
		delete(g.Vertices, vId)
	}

	return exec
}

func F3Prepocess(g *HyperGraph, c map[int32]bool, n int) int {
	remVertices := make(map[int32]bool)

	i := 0
	for _, e := range g.Edges {
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
					c[v] = true
				}
			}

		}
		if i == n {
			break
		}
	}

	for eId, e := range g.Edges {
		for v := range e.V {
			if remVertices[v] {
				delete(g.Edges, eId)
				break
			}
		}
	}
	return n
}

func SmallEdgeDegreeTwoRule(g *HyperGraph, c map[int32]bool) int {
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
				if  l == 3 {
					s3Edge = eId
				} else if l == 2{
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
				if vDeg[b] > 0{
					vDeg[b]--
				} 
			}
			remEdges[h] = true
		}
		c[x] = true
		delete(vDeg, x)
		delete(incMap, x)
		delete(g.Vertices, x)

		for a := range g.Edges[remEdge].V {
			for h := range incMap[a] {
				for b := range g.Edges[h].V {
					if vDeg[b] > 0{
						vDeg[b]--
					} 
				}
				remEdges[h] = true
			}
			delete(incMap, a)
			delete(vDeg, a)
			delete(g.Vertices, a)
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
			for ep := range e.V{
				eArr[k] = ep
				k++
			}	
	
			for i, vert := range eArr {
				// fix y and z
				y := vert
				z := eArr[(i+1) % 2]
	
				var f_0 int32 = -1

				for f := range incMap[y] {
					if len(g.Edges[f].V) != 3 {
						continue
					}
					// f has size 3

					// Filter size e\cap f = 1
					if g.Edges[f].V[z] {
						continue
					}
					
					for _g := range incMap[z] {
						cond := true
						for ep := range g.Edges[_g].V {
							if ep == z {
								continue
							}
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
						delete(g.Vertices, a)
						c[a] = true
					}
	
					for h := range incMap[z] {
						remEdges[h] = true
					}
					delete(incMap, z)
					delete(g.Vertices, z)
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

func mapToSlice[K comparable, V any](m map[K]V) []K {
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
