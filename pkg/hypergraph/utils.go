package hypergraph

import (
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

func getHash(arr []int32) uint32 {
	h := xxhash.New32()
    
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := "|"

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}
	r := strings.NewReader(in)
	io.Copy(h, r)

	return h.Sum32();
}	

func GenerateTestGraph(numVertices int32, numEdges int32, tinyEdges bool) HyperGraph {
	g := NewHyperGraph()
	var tinyEdgeProb float32 = 0.01
	if !tinyEdges {
		tinyEdgeProb = 0.0
	}
	
	var i int32 = 0

	for ; i < numVertices; i++ {
		g.AddVertex(i, 0)
	}

	i = 0

	for ; i < numEdges; i++ {
		d := 1
		r := rand.Float32()
		if r > tinyEdgeProb && r < 0.6 {
			d = 2
		} else if r >= 0.6 {
			d = 3
		}
		eps := make(map[int32]bool)
		for j := 0; j < d; j++ {
			val := rand.Int31n(numVertices)
			for eps[val] {
				val = rand.Int31n(numVertices)
			}
			eps[val] = true
		}
		g.AddEdgeMap(eps)
	}

	return g
}

func SetMinus(e Edge, elem int32) ([]int32, bool) {
	arr := []int32{}
	lenBefore := len(e.v)

	for v := range e.v {
		if v == elem {
			continue
		}
		arr = append(arr, v)
	}

	lenAfter := len(arr)

	return arr, lenBefore != lenAfter
}