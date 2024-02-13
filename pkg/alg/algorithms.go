package alg

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

//go:embed minedgecover.py
var scriptCode string

var Labels = []string{"kTiny", "kVertDom", "kEdgeDom", "kSmall", "kTri", "kExtTri", "kApVertDom", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kFallback"}

var Ratios = map[string]IntTuple{
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

func LoggingThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, graphtype string, masterfilename string, iteration int, outdir string) map[string]int {

	header := "Ratio;"
	header += strings.Join(Labels, ";") + "\n"

	os.Mkdir(outdir, 0700)

	logfilename := fmt.Sprintf("%s/%s_%.2f_%d.csv", outdir, graphtype, float64(len(g.Edges))/float64(len(g.Vertices)), iteration)
	logfile, err := os.Create(logfilename)
	if err != nil {
		log.Fatalf("Could not create file %s", logfilename)
	}
	logWriter := bufio.NewWriter(logfile)

	fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
	masterfile, err := os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			masterfile, _ = os.Create(fMasterFilename)
			masterfile.WriteString(header)
		} else {
			log.Fatalf("Could not open file %s", fMasterFilename)
		}
	}

	defer logfile.Close()
	defer masterfile.Close()

	logWriter.WriteString(header)

	execs := MakeExecs()
	msg := ""

	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, 0)
		execs["kFallback"] += hypergraph.F3TargetLowDegree(g, c)

		msg = fmt.Sprintf("%f;", GetRatio(execs))
		for _, v := range Labels {
			msg += fmt.Sprintf("%d;", execs[v])
		}
		msg = msg[:len(msg)-1] + "\n"
		logWriter.WriteString(msg)
	}
	masterfile.WriteString(msg)
	logWriter.Flush()
	return execs
}

func ThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, prio int) map[string]int {
	execs := MakeExecs()
	f3 := 0

	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, prio)
		prio = 0
		kFallback := hypergraph.F3TargetLowDegree(g, c)
		execs["kFallback"] += kFallback
		f3++
		//prio = nextPrio
	}
	return execs
}

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

	switch prio {
	case 2:
		exec := hypergraph.SmallTriangleRule(g, c)
		execs["kTri"] += exec
	}

	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kVertDom := hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom := hypergraph.EdgeDominationRule(g)
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule2(g, c)
		//kSmallEdgeDegTwo := hypergraph.SmallEdgeDegreeTwoRule(g, c)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kExtTri := hypergraph.ExtendedTriangleRule(g, c)
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		//execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom+kExtTri == 0 {
			break
		}
	}

	fmt.Println(execs)

	return execs
}

func ThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()
	ApplyRules(g, c, execs, 0)
	expDepth := 2

	if len(g.Edges) == 0 {
		return execs
	}
	fmt.Println(execs)

	e := hypergraph.F3TargetLowDegreeDetect(g)
	if e != -1 {
		execs["kFallback"] += 1
		for v := range g.Edges[e].V {
			c[v] = true
		}
	}

	gf := hypergraph.GetFrontierGraph(g, expDepth, e)
	fmt.Println(len(gf.Edges))

	for len(gf.IncMap) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesFrontier(gf, g, c, execs, expand)
		if len(expand) > 0 {
			oldSizeE := len(gf.Edges)
			hypergraph.ExpandFrontier(gf, g, expDepth, expand)
			fmt.Printf("Expand added %d new edges\n", len(gf.Edges)-oldSizeE)
			continue
		}

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			fmt.Println("Could not find size 3 edge")
			continue
		}

		expand = make(map[int32]bool)

		for v := range g.Edges[e].V {
			if gf.VertexFrontier[v] {
				expand[v] = true
			}
			c[v] = true
		}

		if _, ex := gf.Edges[e]; !ex {
			hypergraph.F3_ExpandFrontier(gf, g, e, expDepth)
		} else if len(expand) > 0 {
			hypergraph.ExpandFrontier(gf, g, expDepth, expand)
			for v := range g.Edges[e].V {
				for e := range g.IncMap[v] {
					gf.F_RemoveEdge(e, g)
				}
			}
		} else {
			for v := range g.Edges[e].V {
				for e := range g.IncMap[v] {
					gf.F_RemoveEdge(e, g)
				}
			}
		}
		execs["kFallback"] += 1

		fmt.Println(len(gf.Edges), len(gf.IncMap), len(g.Edges), execs["kFallback"])
	}
	fmt.Println(len(gf.Edges), len(gf.IncMap), len(g.Edges))

	return execs
}

func ApplyRulesFrontier(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	for {
		kVertDom := hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kTiny := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule(gf, g, c, expand)
		//kSmallEdgeDegTwo, l9 := hypergraph.S_SmallEdgeDegreeTwoRule(gf,g, c)
		kTri := hypergraph.S_SmallTriangleRule(gf, g, c, expand)
		kExtTri := hypergraph.S_ExtendedTriangleRule(gf, g, c, expand)
		kSmall := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		//execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		if kTiny+kTri+kSmall+kEdgeDom+kVertDom+kExtTri+kApVertDom+kApDoubleVertDom == 0 {
			break
		}
	}
}

func ApplyRulesRand(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

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

		kApVertDom := 0
		kApDoubleVertDom := 0
		kSmallEdgeDegTwo := 0
		kTri := 0
		kExtTri := 0
		kSmall := 0

		perm := make([]int, 6)
		for i := range perm {
			perm[i] = i
		}

		Shuffle[int](perm)

		for i := 0; i < 6; i++ {
			switch perm[i] {
			case 0:
				kApVertDom = hypergraph.ApproxVertexDominationRule(g, c)
			case 1:
				kApDoubleVertDom = hypergraph.ApproxDoubleVertexDominationRule2(g, c)
			case 2:
				kSmallEdgeDegTwo = hypergraph.SmallEdgeDegreeTwoRule(g, c)
			case 3:
				kTri = hypergraph.SmallTriangleRule(g, c)
			case 4:
				kExtTri = hypergraph.ExtendedTriangleRule(g, c)
			case 5:
				kSmall = hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
			}
		}

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kEdgeDom+kVertDom+kTri+kSmall+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kExtTri == 0 {
			break
		}
	}

	fmt.Println(execs)

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
