package hypergraph

import (
	"container/list"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

func getSubsetsRec(arr []int32, i int, end int, s int, data []int32, index int, subsets *list.List){
    if index ==  s{
        subset := make([]int32, s)    
        for j, val := range data {
            subset[j] = val
        }
        subsets.PushBack(subset)
        return
    }

    if i >= end {
        return
    }

    data[index] = arr[i]
    
    getSubsetsRec(arr, i+1, end, s, data, index+1, subsets)
    getSubsetsRec(arr, i+1, end, s, data, index, subsets)
}

func getHash(arr []int32) uint32 {
	h := xxhash.New32()
    
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := ""

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}
	r := strings.NewReader(in)
	io.Copy(h, r)

	return h.Sum32();
}	