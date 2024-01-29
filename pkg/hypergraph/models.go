package hypergraph

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

func TestGraph(n int32, m int32, tinyEdges bool) *HyperGraph {
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

		hash := GetHash(epsArr)
		if !edgeHashes[hash] {
			edgeHashes[hash] = true
			g.AddEdgeMap(eps)
		}
	}

	return g
}

func UniformTestGraph(n int32, m int32, u int) *HyperGraph {
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

		hash := GetHash(epsArr)
		if !edgeHashes[hash] {
			edgeHashes[hash] = true
			g.AddEdgeMap(eps)
		}
	}
	return g
}

func FixDistTestGraph(n int32, m int32, dist []int) *HyperGraph {
	g := NewHyperGraph()

	sum := 0
	for _, val := range dist {
		sum += val
	}

	edgeHashes := make(map[string]bool)

	var tinyProb float32 = float32(dist[0]) / float32(sum)
	var smallProb float32 = float32(dist[1]) / float32(sum)

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

		hash := GetHash(epsArr)
		if !edgeHashes[hash] {
			edgeHashes[hash] = true
			g.AddEdgeMap(eps)
		}
	}

	return g
}

func PrefAttachmentGraph(n int32, p float64, maxEdgesize int32) *HyperGraph {
	var initSize int32 = 5
	g := TestGraph(initSize, initSize, true)
	var vCounter int32 = initSize

	for vCounter < n {
		size := 2 + rand.Int31n(maxEdgesize-2+1)
		if rand.Float64() < p {
			g.AddEdge(append(selectEndpoints(g, size-1), vCounter)...)
			g.AddVertex(vCounter, 0)
			vCounter++
		} else {
			g.AddEdge(selectEndpoints(g, size)...)
		}
		fmt.Printf("%d/%d Vertices added\r", vCounter, n)
	}
	fmt.Println()

	return g
}

func ModPrefAttachmentGraph(n int, r int, p float64, alpha float64) *HyperGraph {
	g := NewHyperGraph()
	c := make([][]int32, r)
	pcum := generate3PMatrix(r, alpha)
	vc := 0

	for ; vc < r; vc++ {
		c[vc] = []int32{}
		g.AddVertex(int32(vc), vc)
		g.AddEdge(int32(vc))
		c[vc] = append(c[vc], int32(vc))
	}

	for vc < n {
		//fmt.Printf("%d/%d\r", vc, n)
		if rand.Float64() < p {
			roll := rand.Intn(r)
			g.AddVertex(int32(vc), roll)
			g.AddEdge(int32(vc))
			c[roll] = append(c[roll], int32(vc))
			vc++
		} else {
			roll := rand.Float64()
			found := false

			cArr := make([]int, 3)

			for i := 0; i < r; i++ {
				for j := 0; j < r; j++ {
					for k := 0; k < r; k++ {
						if roll <= pcum[i][j][k] {
							found = true
							cArr[0] = i
							cArr[1] = j
							cArr[2] = k
							break
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}

			if !found {
				log.Panic("This should not be possible")
			}

			// we ignore the X distribution, resulting in a 3-uniform hypergraph
			eps := make(map[int32]bool)

			for _, cId := range cArr {
				dSum := 0
				for _, v := range c[cId] {
					dSum += int(g.Deg(v))
				}

				roll := rand.Intn(dSum)
				dCounter := 0

				for _, v := range c[cId] {
					dCounter += int(g.Deg(v))
					if roll <= dCounter {
						eps[v] = true
						break
					}
				}
			}
			g.AddEdgeMap(eps)
		}
	}

	return g
}

func UniformERGraph(n int, p float64, evr float64, size int) *HyperGraph {
	g := NewHyperGraph()
	nArr := make([]int32, n)

	for i := 0; i < n; i++ {
		g.AddVertex(int32(i), 0)
		nArr[i] = int32(i)
	}

	if evr > 0 {
		p = float64(n) * evr / float64(binomialCoefficient(n, size))
	}

	// Dont actually compute all of them but compute them one at at time

	getSubsetsRec2(nArr, size, func(arg []int32) {
		if rand.Float64() < p {
			g.AddEdge(arg...)
		}
	})

	return g
}

func UniformERGraphCallback(n int, p float64, evr float64, size int, callback func(edge []int32)) {
	nArr := make([]int32, n)

	for i := 0; i < n; i++ {
		nArr[i] = int32(i)
	}

	if evr > 0 {
		p = float64(n) * evr / float64(binomialCoefficient(n, size))
	}

	// Dont actually compute all of them but compute them one at at time

	getSubsetsRec2(nArr, size, func(arg []int32) {
		if rand.Float64() < p {
			callback(arg)
		}
	})
}

func generate3PMatrix(r int, alpha float64) [][][]float64 {
	p := make([][][]float64, r)

	for i := 0; i < r; i++ {
		p[i] = make([][]float64, r)
		for j := 0; j < r; j++ {
			p[i][j] = make([]float64, r)
		}
	}

	single := (1.0 - alpha) / float64(r)
	u := alpha / (math.Pow(float64(r), 3) - float64(r))

	p[0][0][0] = single

	for i := 0; i < r; i++ {
		for j := 0; j < r; j++ {
			for k := 0; k < r; k++ {
				if i == j && j == k {
					p[i][j][k] = single
				} else {
					p[i][j][k] = u
				}
			}
		}
	}

	for i := 0; i < r; i++ {
		for j := 0; j < r; j++ {
			for k := 0; k < r; k++ {
				if k == 0 && j == 0 && i == 0 {
					continue
				} else if k == 0 && j == 0 {
					p[i][j][k] = p[i-1][r-1][r-1] + p[i][j][k]
				} else if k == 0 {
					p[i][j][k] = p[i][j-1][r-1] + p[i][j][k]
				} else {
					p[i][j][k] = p[i][j][k-1] + p[i][j][k]
				}
			}
		}
	}

	return p
}

func selectEndpoints(g *HyperGraph, size int32) []int32 {
	endpoints := []int32{}
	ids := make([]int32, len(g.Vertices))
	i := 0
	for vId := range g.Vertices {
		ids[i] = vId
		i++
	}

	for i := 0; i < int(size); i++ {
		pSum := make([]int, len(ids))
		pSum[0] = g.Deg(ids[0])

		// recalculate the cumulative probability array
		for j := 1; j < len(ids); j++ {
			pSum[j] = pSum[j-1] + g.Deg(ids[j])
		}

		r := rand.Intn(pSum[len(pSum)-1])

		for k := 0; k < len(pSum); k++ {
			if r <= pSum[k] {
				endpoints = append(endpoints, int32(k))
				ids = append(ids[:k], ids[k+1:]...)
				break
			}
		}
	}

	return endpoints
}
