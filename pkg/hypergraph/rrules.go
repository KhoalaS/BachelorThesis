package hypergraph

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
)

// Currently the rules will output a new hypergraph struct.
// A implementation that only manipulates the partial solution C and
// computes the 'current graph' derived from C is possibly better.

// Time Complexity: |E|^2 * d

var wg sync.WaitGroup

func batchSubComp(g HyperGraph, subEdges []int32, domEdges map[int32]bool, done chan<- map[int32]bool) {
	runtime.LockOSThread()
	defer wg.Done()

	remEdges := make(map[int32]bool)

	for _, eId := range subEdges {
		for compId := range domEdges {
			if remEdges[compId] {
				continue
			}
			subset := true
			for vId := range g.Edges[eId].v {
				if !g.Edges[compId].v[vId] {
					subset = false
					break
				}
			}
			if subset {
				remEdges[compId] = true
			}
		}
	}
	done <- remEdges

	runtime.UnlockOSThread()
}

func EdgeDominationRule(g HyperGraph, c map[int32]bool) {	
	subEdges := make(map[uint32]bool)
	domEdges := make(map[int32]bool)
	remEdges := make(map[int32]bool)


	for eId, e := range g.Edges {
		if len(e.v) == 2 {
			eHash := e.getHash()
			subEdges[eHash] = true
		} else {
			domEdges[eId] = true
		} 
	}

	epArr := []int32{}


	for eId := range domEdges {
		for ep := range g.Edges[eId].v {
			epArr = append(epArr, ep)
		}

		// compute all subsets of edge with id eId
		subsets := list.New()

		for s := g.Degree-1; s > 0; s--{
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

	fmt.Println(len(remEdges))

	for eId := range remEdges {
		delete(g.Edges, eId)
	}	
	/*
		for _, e := range g.Edges {
		if len(e.v) == int(g.Degree) {
			continue
		}
		counter++
		fmt.Printf("%d/%d\r", counter, l)
		for cId, comp := range g.Edges {
			if remEdges[cId] || len(comp.v) <= len(e.v){
				continue
			}

			subset := true

			for id := range e.v {
				if !comp.v[id] {
					subset = false
					break
				}
			}

			if subset {
				remEdges[cId] = true
			}
		}
	}
	*/
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

	fmt.Println(len(remVertices))

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
// Currently a lot of overlap. 

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

