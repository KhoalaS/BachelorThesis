package hypergraph

import (
	"sort"
	"strconv"
)

func getHash(arr []int32) string {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := "|"

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}

	return in
	
}


func SetMinus(e Edge, elem int32) ([]int32, bool) {
	arr := []int32{}
	lenBefore := len(e.V)

	for v := range e.V {
		if v == elem {
			continue
		}
		arr = append(arr, v)
	}

	lenAfter := len(arr)

	return arr, lenBefore != lenAfter
}
