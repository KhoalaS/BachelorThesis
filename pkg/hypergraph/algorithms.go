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
func twoSumOld(arr []IdValueHolder, t int) ([][]int32) { 
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
func twoSum(items map[int32]int32, t int32) ([]int32, bool) { 
    lookup := make(map[int32]int32)

    for key, val := range items {
        if _, ex := lookup[t - val]; ex {
            return []int32{key, lookup[t - val]}, true
        } else {
            lookup[val] = key
        }
    }
    return nil, false
}

type TwoSumSolution struct {
    a int32
    b int32

}