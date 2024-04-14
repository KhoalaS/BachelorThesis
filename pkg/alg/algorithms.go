package alg

import (
	"fmt"
	//"log"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

var Labels = []string{"kTiny", "kVertDom", "kEdgeDom", "kSmall", "kTri", "kExtTri", "kApVertDom", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kSmallEdgeDegTwo2", "kFallback"}

var Ratios = map[string]IntTuple{
	"kTiny":             {A: 1, B: 1},
	"kSmall":            {A: 1, B: 2},
	"kTri":              {A: 2, B: 3},
	"kExtTri":           {A: 2, B: 4},
	"kApVertDom":        {A: 1, B: 2},
	"kApDoubleVertDom":  {A: 1, B: 2},
	"kSmallEdgeDegTwo":  {A: 2, B: 4},
	"kSmallEdgeDegTwo2": {A: 2, B: 3},
	"kFallback":         {A: 1, B: 3},
}

func ThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()

	for len(g.Edges) > 0 {
		ApplyRules(g, c, execs)
		kFallback := hypergraph.F3TargetLowDegree(g, c)
		execs["kFallback"] += kFallback
	}
	return execs
}

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int) {
	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kVertDom := hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom := hypergraph.EdgeDominationRule(g)
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.SmallEdgeDegreeTwoRule(g, c)
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
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom+kExtTri+kSmallEdgeDegTwo+kSmallEdgeDegTwo2 == 0 {
			break
		}
	}
}

func ThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()
	ApplyRules(g, c, execs)
	expDepth := 2

	if len(g.Edges) == 0 {
		return execs
	}

	e := hypergraph.F3TargetLowDegreeDetect(g)
	if e != -1 {
		execs["kFallback"] += 1
		for v := range g.Edges[e].V {
			c[v] = true
		}
	}

	gf := hypergraph.F3_ExpandFrontier(g, e, expDepth)

	for len(g.Edges) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesFrontier(gf, g, c, execs, expand)
		if len(expand) > 0 {
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			continue
		}

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			continue
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		execs["kFallback"] += 1
	}
	return execs
}

func ApplyRulesFrontier(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	for {
		kTiny := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kVertDom := hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
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
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kTri+kSmall+kEdgeDom+kVertDom+kExtTri+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kSmallEdgeDegTwo2 == 0 {
			break
		}
	}
}

func GreedyHighDeg(g *hypergraph.HyperGraph, c map[int32]bool) {

	for len(g.Edges) > 0 {
		max := 0
		var remVertex int32 = -1

		for v := range g.Vertices {
			d := g.Deg(v)
			if d > max {
				max = d
				remVertex = v
			}
		}

		for e := range g.IncMap[remVertex] {
			g.RemoveEdge(e)
		}
		c[remVertex] = true
	}
}

func ApplyRulesSingle(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool, exact bool) {
	if exact {
		kTiny := 0
		kVertDom := 0
		kEdgeDom := 0

		for {
			old := kTiny + kVertDom + kEdgeDom

			kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
			kVertDom += hypergraph.S_VertexDominationRule(gf, g, c, expand)
			kEdgeDom += hypergraph.S_EdgeDominationRule(gf, g, expand)

			new := kTiny + kVertDom + kEdgeDom
			if old == new {
				break
			}
		}

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
	}

	kApVertDom := hypergraph.FS_ApproxVertexDominationRule(gf, g, c, expand)
	if kApVertDom > 0 {
		execs["kApVertDom"] += kApVertDom
		return
	}
	kApDoubleVertDom := hypergraph.FS_ApproxDoubleVertexDominationRule(gf, g, c, expand)
	if kApDoubleVertDom > 0 {
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		return
	}
	kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.FS_SmallEdgeDegreeTwoRule(gf, g, c, expand)
	if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
		return
	}
	kTri := hypergraph.FS_SmallTriangleRule(gf, g, c, expand)
	if kTri > 0 {
		execs["kTri"] += kTri
		return
	}
	kExtTri := hypergraph.FS_ExtendedTriangleRule(gf, g, c, expand)
	if kExtTri > 0 {
		execs["kExtTri"] += kExtTri
		return
	}
	kSmall := hypergraph.FS_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
	if kSmall > 0 {
		execs["kSmall"] += kSmall
		return
	}
}

func PreProcessOnly(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	kTiny := 0
	kVertDom := 0
	kEdgeDom := 0

	for {
		old := kTiny + kVertDom + kEdgeDom

		kTiny += hypergraph.FS_TinyEdgeRule(g, c, expand)
		kVertDom += hypergraph.FS_VertexDominationRule(g, expand)
		kEdgeDom += hypergraph.FS_EdgeDominationRule(g, expand)

		new := kTiny + kVertDom + kEdgeDom
		if old == new {
			break
		}
	}
	execs["kTiny"] += kTiny
	execs["kVertDom"] += kVertDom
	execs["kEdgeDom"] += kEdgeDom
}

func ThreeHS_F3ApprPolyFrontierSingle(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()
	expand := make(map[int32]bool)

	PreProcessOnly(g, c, execs, expand)

	if len(g.Edges) == 0 {
		return execs
	}

	expDepth := 1

	gf := hypergraph.ExpandFrontier(g, expDepth, expand)

	for len(g.Edges) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesSingle(gf, g, c, execs, expand, true)

		if len(expand) > 0 {
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			continue
		} else {
			ApplyRulesSingle(g, g, c, execs, expand, false)
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			if len(expand) > 0 {
				continue
			}
		}

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			fmt.Println("No size 3 edge")
			continue
		}

		//log.Default().Println("F3")

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		execs["kFallback"] += 1
	}
	return execs
}
