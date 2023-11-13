package hypergraph

import "testing"

func TestTwoSum(t *testing.T){
	val0 := IdValueHolder{Id: 0, Value: 2}
	val1 := IdValueHolder{Id: 1, Value: 4}
	val2 := IdValueHolder{Id: 2, Value: 4}
	val3 := IdValueHolder{Id: 3, Value: 6}

	arr := []IdValueHolder{val0, val1, val2, val3}

	solutions := twoSum(arr, 10)
	sol := map[int32]bool{2: true, 1:true, 3:true}
	for _, val := range solutions {
		for _, id := range val {
			if !sol[id] {
				t.Fatalf("ID %d is not part of the solution", id)
			}
		}
	}
}