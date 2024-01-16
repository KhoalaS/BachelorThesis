package hypergraph

import (
	"container/list"
)

func getSubsetsRec(arr []int32, s int, subsets *list.List) {
	data := make([]int, s)
	n := len(arr)
	last := s - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			data[i] = j
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
			data[i] = int32(j)
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
