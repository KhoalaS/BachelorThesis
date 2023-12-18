package alg

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/KhoalaS/BachelorThesis/pkg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

//go:embed minedgecover.py
var scriptCode string

var Ratios = map[string]pkg.IntTuple{
	"kTiny":            {A: 1, B: 1},
	"kSmall":           {A: 2, B: 1},
	"kTri":             {A: 3, B: 2},
	"kApVertDom":       {A: 2, B: 1},
	"kApDoubleVertDom": {A: 2, B: 1},
}

func ThreeHS_2ApprBranchOnly(g *hypergraph.HyperGraph, c map[int32]bool, K int) bool {
	//_, k := ApplyRules(g, c, K)

	if K < 0 {
		return false
	}

	

	return true
}

func ThreeHS_2ApprGeneral(g *hypergraph.HyperGraph, c map[int32]bool, K int) (bool, map[int32]bool) {
	execs, k := ApplyRules(g, c, K)
	fmt.Println(execs)

	if k < 0 {
		return false, make(map[int32]bool)
	}

	if len(g.Edges) > 0 {
		g_n := g.Copy()
		c_n := make(map[int32]bool)
		for key, val := range c {
			c_n[key] = val
		}

		v, ex := PotentialTriangle(g)
		if ex {
			delete(g_n.Vertices, v)
			for _, e := range g_n.Edges {
				if e.V[v] {
					delete(e.V, v)
				}
			}			
			ThreeHS_2ApprGeneral(g_n, c_n, k)
		} else if g.IsSimple() {
			cover := MinEdgeCover(g)
			if k - len(cover) > 0 {
				for _, w := range cover {
					c[w] = true					
				}
				return true, c
			}else{
				return false, make(map[int32]bool)
			}
		}
	}

	return true, c
}



func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, K int) (map[string]int, int) {

	execs := make(map[string]int)

	k := K

	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom := hypergraph.EdgeDominationRule(g, c)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kVertDom := hypergraph.VertexDominationRule(g,c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kApVertDom := hypergraph.ApproxVertexDominationRule3(g, c, false)
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		//kApDoubleVertDom := 0

		//log.Default().Println("#Edges: ", g.Edges)

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom == 0 {
			break
		}
	}

	k -= execs["kTiny"] * Ratios["kTiny"].B
	k -= execs["kTri"] * Ratios["kTri"].B
	k -= execs["kApVertDom"] * Ratios["kApVertDom"].B
	k -= execs["kSmall"] * Ratios["kSmall"].B
	k -= execs["kApDoubleVertDom"] * Ratios["kApDoubleVertDom"].B

	//m, err := os.Create("mem_main.prof")
	//if err != nil {
	//	log.Fatal("could not create memory profile: ", err)
	//}
	//defer m.Close()
	//if err := pprof.WriteHeapProfile(m); err != nil {
	//	log.Fatal("could not write memory profile: ", err)
	//}

	return execs, k
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
						return v, true
					}
				}
			}else{
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
			for eId := range incList[v]{
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
	output = output[1:len(output)-1]
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