package hypergraph

import (
	"container/list"
)

func getSubsetsRec(arr []int32, s int, subsets *list.List) {
	data := make([]int32, s)
	n := len(arr)
	last := s - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			data[i] = arr[j]
			if i == last {
				sub := make([]int32, s)
				for k, val := range data {
					sub[k] = int32(val)
				}
				subsets.PushBack(sub)
			} else {
				rc(i+1, j+1)
			}
		}
	}
	rc(0, 0)
}

// callback instead of list
func getSubsetsRec2(arr []int32, s int, do func(arg []int32)) {
	data := make([]int32, s)
	n := len(arr)
	last := s - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			data[i] = arr[j]
			if i == last {
				do(data)
			} else {
				rc(i+1, j+1)
			}
		}
	}
	rc(0, 0)
}

// Time Complexity: n
func twoSum(items map[int32]int32, t int32) ([]int32, bool) {
	lookup := make(map[int32]int32)

	for key, val := range items {
		if _, ex := lookup[t-val]; ex {
			return []int32{key, lookup[t-val]}, true
		} else {
			lookup[val] = key
		}
	}
	return nil, false
}

func twoSumAll(items map[int32]int32, t int32, callback func(x0 int32, x1 int32)) {
	lookup := make(map[int32]map[int32]bool)

	for key, val := range items {
		if _, ex := lookup[t-val]; ex {
			for p := range lookup[t-val] {
				callback(key, p)
			}
		}
		if _, ex := lookup[val]; !ex {
			lookup[val] = make(map[int32]bool)
		}
		lookup[val][key] = true
	}
}

func TriangleDetection(adjList map[int32]map[int32]bool) *HyperGraph {
	//defer LogTime(time.Now(), "SmallTriangleRule")
	g := NewHyperGraph()
	hashes := make(map[string]bool)
	exec := 0

	for x, val := range adjList {
		if len(val) < 2 {
			continue
		}
		arr := setToSlice(val)
		s := 2

		getSubsetsRec2(arr, s, func(subset []int32) {
			if adjList[subset[0]][subset[1]] || adjList[subset[1]][subset[0]] {
				remSet := []int32{subset[0], subset[1], x}
				hash := GetHash(remSet...)
				if !hashes[hash] {
					exec++
					g.AddEdge(remSet...)
					for _, v := range remSet {
						g.AddVertex(v, 0)
					}
					hashes[hash] = true
				}
			}
		})
	}
	return g
}

func P3Detection(g *HyperGraph) *HyperGraph {
	h := NewHyperGraph()

	edgeHashes := make(map[string]int32)
	hashes := make(map[string]bool)

	for eId, e := range g.Edges {
		edgeHashes[e.getHash()] = eId
	}

	for u := range g.Vertices {
		for v := range g.AdjCount[u] {
			for w := range g.AdjCount[v] {
				if w == u {
					continue
				}
				if _, ex := g.AdjCount[u][w]; ex {
					continue
				}
				hash := GetHash(u, v, w)
				if _, ex := hashes[hash]; !ex {
					h.AddVertex(u, 0)
					h.AddVertex(v, 0)
					h.AddVertex(w, 0)
					h.AddEdge(u, v, w)
					hashes[hash] = true
				}
			}
		}
	}

	return h
}
