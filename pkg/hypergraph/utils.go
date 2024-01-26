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

func GetFrontierGraph(g *HyperGraph, incMap map[int32]map[int32]bool, level int, remId int32) *HyperGraph {
	g2 := NewHyperGraph()
	frontier := make(map[int32]bool)
	remEdge := g.Edges[remId]
	hashes := make(map[string]bool)

	for v := range g.Edges[remId].V {
		for e := range incMap[v] {
			for w := range g.Edges[e].V {
				if !g.Edges[remId].V[w] {
					frontier[w] = true
					g2.AddVertex(w, 0)
				}
			}
		}
	}

	for i := 0; i < level; i++ {
		nextFrontier := make(map[int32]bool)
		for v := range frontier {
			for e := range incMap[v] {
				found := true
				for w := range g.Edges[e].V {
					if remEdge.V[w] {
						found = false
						break
					}
				}
				if found {
					hash := g.Edges[e].getHash()
					if !hashes[hash] {
						hashes[hash] = true
						g2.AddEdgeMap(g.Edges[e].V)
						for w := range g.Edges[e].V {
							if !frontier[w] {
								g2.AddVertex(w, 0)
								nextFrontier[w] = true
							}
						}
					}
				}
			}
		}
		frontier = nextFrontier
	}

	return g2
}
