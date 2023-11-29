package main

import (
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
	"fmt"
)

func main(){
	g := hypergraph.GenerateTestGraph(10000, 100000, false)
	g.RemoveDuplicate()
	fmt.Println("After remove:", len(g.Edges))
	c := make(map[int32]bool)
	
	kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
	fmt.Println("|After Tiny Edge Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	fmt.Println(kTiny)

	kTri := hypergraph.SmallTriangleRule(g, c)
	fmt.Println("|After Small Triangle Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	fmt.Println(kTri)

	hypergraph.EdgeDominationRule(g, c)	
	fmt.Println("|After Edge Domination Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))

	kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
	fmt.Println("|After Small Edge Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	fmt.Println(kSmall)

	kApVertDom :=  hypergraph.ApproxVertexDominationRule3(g, c)
	fmt.Println("|After Approx Vertex Domination Rule3|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	fmt.Println(kApVertDom)
	
	approxFactor :=  float64(kTiny*1 + kSmall*2 + kApVertDom*2 + kTri*3)/float64(kTiny*1 + kSmall*1 + kApVertDom*1 + kTri*2)
	fmt.Println(approxFactor)
	
	//fmt.Println(len(c))

	//for ex := hypergraph.ApproxVertexDominationRule2(g, c); ex; ex = hypergraph.ApproxVertexDominationRule2(g, c){}
	//fmt.Println("|After Approx Vertex Domination Rule2|")
	//fmt.Println(len(c))
	//fmt.Println(len(g.Edges))
	
	fmt.Printf("Graph g is simple: %v\n", g.IsSimple())
}