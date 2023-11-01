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
	edges = append(edges, hypergraph.NewEdge(2,3))
	edges = append(edges, hypergraph.NewEdge(4))

	g := hypergraph.NewHyperGraph(vertices, edges)
	c := make(map[int]bool)
	g_1, c_1 := hypergraph.RemoveEdges(g, c, hypergraph.TINY)
	fmt.Println(g.IsSimple())
	fmt.Println(g_1)
	fmt.Println(g)
	fmt.Println(c_1)
	fmt.Println(c)
}