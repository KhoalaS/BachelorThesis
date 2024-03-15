package alg

import (
	"fmt"
	"math/rand"
)

func MakeExecs() map[string]int{
	execs := make(map[string]int)
	for _, k := range Labels {
		execs[k] = 0
	}
	return execs
}

func GetRatio(execs map[string]int) float64 {
	var nom float64 = 0
	var denom float64 = 0

	for key, val := range execs {
		nom += float64(Ratios[key].A * val)
		denom += float64(Ratios[key].B * val)
	}
	fmt.Println("GetRatio:", nom, denom)
	return nom / denom
}

func Shuffle[V any](arr []V){
	var t V
	for i:=len(arr)-1; i>0; i--{
		j := rand.Intn(i+1)
		t = arr[i]
		arr[i] = arr[j]
		arr[j] = t
	}
}