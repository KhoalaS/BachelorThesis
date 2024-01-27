package hypergraph

import (
	"container/list"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoSum(t *testing.T) {
	values := map[int32]int32{0: 4, 1: 3, 2: 1}

	solution, _ := twoSum(values, int32(7))
	sol := map[int32]bool{0: true, 1: true, 2: true}

	for _, val := range solution {
		assert.Equal(t, true, sol[val])
	}
}

func TestGetSubsetsRec(t *testing.T) {
	//TODO port to string hashing
	//hashes := map[uint32]bool{276588876: true, 3284138328: true, 977105573: true}

	arr := []int32{0, 1, 2}
	subsets := list.New()
	size := 2
	sol := map[string]bool{"|0|1|": true, "|0|2|": true, "|1|2|": true}

	getSubsetsRec(arr, size, subsets)

	assert.Equal(t, 3, subsets.Len())

	for item := subsets.Front(); item != nil; item = item.Next() {
		assert.Equal(t, true, sol[GetHash(item.Value.([]int32))])
	}
}

func BenchmarkTwoSum(b *testing.B) {
	size := 100000
	r := 1000
	arr := make(map[int32]int32)
	for i := 0; i < size; i++ {
		arr[int32(i)] = int32(rand.Intn(r))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		twoSum(arr, 451)
	}
}

func BenchmarkGetSubsetsRec(b *testing.B) {
	size := 1000
	arr := make([]int32, size)
	for i := 0; i < size; i++ {
		arr[i] = int32(i)
	}
	subsetSize := 2

	f, err := makeProfile("subsetRec")
	if err != nil {
		b.Fatal("Could not create profile")
	}
	defer stopProfiling(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subsets := list.New()
		getSubsetsRec(arr, subsetSize, subsets)
	}
}
