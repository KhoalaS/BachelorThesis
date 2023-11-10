package main

import (
	"com/khoa/thesis/pkg/hypergraph"
	"fmt"
)

func main(){

	vertices := []hypergraph.Vertex{}
	edges := []hypergraph.Edge{}
	
	for i := 0; i<5; i++ {
		vertices = append(vertices, hypergraph.NewVertex(i,0))
	}

	edges = append(edges, hypergraph.NewEdge(0,1,2))
	edges = append(edges, hypergraph.NewEdge(0,1))
	edges = append(edges, hypergraph.NewEdge(2,3))
	edges = append(edges, hypergraph.NewEdge(4))

	g := hypergraph.NewHyperGraph(vertices, edges)
	c := make(map[int]bool)
	g_1, _ := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
	g_2 := hypergraph.EdgeDominationRule(g_1)
	
	fmt.Printf("Graph g is simple: %v\n", g.IsSimple())
	g.Print()

	fmt.Println("|After Tiny Edge Rule|")
	g_1.Print()

	fmt.Println("|After Edge Domination Rule|")
	g_2.Print()
}