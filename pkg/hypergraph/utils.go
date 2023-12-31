package hypergraph

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

func getHash(arr []int32) string {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := "|"

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}

	return in
	
}

func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph {
	g := NewHyperGraph()

	edgeHashes := make(map[string]bool)

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

			vertReroll := 0
			for ex && vertReroll < 50 {
				val = rand.Int31n(n)
				epsArr[j] = val
				_, ex = eps[val]
				vertReroll++
			}
			epsArr[j] = val
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

func GenerateUniformTestGraph(n int32, m int32, u int) *HyperGraph {
	g := NewHyperGraph()

	edgeHashes := make(map[string]bool)

	var i int32 = 0

	for ; i < n; i++ {
		g.AddVertex(i, 0)
	}

	i = 0

	for ; i < m; i++ {
		eps := make(map[int32]bool)
		epsArr := make([]int32, u)
		for j := 0; j < u; j++ {
			val := rand.Int31n(n)
			_, ex := eps[val]

			vertReroll := 0
			for ex && vertReroll < 1000 {
				val = rand.Int31n(n)
				epsArr[j] = val
				_, ex = eps[val]
				vertReroll++
			}
			epsArr[j] = val
			eps[val] = true
		}

		if len(eps) != u {
			continue
		}

		hash := getHash(epsArr)
		if !edgeHashes[hash] {
			edgeHashes[hash] = true
			g.AddEdgeMap(eps)
		}
	}
	return g
}

func GenerateFixDistTestGraph(n int32, m int32, dist []int) *HyperGraph {
	g := NewHyperGraph()

	sum := 0
	for _, val := range dist {
		sum += val
	}

	edgeHashes := make(map[string]bool)

	var tinyProb float32 = float32(dist[0])/float32(sum)
	var smallProb float32 = float32(dist[1])/float32(sum)

	var i int32 = 0

	for ; i < n; i++ {
		g.AddVertex(i, 0)
	}

	i = 0

	for ; i < m; i++ {
		d := 1
		r := rand.Float32()
		if r > tinyProb && r < tinyProb+smallProb {
			d = 2
		} else if r >= tinyProb+smallProb {
			d = 3
		}

		eps := make(map[int32]bool)
		epsArr := make([]int32, d)
		for j := 0; j < d; j++ {
			val := rand.Int31n(n)
			_, ex := eps[val]

			vertReroll := 0
			for ex && vertReroll < 1000 {
				val = rand.Int31n(n)
				epsArr[j] = val
				_, ex = eps[val]
				vertReroll++
			}
			epsArr[j] = val
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

func GeneratePrefAttachmentGraph(n int32, p float64, maxEdgesize int32) *HyperGraph {
	var initSize int32 = 5
	g := GenerateTestGraph(initSize, initSize, true)
	var vCounter int32 = initSize

	for vCounter < n {
		size := 2 + rand.Int31n(maxEdgesize-2+1)
		if rand.Float64() < p {
			g.AddEdgeArr(append(selectEndpoints(g ,size-1), vCounter))
			g.AddVertex(vCounter, 0)
			vCounter++
		}else{
			g.AddEdgeArr(selectEndpoints(g ,size))
		}
		fmt.Printf("%d/%d Vertices added\r", vCounter, n)
	}
	fmt.Println()


	return g
}

func selectEndpoints(g *HyperGraph, size int32) []int32 {
	endpoints := []int32{}
	ids := make([]int32, len(g.Vertices))
	i := 0
	for vId := range g.Vertices {
		ids[i] = vId
		i++
	}

	for i:=0; i<int(size); i++ {
		pSum := make([]int32, len(ids))
		pSum[0] = g.VDeg[ids[0]]
		
		for j:=1; j<len(ids); j++{
			pSum[j] = pSum[j-1] + g.VDeg[ids[j]] 
		}

		r := rand.Int31n(pSum[len(pSum)-1]+1)

		for k:=0; k<len(pSum); k++ {
			if r <= pSum[k] {
				endpoints = append(endpoints, int32(k))
				ids = append(ids[:k], ids[k+1:]...)
				break
			}
		}
	}

	return endpoints
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
