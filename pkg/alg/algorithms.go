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

func ThreeHS_2ApprBranchOnly(g *hypergraph.HyperGraph, c map[int32]bool, k int) bool {
	//_, k := ApplyRules(g, c, K)

	//c := make(map[int32]bool)

	if k < 0 {
		return false
	}



	return true
}

func ThreeHS_2ApprGeneral(g *hypergraph.HyperGraph, c map[int32]bool, K int) bool {
	_, k := ApplyRules(g, c, K)
	if k < 0 {
		return false
	}

	return true
}



func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, K int) (map[string]int, int) {

	execs := make(map[string]int)

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

		log.Default().Println("#Edges: ", g.Edges)

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

	//m, err := os.Create("mem_main.prof")
	//if err != nil {
	//	log.Fatal("could not create memory profile: ", err)
	//}
	//defer m.Close()
	//if err := pprof.WriteHeapProfile(m); err != nil {
	//	log.Fatal("could not write memory profile: ", err)
	//}

	return execs, K
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
		log.Fatal(err)
	}
	defer os.Remove(script.Name())

	_, err = script.WriteString(scriptCode)
	if err != nil {
		log.Fatal(err)
	}

	pyCmd := exec.Command("python3", script.Name(), f.Name())

	out, err := pyCmd.CombinedOutput()
	if err != nil {
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