package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func flagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	in := flag.String("i", "", "path to input graph file")
	n := flag.Int("n", 1000, "number of vertices")
	p := flag.Float64("p", 0.5, "probability of adding an edge")
	evr := flag.Float64("evr", 0.0, "targetted edge/vertex ratio, takes priority over p")
	logging := flag.Int("log", 1, "log the number of rule executions, do log many runs")
	outdir := flag.String("d", "./data", "output directory")

	flag.Parse()

	adjList := make(map[int32]map[int32]bool)

	graphtype := "CUSTOM"
	if !flagPassed("i") {
		graphtype = "ER"
	}
	var g *hypergraph.HyperGraph

	timestamp := time.Now().Unix()
	masterfilename := fmt.Sprintf("master_%s_%d.csv", graphtype, timestamp)

	for i := 0; i < *logging; i++ {
		adjList = make(map[int32]map[int32]bool)
		var a int
		var b int

		if len(*in) > 0 {
			hypergraph.ReadFromFileSimpleCallback(*in, func(line string) {
				if line[0] == '#' {
					return
				}
				spl := strings.Fields(line)
				a, _ = strconv.Atoi(spl[0])
				b, _ = strconv.Atoi(spl[1])
				if _, ex := adjList[int32(a)]; !ex {
					adjList[int32(a)] = make(map[int32]bool)
				}
				if _, ex := adjList[int32(b)]; !ex {
					adjList[int32(b)] = make(map[int32]bool)
				}
				adjList[int32(a)][int32(b)] = true
				adjList[int32(b)][int32(a)] = true
			})
		} else {
			hypergraph.UniformERGraphCallback(*n, *p, *evr, 2, func(edge []int32) {
				if _, ex := adjList[edge[0]]; !ex {
					adjList[edge[0]] = make(map[int32]bool)
				}
				if _, ex := adjList[edge[1]]; !ex {
					adjList[edge[1]] = make(map[int32]bool)
				}
				adjList[edge[0]][edge[1]] = true
				adjList[edge[1]][edge[0]] = true

			})
			graphtype = "ER"
		}

		fmt.Println("Start Triangle detection and problem reduction")
		g = hypergraph.TriangleDetection(adjList)

		fmt.Println(len(g.Vertices))

		c := make(map[int32]bool)

		var execs map[string]int

		fmt.Println("Start 3-HS algorithm")
		if flagPassed("log") {
			execs = alg.LoggingThreeHS_F3ApprPoly(g, c, graphtype, masterfilename, i, *outdir)
		} else {
			execs = alg.ThreeHS_F3ApprPoly(g, c, 0)
		}
		fmt.Println(execs)
		fmt.Println("Est. Approximation Factor:", alg.GetRatio(execs))
	}
}
