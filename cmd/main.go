package main

import (
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
	"fmt"
)

func main(){

	g := hypergraph.GenerateTestGraph(1000000, 2000000)

	c := make(map[int32]bool)
	
	hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
	
	fmt.Println("|After Tiny Edge Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	/*
	hypergraph.EdgeDominationRule(g, c)
	
	fmt.Println("|After Edge Domination Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))

	*/

	hypergraph.ApproxVertexDominationRule2(g, c)
	fmt.Println("|After Approx Vertex Domination Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	
	fmt.Printf("Graph g is simple: %v\n", g.IsSimple())
}