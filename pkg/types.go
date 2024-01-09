package pkg

type IntTuple struct {
	A int
	B int
}

type Set [K comparable, V any]map[K]V

func NewSet[K comparable, V any]() Set[K,V] {
	return Set[K, V]{}
}