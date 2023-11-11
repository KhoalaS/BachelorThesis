package hypergraph

import (
	"container/list"
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