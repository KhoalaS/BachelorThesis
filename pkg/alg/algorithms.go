package alg

import (
	"log"
	"github.com/KhoalaS/BachelorThesis/pkg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

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
}