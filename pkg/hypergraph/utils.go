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

	return h.Sum32()
}

func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph {
	g := NewHyperGraph()

	edgeHashes := make(map[uint32]bool)

	var tinyEdgeProb float32 = 0.01
	if !tinyEdges {
		tinyEdgeProb = 0.0
	}

	var i int32 = 0

	for ; i < n; i++ {
		g.AddVertex(i, 0)
	}

	i = 0

	for ; i < m; i++ {
		d := 1
		r := rand.Float32()
		if r > tinyEdgeProb && r < 0.6 {
			d = 2
		} else if r >= 0.6 {
			d = 3
		}

		eps := make(map[int32]bool)
		epsArr := make([]int32, d)
		for j := 0; j < d; j++ {
			val := rand.Int31n(n)
			_, ex := eps[val]
			epsArr[j] = val

			vertReroll := 0
			for ex && vertReroll < 50 {
				val = rand.Int31n(n)
				epsArr[j] = val
				_, ex = eps[val]
				vertReroll++
			}
			eps[val] = true
		}

		if len(eps) != d {
			break
		}

		hash := getHash(epsArr)
		if !edgeHashes[hash] {
			edgeHashes[hash] = true
			g.AddEdgeMap(eps)
		}
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
