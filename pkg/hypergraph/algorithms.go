package hypergraph

import (
	"container/list"
)

func getSubsetsRec(arr []int32, s int, subsets *list.List) {
	data := make([]int32, s)
	getSubsetsRecMain(arr, 0, len(arr), s, data, 0, subsets)
}

func getSubsetsRecMain(arr []int32, i int, n int, s int, data []int32, index int, subsets *list.List) {
	if index == s {
		subset := make([]int32, s)
		for j := 0; j < index; j++ {
			subset[j] = data[j]
		}
		subsets.PushBack(subset)
		return
	}

	if i >= n {
		return
	}

	data[index] = arr[i]

	getSubsetsRecMain(arr, i+1, n, s, data, index+1, subsets)
	getSubsetsRecMain(arr, i+1, n, s, data, index, subsets)
}

// slightly more memory efficient version of getSubsetsRec
// uses fixed size subsets array instead of list
func getSubsetsRec2(arr []int32, i int, n int, s int, data []int32, index int, subsets *[][]int32, subsetsIndex *int) {
	if index == s {
		subset := make([]int32, s)
		for j := 0; j < index; j++ {
			subset[j] = data[j]
		}
		(*subsets)[*subsetsIndex] = subset
		*subsetsIndex++
		return
	}

	if i >= n {
		return
	}

	data[index] = arr[i]

	getSubsetsRec2(arr, i+1, n, s, data, index+1, subsets, subsetsIndex)
	getSubsetsRec2(arr, i+1, n, s, data, index, subsets, subsetsIndex)
}

// using generics and a callback to do something with the subset
func getSubsetsRec3[V any](arr []V, s int, do func(arg []V)) {
	data := make([]V, s)
	getSubsetsRec3Main[V](arr, 0, len(arr), s, data, 0, do)

}

func getSubsetsRec3Main[V any](arr []V, i int, n int, s int, data []V, index int, do func(arg []V)) {
	if index == s {
		do(data)
		return
	}

	if i >= n {
		return
	}

	data[index] = arr[i]

	getSubsetsRec3Main(arr, i+1, n, s, data, index+1, do)
	getSubsetsRec3Main(arr, i+1, n, s, data, index, do)
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