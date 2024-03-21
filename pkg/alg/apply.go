package alg

import "github.com/KhoalaS/BachelorThesis/pkg/hypergraph"

var exactIter = 3

func ApplyRulesStr1(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int) map[string]int {
	for {
		kVertDom := 0
		kTiny := 0
		kEdgeDom := 0
		for i := 0; i < exactIter; i++ {
			kVertDom += hypergraph.VertexDominationRule(g, c)
			kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
			kEdgeDom += hypergraph.EdgeDominationRule(g)
		}
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		kVertDom += hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom += hypergraph.EdgeDominationRule(g)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		kVertDom += hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kEdgeDom += hypergraph.EdgeDominationRule(g)
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
	return execs
}

func ApplyRulesFrontierStr1(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	for {
		kVertDom := 0
		kTiny := 0
		kEdgeDom := 0
		for i := 0; i < exactIter; i++ {
			kVertDom += hypergraph.S_VertexDominationRule(gf, g, c, expand)
			kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
			kEdgeDom += hypergraph.S_EdgeDominationRule(gf, g, expand)
		}
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		kVertDom += hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kEdgeDom += hypergraph.S_EdgeDominationRule(gf, g, expand)
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
		kVertDom += hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kEdgeDom += hypergraph.S_EdgeDominationRule(gf, g, expand)
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

func ApplyRulesStr2(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int) map[string]int {
	for {
		kVertDom := hypergraph.VertexDominationRule(g, c)
		execs["kVertDom"] += kVertDom
		
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		if kTiny > 0 {
			execs["kTiny"] += kTiny
			continue
		}
		kEdgeDom := hypergraph.EdgeDominationRule(g)
		if kEdgeDom > 0 {
			execs["kEdgeDom"] += kEdgeDom
			continue
		}
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		if kApVertDom > 0 {
			execs["kApVertDom"] += kApVertDom
			continue
		}
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		if kApDoubleVertDom > 0 {
			execs["kApDoubleVertDom"] += kApDoubleVertDom
			continue
		}
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.SmallEdgeDegreeTwoRule(g, c)
		if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
			execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
			execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
			continue
		}
		kTri := hypergraph.SmallTriangleRule(g, c)
		if kTri > 0 {
			execs["kTri"] += kTri
			continue
		}
		kExtTri := hypergraph.ExtendedTriangleRule(g, c)
		if kExtTri > 0 {
			execs["kExtTri"] += kExtTri
			continue
		}
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		if kSmall > 0 {
			execs["kSmall"] += kSmall
			continue
		}

		break
	}
	return execs
}

func ApplyRulesFrontierStr2(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	for {
		kVertDom := hypergraph.S_VertexDominationRule(gf, g, c, expand)
		execs["kVertDom"] += kVertDom

		kTiny := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		if kTiny > 0 {
			execs["kTiny"] += kTiny
			continue
		}
		kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		if kEdgeDom > 0 {
			execs["kEdgeDom"] += kEdgeDom
			continue
		}
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		if kApVertDom > 0 {
			execs["kApVertDom"] += kApVertDom
			continue
		}
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
		if kApDoubleVertDom > 0 {
			execs["kApDoubleVertDom"] += kApDoubleVertDom
			continue
		}
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
		if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
			execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
			execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
			continue
		}
		kTri := hypergraph.S_SmallTriangleRule(gf, g, c, expand)
		if kTri > 0 {
			execs["kTri"] += kTri
			continue
		}
		kExtTri := hypergraph.S_ExtendedTriangleRule(gf, g, c, expand)
		if kExtTri > 0 {
			execs["kExtTri"] += kExtTri
			continue
		}
		kSmall := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
		if kSmall > 0 {
			execs["kSmall"] += kSmall
			continue
		}

		break
	}
}