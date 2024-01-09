package alg

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/KhoalaS/BachelorThesis/pkg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

//go:embed minedgecover.py
var scriptCode string

var Ratios = map[string]pkg.IntTuple{
	"kTiny":            {A: 1, B: 1},
	"kSmall":           {A: 2, B: 1},
	"kTri":             {A: 3, B: 2},
	"kExtTri":          {A: 4, B: 2},
	"kApVertDom":       {A: 2, B: 1},
	"kApDoubleVertDom": {A: 2, B: 1},
	"kSmallEdgeDegTwo": {A: 4, B: 2},
	"kFallback":        {A: 3, B: 1},
}

func ThreeHS_2ApprBranchOnly(g *hypergraph.HyperGraph, c map[int32]bool, K int) bool {
	//_, k := ApplyRules(g, c, K)

	if K < 0 {
		return false
	}

	return true
}

func ThreeHS_2ApprGeneral(g *hypergraph.HyperGraph, c map[int32]bool, K int, execs map[string]int) (bool, map[int32]bool, map[string]int) {
	nExecs := ApplyRules(g, c, execs, 0)

	// placeholder
	k := 0

	c_n := make(map[int32]bool)
	execs_n := make(map[string]int)

	for key, val := range nExecs {
		execs_n[key] = val
	}

	for key, val := range c {
		c_n[key] = val
	}

	if k < 0 {
		return false, c, nExecs
	}

	if len(g.Edges) > 0 {
		g_n := g.Copy()

		v, ex := PotentialTriangle(g)

		// TODO General Branching
		// This is only the Potential Triangle Situation preferred branch

		if ex {
			delete(g_n.Vertices, v)
			for _, e := range g_n.Edges {
				if e.V[v] {
					delete(e.V, v)
				}
			}
			return ThreeHS_2ApprGeneral(g_n, c_n, k, execs_n)
		} else if g.IsSimple() {
			cover := MinEdgeCover(g)
			if k-len(cover) > 0 {
				for _, w := range cover {
					c[w] = true
				}
				return true, c, execs_n
			} else {
				return false, map[int32]bool{}, map[string]int{}
			}
		} else {
			fmt.Printf("Missing Branching, %d Edges are left\n", len(g.Edges))
			return false, c_n, execs_n
		}
	}

	return true, c_n, execs_n
}

func ThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) (bool, map[int32]bool, map[string]int) {
	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, prio)
		prio = 0

		execs["kFallback"] += hypergraph.F3Prepocess(g, c, 1)
	}
	return true, c, execs
}

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

	switch prio {
	case 2:
		exec := hypergraph.SmallTriangleRule(g, c)
		execs["kTri"] += exec
	}

	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom := hypergraph.EdgeDominationRule(g)
		kVertDom := hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c, false)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		kSmallEdgeDegTwo := hypergraph.SmallEdgeDegreeTwoRule(g, c)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kExtTri := hypergraph.ExtendedTriangleRule(g, c)
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		//kExtTri := 0

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom+kSmallEdgeDegTwo+kExtTri == 0 {
			break
		}
	}

	return execs
}

func PotentialTriangle(g *hypergraph.HyperGraph) (int32, bool) {
	// e = {x, y, z}, f = {x, y, w}, g = {x, w, z}
	// f,g have to share a vertex
	// fix x
	// if there exist at least 3 edges len(e.V)=3 incident to x
	//	(for) iterate over these edges
	//		keep track of the vertices that are vertex-adjacent to y and z
	//		for an edge containing y, check if w is in the z map
	//		if true then we found a pot. Triangle Situation
	//		else add w to the y map

	incList := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		if len(e.V) != 3 {
			continue
		}

		for v := range e.V {
			if _, ex := incList[v]; !ex {
				incList[v] = make(map[int32]bool)
			}
			incList[v][eId] = true
		}
	}

	for v, incEdges := range incList {
		if len(incEdges) < 3 {
			continue
		}
		m0 := make(map[int32]bool)
		m1 := make(map[int32]map[int32]bool)

		for eId := range incEdges {
			setMinus := make([]int32, 2)

			var i int32 = 0
			for w := range g.Edges[eId].V {
				if v == w {
					continue
				}
				setMinus[i] = w
				i++
			}
			if m0[setMinus[0]] && m0[setMinus[1]] {
				for x := range m1[setMinus[0]] {
					if x == setMinus[1] {
						continue
					}
					if m1[setMinus[1]][x] {
						vInc := make([]int32, len(incList[v]))
						j := 0
						for e := range incList[v] {
							vInc[j] = e
							j++
						}
						return v, true
					}
				}
			} else {
				m0[setMinus[0]] = true
				m0[setMinus[1]] = true
				if _, ex := m1[setMinus[0]]; !ex {
					m1[setMinus[0]] = make(map[int32]bool)
				}
				if _, ex := m1[setMinus[1]]; !ex {
					m1[setMinus[1]] = make(map[int32]bool)
				}
				m1[setMinus[0]][setMinus[1]] = true
				m1[setMinus[1]][setMinus[0]] = true
			}
		}
	}
	return -1, false

}

