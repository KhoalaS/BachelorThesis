package main

import (
	"com/khoa/thesis/pkg/hypergraph"
	"fmt"
	"math/rand"
)

func main(){

	g := hypergraph.NewHyperGraph()
	
	var i int32 = 0
	var vSize int32 = 100000
	var eSize int32 = 20000

	for ; i < vSize; i++ {
		g.AddVertex(i, 0)
	}

	i = 0

	for ; i < eSize; i++ {
		d := 1
		r := rand.Float32()
		if r > 0.1 && r < 0.6 {
			d = 2
		} else if r >= 0.6 {
			d = 3
		}
		eps := make(map[int32]bool)
		for j := 0; j < d; j++ {
			val := rand.Int31n(vSize)
			for eps[val] {
				val = rand.Int31n(vSize)
			}
			eps[val] = true
		}
		g.AddEdgeMap(eps)
	}

	c := make(map[int32]bool)

	//g.Print()
	
	hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
	
	fmt.Println("|After Tiny Edge Rule|")
	//g.Print()
	
	hypergraph.EdgeDominationRule(g, c)
	
	fmt.Println("|After Edge Domination Rule|")
	//g.Print()
	
	
	fmt.Printf("Graph g is simple: %v\n", g.IsSimple())
	//g.Print()
	//fmt.Println(c)

}