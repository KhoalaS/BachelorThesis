package hypergraph

import (
	"log"
	"math"
	"sort"
	"strconv"
	"time"
)

func GetHash(arr ...int32) string {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := "|"

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}

	return in

}

func SetMinus(e Edge, elem int32) ([]int32, bool) {
	arr := []int32{}
	lenBefore := len(e.V)

	for v := range e.V {
		if v == elem {
			continue
		}
		arr = append(arr, v)
	}

	lenAfter := len(arr)

	return arr, lenBefore != lenAfter
}

func LogTime(start time.Time, name string) {
	stop := time.Since(start)
	log.Printf("%s took %s\n", name, stop)
}

func binomialCoefficient(n int, k int) int {
	//wenn 2*k > n dann k = n-k
	//ergebnis = 1
	//für i = 1 bis k
	//    ergebnis = ergebnis * (n + 1 - i) / i
	//rückgabe ergebnis
	if 2*k > n {
		k = n - k
	}
	c := 1.0
	for i := 1; i <= k; i++ {
		c = c * float64(n+1-i) / float64(i)
	}
	return int(math.Ceil(c))
}
