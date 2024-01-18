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

func TriangleDetection(g *HyperGraph) []map[int32]bool {
	//defer LogTime(time.Now(), "SmallTriangleRule")
	c := []map[int32]bool{}
	adjList := make(map[int32]map[int32]bool)
	exec := 0

	// Time Compelxity: |E|
	for _, e := range g.Edges {
		for v := range e.V {
			if _, ex := adjList[v]; !ex {
				adjList[v] = make(map[int32]bool)
			}
			for w := range e.V {
				if v == w {
					continue
				}
				adjList[v][w] = true
			}
		}
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
				c = append(c, remSet) 
			}
		}
	}
	return c
}
