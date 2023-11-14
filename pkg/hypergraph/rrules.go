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

		for s := g.Degree - 1; s > 0; s-- {
			getSubsetsRec(epArr, 0, len(epArr), s, make([]int32, s), 0, subsets)
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
	channel := make(chan map[int32]bool, numCPU)

	if lDom < numCPU {
		numCPU = 1
		batchSize = lDom
	}

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

func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool) {
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
	}
}

func ApproxVertexDominationRule2(g HyperGraph, c map[int32]bool) {
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

	fmt.Println(len(remVertices))

	for eId := range remEdges {
		delete(g.Edges, eId)
	}

	for vId := range remVertices {
		delete(g.Vertices, vId)
	}
}

type IdValueHolder struct {
	Id    int32
	Value int32
}
