package hypergraph

import (
	"container/list"
)
func getSubsetsRec(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List){
    if index ==  s{
        subset := make([]int32, s)    
        for j:= 0; j < index; j++ {
            subset[j] = (*data)[j]
        }
        subsets.PushBack(subset)
        return
    }

    if i >= n {
        return
    }

    (*data)[index] = (*arr)[i]
    
    getSubsetsRec(arr, i+1, n, s, data, index+1, subsets)
    getSubsetsRec(arr, i+1, n, s, data, index, subsets)
}

// Time Complexity: n
func twoSum(arr []IdValueHolder, t int) ([][]int32) { 
    N := int32(t)
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

// Time Complexity: n
func twoSumSingleSolution(items map[int32]int32, t int) ([]int32, bool) { 
    N := int32(t)
    lookup := make(map[int32]int32)

    for id, val := range items {
        if _, ex := lookup[N - val]; ex {
            return []int32{id, lookup[N - val]}, true
        } else {
            lookup[val] = id
        }

    }
    return nil, false
}

type TwoSumSolution struct {
    a int32
    b int32

}