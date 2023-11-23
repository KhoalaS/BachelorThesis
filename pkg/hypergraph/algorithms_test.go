package hypergraph

import (
	"container/list"
	"math/rand"
	"testing"
)

func TestTwoSum(t *testing.T) {
	val0 := IdValueHolder{Id: 0, Value: 2}
	val1 := IdValueHolder{Id: 1, Value: 4}
	val2 := IdValueHolder{Id: 2, Value: 4}
	val3 := IdValueHolder{Id: 3, Value: 6}

	arr := []IdValueHolder{val0, val1, val2, val3}

	solutions := twoSum(arr, 10)
	sol := map[int32]bool{2: true, 1: true, 3: true}
	for _, val := range solutions {
		for _, id := range val {
			if !sol[id] {
				t.Fatalf("ID %d is not part of the solution", id)
			}
		}
	}
}

func TestGetSubsetsRec(t *testing.T) {
	hashes := map[uint32]bool{276588876: true, 3284138328: true, 977105573: true}

	arr := []int32{0, 1, 2}
	subsets := list.New()
	size := 2
	data := make([]int32, size)

	getSubsetsRec(&arr, 0, len(arr), size, &data, 0, subsets)
	
	if subsets.Len() != 3 {
		t.Fatalf("Solution has size %d, expected 3.", subsets.Len())
	}
	
	for item := subsets.Front().Next(); item != nil; item = item.Next() {
		if !hashes[getHash(item.Value.([]int32))] {
			t.Fatalf("Solution %d is not a size two subset of [0,1,2].", item.Value.([]int32))
		}
	}
}

func BenchmarkTwoSum(b *testing.B) {
	size := 100000
	r := 1000
	arr := make([]IdValueHolder, size)
	for i := 0; i < size; i++ {
		arr[i] = IdValueHolder{int32(i), int32(rand.Intn(r))}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		twoSum(arr, 451)
	}
}

func BenchmarkGetSubsetsRec(b *testing.B) {
	size := 100
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
		data := make([]int32, subsetSize)
		getSubsetsRec(&arr, 0, size, subsetSize, &data, 0, subsets)
	}
}
