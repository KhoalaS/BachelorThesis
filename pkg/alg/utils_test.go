package alg

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, arr[i], perm[i])
	}
}