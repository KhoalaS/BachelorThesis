package hypergraph

import (
	"container/list"
	"testing"
)

func TestGetSubsetsRec(t *testing.T){
	arr := []int32{1,2,3}
	subsets := list.New()
	s := 2
	getSubsetsRec(arr,0,3,s,make([]int32, s), 0, subsets)
	for e := subsets.Front(); e != nil; e = e.Next() {
        t.Log(e.Value)
    }
}