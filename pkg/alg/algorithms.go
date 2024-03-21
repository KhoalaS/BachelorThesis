package alg

import (
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

var Labels = []string{"kTiny", "kVertDom", "kEdgeDom", "kSmall", "kTri", "kExtTri", "kApVertDom", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kSmallEdgeDegTwo2", "kFallback"}

var Ratios = map[string]IntTuple{
	"kTiny":             {A: 1, B: 1},
	"kSmall":            {A: 2, B: 1},
	"kTri":              {A: 3, B: 2},
	"kExtTri":           {A: 4, B: 2},
	"kApVertDom":        {A: 2, B: 1},
	"kApDoubleVertDom":  {A: 2, B: 1},
	"kSmallEdgeDegTwo":  {A: 4, B: 2},
	"kSmallEdgeDegTwo2": {A: 3, B: 2},
	"kFallback":         {A: 3, B: 1},
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
		//kEdgeDom := hypergraph.EdgeDominationRule(g)
		kEdgeDom := 0
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.SmallEdgeDegreeTwoRule(g, c)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kExtTri := hypergraph.ExtendedTriangleRule(g, c)
		//kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		kSmall := 0

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
	return execs
}

func ThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()
	ApplyRules(g, c, execs, 0)
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
		//kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		kEdgeDom := 0
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
		kTri := hypergraph.S_SmallTriangleRule(gf, g, c, expand)
		kExtTri := hypergraph.S_ExtendedTriangleRule(gf, g, c, expand)
		//kSmall := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
		kSmall := 0

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
