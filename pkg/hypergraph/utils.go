package hypergraph

import (
	"log"
	"math"
	"sort"
	"strconv"
	"time"
)

func GetHash(arr []int32) string {
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


func getRatio(execs map[string]int) float64 {
	var nom float64 = 0
	var denom float64 = 0

	for key, val := range execs {
		nom += float64(Ratios[key].A * val)
		denom += float64(Ratios[key].B * val)
	}
	return nom / denom
}

var Ratios = map[string]IntTuple{
	"kTiny":            {A: 1, B: 1},
	"kSmall":           {A: 2, B: 1},
	"kTri":             {A: 3, B: 2},
	"kExtTri":          {A: 4, B: 2},
	"kApVertDom":       {A: 2, B: 1},
	"kApDoubleVertDom": {A: 2, B: 1},
	"kSmallEdgeDegTwo": {A: 4, B: 2},
	"kFallback":        {A: 3, B: 1},
}
