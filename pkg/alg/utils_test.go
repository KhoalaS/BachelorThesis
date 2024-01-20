package alg

import (
	"math/rand"
	"testing"
)

func TestShuffle(t *testing.T) {
	rand.Seed(42)

	arr := make([]int, 6)
	for i:=0; i<6; i++{
		arr[i] = i
	}

	perm := []int{4,1,3,0,2,5}

	Shuffle[int](arr)

	for i:=0; i<6; i++{
		if arr[i] != perm[i] {		
			t.Errorf("Expected %d got %d", perm[i], arr[i])	
		}
	}
}