package hypergraph

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
)

func batchSubComp(wg *sync.WaitGroup, g HyperGraph, subEdges map[uint32]bool, domEdges []int32, done chan<- map[int32]bool) {
	runtime.LockOSThread()
	defer wg.Done()

	remEdges := make(map[int32]bool)

	epArr := []int32{}

	for _, eId := range domEdges {
		for ep := range g.Edges[eId].v {
			epArr = append(epArr, ep)
		}

		// compute all subsets of edge with id eId
		subsets := list.New()
		
		for s := 2; s > 0; s-- {
			data := make([]int32, s);
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

func EdgeDominationRule(g HyperGraph, c map[int32]bool) {
	var wg sync.WaitGroup

	subEdges := make(map[uint32]bool)
	domEdges := []int32{}

	for eId, e := range g.Edges {
		if len(e.v) == 2 {
			eHash := e.getHash()
			subEdges[eHash] = true
		} else {
			domEdges = append(domEdges, eId)
		}
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
			delete(g.Edges, eId)
		}
	}

}

// Time Complexity: |E| * d

func RemoveEdgeRule(g HyperGraph, c map[int32]bool, t int) {
	remEdges := make(map[int32]bool)
	remVertices := make(map[int32]bool)

	for _, e := range g.Edges {
		if len(e.v) == t {
			for v := range e.v {
				remVertices[v] = true
				c[v] = true
			}
		}
	}

	for id, e := range g.Edges {
		for v := range e.v {
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
	fmt.Println("delete finished")
}

// Complexity: (|E| * d)^2
// What can be done fast:
// - Extract size-3 edges
// - compute subsets
// -

// New Algorithm
// Iterate over all edges, extracting all edges degree 3
// associate every vertex with its other vertices in an edge
// vSub := map[int32]map[int32][]int32
// sum the occurences of these size two sets up in a map
// if the map contains at least two values a,b such that a+b=len(vSub)+1 then we know, that there
// exists two vertices that are part of every edge that x is an element of
// Complexity:
// |E| + |E| + d

func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool) bool {
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	var yz Edge
	var xDom int32 = -1

	for id, edge := range g.Edges {
		if len(edge.v) < 3 {
			continue
		}

		for x := range edge.v {

			cond := true

			for idComp, edgeComp := range g.Edges {
				if id == idComp {
					continue
				}
				if edgeComp.v[x] {
					sum := 0
					for vertex := range edge.v {
						if edgeComp.v[vertex] {
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
		for vertex := range yz.v {
			if vertex != xDom {
				remVertices[vertex] = true
				c[vertex] = true
				for eId, edge := range g.Edges {
					if edge.v[vertex] {
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

func ApproxVertexDominationRule2(g HyperGraph, c map[int32]bool) bool {
	vSub := make(map[int32]map[uint32]bool)
	vSubCount := make(map[int32]map[int32]int32)
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	// Time Complexity: |E| * d^2
	for _, e := range g.Edges {
		for vId0 := range e.v {
			sub := []int32{}

			if _, ex := vSubCount[vId0]; !ex {
				vSubCount[vId0] = make(map[int32]int32)
				vSub[vId0] = make(map[uint32]bool)
			}

			for vId1 := range e.v {
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
		solutions := twoSum(arr, len(vSub[vId])+1)
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
		for v := range e.v {
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

func ApproxVertexDominationRule3(g HyperGraph, c map[int32]bool) {
	vSub := make(map[int32]map[uint32]bool)
	vSubCount := make(map[int32]map[int32]int32)
	remVertices := make(map[int32]bool)
	adjList := make(map[int32]map[int32]bool)

	// Time Complexity: |E| * d^2
	for eId, e := range g.Edges {
		for vId0 := range e.v {
			sub := []int32{}
			if _, ex := adjList[vId0]; !ex {
				adjList[vId0] = make(map[int32]bool)
			}
			adjList[vId0][eId] = true

			if _, ex := vSubCount[vId0]; !ex {
				vSubCount[vId0] = make(map[int32]int32)
				vSub[vId0] = make(map[uint32]bool)
			}

			for vId1 := range e.v {
				if vId0 != vId1 {
					sub = append(sub, vId1)
					vSubCount[vId0][vId1]++
				}
			}

			subHash := getHash(sub)
			vSub[vId0][subHash] = true
		}
	}


	// Time Complexity: |V| * (|V| + 4c)
	for ; true; {
		solFound := false
		for vId, count := range vSubCount {
			if c[vId] {
				continue
			}
			
			arr := make([]IdValueHolder, len(count))
			i := 0
			for id, val := range count {
				if val == 0 {
					continue
				}
				arr[i] = IdValueHolder{Id: id, Value: val}
				i++
			}
			arr = arr[0:i]

			target := 0
			for _, m := range vSub[vId] {
				if m {
					target++
				}
			}
			solutions := twoSum(arr, target+1)
	
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
						c[v] = true
						remVertices[v] = true
						for remEdge := range adjList[v] {
							for w := range g.Edges[remEdge].v {
								if w == v {
									continue
								}
								subEdge, succ := SetMinus(g.Edges[remEdge], w)
								for _, u := range subEdge {
									vSubCount[w][u]--
								}
								if succ {
									vSub[w][getHash(subEdge)] = false
								}
							}
							delete(g.Edges, remEdge)
						}
						delete(adjList, v)
						delete(vSub, v)
						delete(vSubCount, v)
					}
					solFound = true
					break
				}
			}
		}
		if !solFound {
			break
		}
	}

	for vId := range remVertices {
		delete(g.Vertices, vId)
	}
}

func SmallTriangleRule(g HyperGraph, c map[int32]bool) {
	adjList := make(map[int32]map[int32]bool)
	remVertices := make(map[int32]bool)
	remEdges := make(map[int32]bool)

	// Time Compelxity: |E|
	for _, e := range g.Edges {
		if len(e.v) != 2 {
			continue
		}
		arr := mapToSlice(e.v)

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
				remSet := []int32{subset[0], subset[1], x}
				for _, y := range remSet {
					c[y] = true
					remVertices[y] = true
					for z := range adjList[y] {
						for _, u := range remSet {
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
		for v := range e.v {
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
}

func mapToSlice[K comparable, V any ](m map[K]V) []K {
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
