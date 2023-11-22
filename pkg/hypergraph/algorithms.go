package hypergraph

import (
	"container/list"
)

func getSubsetsRec(arr *[]int32, i int, end int, s int, data *[]int32, index int, subsets *list.List){
    if index ==  s{
        subset := make([]int32, s)    
        for j, val := range *data {
            subset[j] = val
        }
        subsets.PushBack(subset)
        return
    }

    if i >= end {
        return
    }

    (*data)[index] = (*arr)[i]
    
    getSubsetsRec(arr, i+1, end, s, data, index+1, subsets)
    getSubsetsRec(arr, i+1, end, s, data, index, subsets)
}

func twoSum(arr []IdValueHolder, n int) ([][]int32) { 
    N := int32(n)
    lookup := make(map[int32][]IdValueHolder)
    solutions := [][]int32{}

    for _, val := range arr {
        if _, ex := lookup[N - val.Value]; ex {
            for _, comp := range lookup[N - val.Value] {
                solutions = append(solutions, []int32{val.Id, comp.Id})
            }
        } else {
            if _, ex := lookup[val.Value]; !ex {
                lookup[val.Value] = []IdValueHolder{}
            }
            lookup[val.Value] = append(lookup[val.Value], val)
        }

    }
    return solutions
}