func ParallelPotentialTriangle(g *hypergraph.HyperGraph) (int32, bool) {
	// e = {x, y, z}, f = {x, y, w}, g = {x, w, z}
	// f,g have to share a vertex
	// fix x
	// if there exist at least 3 edges len(e.V)=3 incident to x
	//	(for) iterate over these edges
	//		keep track of the vertices that are vertex-adjacent to y and z
	//		for an edge containing y, check if w is in the z map
	//		if true then we found a pot. Triangle Situation
	//		else add w to the y map

	incList := make(map[int32]map[int32]bool)
	incIndices := []int32{}

	for eId, e := range g.Edges {
		if len(e.V) != 3 {
			continue
		}

		for v := range e.V {
			if _, ex := incList[v]; !ex {
				incList[v] = make(map[int32]bool)
				incIndices = append(incIndices, v)
			}
			incList[v][eId] = true
		}
	}

	var wg sync.WaitGroup

	numCPU := runtime.NumCPU()
	lInc := len(incIndices)
	batchSize := lInc / numCPU

	if lInc < numCPU {
		numCPU = 1
		batchSize = lInc
	}

	result := make(chan int32)
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	for i := 0; i < lInc/batchSize; i++ {
		wg.Add(1)
		start := i * batchSize
		end := start + batchSize
		if lInc-end < batchSize {
			end = lInc
		}
		go findPotentialTriangle(i, ctx, &wg, incIndices[start:end], incList, g, result)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	resultValue, ok := <-result
	cancelCtx()

	if ok {
		return resultValue, true
	}

	return -1, false
}

func findPotentialTriangle(id int, ctx context.Context, wg *sync.WaitGroup,
	incIndices []int32, incList map[int32]map[int32]bool, g *hypergraph.HyperGraph,
	result chan<- int32) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer wg.Done()

	for _, v := range incIndices {
		select {
		case <-ctx.Done():
			return
		default:
			if len(incList[v]) < 3 {
				continue
			}
			m0 := make(map[int32]bool)
			m1 := make(map[int32]map[int32]bool)

			for eId := range incList[v] {
				setMinus := make([]int32, 2)

				var i int32 = 0
				for w := range g.Edges[eId].V {
					if v == w {
						continue
					}
					setMinus[i] = w
					i++
				}
				if m0[setMinus[0]] && m0[setMinus[1]] {
					for x := range m1[setMinus[0]] {
						if x == setMinus[1] {
							continue
						}
						if m1[setMinus[1]][x] {
							result <- v
							return
						}
					}
				} else {
					m0[setMinus[0]] = true
					m0[setMinus[1]] = true
					if _, ex := m1[setMinus[0]]; !ex {
						m1[setMinus[0]] = make(map[int32]bool)
					}
					if _, ex := m1[setMinus[1]]; !ex {
						m1[setMinus[1]] = make(map[int32]bool)
					}
					m1[setMinus[0]][setMinus[1]] = true
					m1[setMinus[1]][setMinus[0]] = true
				}
			}
		}
	}
}

func MinEdgeCover(g *hypergraph.HyperGraph) []int32 {
	sol := []int32{}
	incList := make(map[int32]map[int32]bool)

	for eId, e := range g.Edges {
		for v := range e.V {
			if _, ex := incList[v]; !ex {
				incList[v] = make(map[int32]bool)
			}
			incList[v][eId] = true
		}
	}

	f, err := os.CreateTemp("", "SimpleGraph_*.txt")
	if err != nil {
		log.Fatal("Could not create temp graph file:", f.Name())
	}
	defer os.Remove(f.Name())

	for v, val := range incList {
		if len(val) == 2 {

			e := []int32{}
			for eId := range incList[v] {
				e = append(e, eId)
			}
			f.WriteString(fmt.Sprintf("%d,%d,%d\n", v, e[0], e[1]))
		}
	}
	f.Close()

	script, err := os.CreateTemp("", "minedgecover_*.py")
	if err != nil {
		log.Default().Println("Could not create temp python file")
		log.Fatal(err)
	}
	defer os.Remove(script.Name())

	_, err = script.WriteString(scriptCode)
	if err != nil {
		log.Default().Println("Could not write embed python code to temp file")
		log.Fatal(err)
	}

	pyCmd := exec.Command("python3", script.Name(), f.Name())

	out, err := pyCmd.CombinedOutput()
	if err != nil {
		log.Default().Println("Could not read output")
		log.Fatal(err)
	}

	output := string(out)
	output = strings.Trim(output, "\n")
	output = output[1 : len(output)-1]
	edgesStr := strings.Split(output, ",")

	for _, e := range edgesStr {
		if len(e) == 0 {
			continue
		}
		eInt, err := strconv.Atoi(e)
		if err != nil {
			log.Fatal(err)
		}
		sol = append(sol, int32(eInt))
	}

	return sol
}
