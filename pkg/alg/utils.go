package alg

import (
	"math"
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
	var num float64 = 0
	var denom float64 = 0

	for key, val := range execs {
		num += float64(Ratios[key].B * val)
		denom += float64(Ratios[key].A * val)
	}
	return num / denom
}

func GetEstOpt(execs map[string]int) int {
	opt := 0

	for key, val := range execs {
		opt += Ratios[key].A * val
	}
	return opt
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

func RoundUp(val float64, decimals int) float64{
	mul := math.Pow(10, float64(decimals))
	return math.Ceil(val * mul) / mul
}