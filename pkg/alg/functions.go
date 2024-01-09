package alg

import "math"

func Exp(x int) int {
	return int(math.Pow(2, float64(x)))
}

func Linear(x int) int {
	return x
}

func Sqrt(x int) int {
	return int(math.Floor(math.Sqrt(float64(x))))
}

func Const1() int {
	return 1
}