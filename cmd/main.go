package main

import (
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
	"fmt"
)

func main(){

	g := hypergraph.GenerateTestGraph(10000, 4000000)

	c := make(map[int32]bool)
	
	hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
	
	fmt.Println("|After Tiny Edge Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))

	hypergraph.EdgeDominationRule(g, c)
	
	fmt.Println("|After Edge Domination Rule|")
	fmt.Println(len(c))
	fmt.Println(len(g.Edges))
	
	fmt.Printf("Graph g is simple: %v\n", g.IsSimple())
}