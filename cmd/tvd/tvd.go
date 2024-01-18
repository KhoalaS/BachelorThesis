package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func ReduceToHS(t []map[int32]bool) *hypergraph.HyperGraph{
	g := hypergraph.NewHyperGraph()

	hashes := make(map[string]bool)
	arr := make([]int32, 3)

	for _, triangle := range t {
		i := 0
		for v := range triangle {
			arr[i] = v
			i++
		}
		hash := hypergraph.GetHash(arr)
		if hashes[hash] {
			continue
		}else{
			for v := range triangle {
				g.AddVertex(v, 0)
			}
			g.AddEdge(arr...)
			hashes[hash] = true
		}
	}
	return g
}

func main(){
	in := flag.String("i", "", "path to input graph file")
	n := flag.Int("n", 1000, "number of vertices")
	p := flag.Float64("p", 0.5, "probability of adding an edge")
	evr := flag.Float64("evr", 0.0, "targetted edge/vertex ratio, takes priority over p")
	logging := flag.Int("log", 0, "log the number of rule executions, do log many runs")

	flag.Parse()

	var g *hypergraph.HyperGraph

	graphtype := "CUSTOM"
	if len(*in) > 0 {
		g = hypergraph.ReadFromFileSimple(*in)
	}else{
		g = hypergraph.UniformERGraph(*n, *p, *evr, 2)
		graphtype = "ER"
	}

	t := hypergraph.TriangleDetection(g)

	g = ReduceToHS(t)

	fmt.Printf("Found %d triangles\n",len(g.Edges))
	c := make(map[int32]bool)

	if *logging > 0 {
		l_evr := float64(*evr)
		if *evr == 0 {
			l_evr = float64(len(g.Edges)) / float64(len(g.Vertices))
		}
		timestamp := time.Now().Unix()
		masterfilename := fmt.Sprintf("master_%s_%.2f_%d.csv", graphtype, l_evr, timestamp)
		for i := 0; i < *logging; i++ {
			alg.LoggingThreeHS_F3ApprPoly(g, c, graphtype, masterfilename, i)
			if i == *logging - 1 {
				break
			}

			if len(*in) > 0 {
				g = hypergraph.ReadFromFileSimple(*in)
			}else{
				g = hypergraph.UniformERGraph(*n, *p, *evr, 2)
			}
			
			t = hypergraph.TriangleDetection(g)
			g = ReduceToHS(t)
			c = make(map[int32]bool)
		}
	}else{
		execs := alg.ThreeHS_F3ApprPoly(g, c, 0)
		fmt.Println(execs)
		fmt.Println("Est. Approximation Factor:", alg.GetRatio(execs))
	}

}

