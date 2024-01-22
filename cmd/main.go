package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func makeHypergraph(input string, u int, f string, n int, m int, prefAttach float64, prefAttachMod bool, er bool, evr int) (*hypergraph.HyperGraph, string) {
	var g *hypergraph.HyperGraph
	graphtype := "STD"

	if len(strings.Trim(input, " ")) > 0 {
		g = hypergraph.ReadFromFile(strings.Trim(input, " "))
		graphtype = "CUSTOM"
	} else if u > 0 {
		g = hypergraph.UniformTestGraph(int32(n), int32(m), u)
		graphtype = "---"
	} else if len(f) > 0 {
		spl := strings.Split(f, ",")
		ratios := make([]int, len(spl))
		for i, val := range spl {
			valInt, _ := strconv.Atoi(val)
			ratios[i] = valInt
		}
		g = hypergraph.FixDistTestGraph(int32(n), int32(m), ratios)
		graphtype = "FIX"
	} else if prefAttach > 0 {
		g = hypergraph.PrefAttachmentGraph(int32(n), prefAttach, 3)
		graphtype = "PREF"
	} else if prefAttachMod {
		g = hypergraph.ModPrefAttachmentGraph(int(n), 5, 0.5, 0.21)
		graphtype = "---"
	} else if er {
		g = hypergraph.UniformERGraph(int(n), 0.0, float64(evr), 3)
		graphtype = "ERU3"
	} else {
		g = hypergraph.TestGraph(int32(n), int32(m), true)
	}
	return g, graphtype
}

func main() {
	input := flag.String("i", "", "Filepath to input file.")
	n := flag.Int("n", 10000, "Number of vertices if no graph file supplied.")
	m := flag.Int("m", 20000, "Number of edges if no graph file supplied.")
	u := flag.Int("u", 0, "Generate a u-uniform graph.")
	f := flag.String("f", "", "Generate a random hypergraph with fixed ratios for the edge sizes.")
	evr := flag.Int("evr", 0, "Maximum ratio |E|/|V| to compute for random graphs.")
	profile := flag.Bool("prof", false, "Make CPU profile")
	export := flag.String("o", "", "Export the generated graph with the given string as filename. The will create a 'graphs' folder where the file is located.")
	exportSimple := flag.String("os", "", "Export the generated graph to the given filepath.")
	prefAttach := flag.Float64("pa", 0.0, "Generate a random preferential attachment hypergraph with given float as probablity to add a new vertex.")
	prefAttachMod := flag.Bool("pamod", false, "")
	er := flag.Bool("er", false, "")
	logging := flag.Int("log", 0, "")
	outdir := flag.String("d", "./data", "")

	flag.Parse()

	g, graphtype := makeHypergraph(*input, *u, *f, *n, *m, *prefAttach, *prefAttachMod, *er, *evr)
	c := make(map[int32]bool)

	if len(*export) > 0 {
		hypergraph.WriteToFile(g, *export)
		return
	}

	if len(*exportSimple) > 0 {
		hypergraph.WriteToFileSimple(g, *exportSimple)
		return
	}

	fmt.Println("Start Algorithm")

	if *profile {
		fmt.Println("Start CPU profile...")
		f, err := os.Create("profiles/benchmark_main.prof")
		if err != nil {
			return
		}

		pprof.StartCPUProfile(f)
	}

	prio := 0
	var execs map[string]int

	if *logging > 0 {
		l_evr := float64(*evr)
		if *evr == 0 {
			l_evr = float64(len(g.Edges)) / float64(len(g.Vertices))
		}
		t := time.Now().Unix()
		masterfilename := fmt.Sprintf("master_%s_%.2f_%d.csv", graphtype, l_evr, t)
		for i := 0; i < *logging; i++ {
			alg.LoggingThreeHS_F3ApprPoly(g, c, graphtype, masterfilename, i, *outdir)
			if i == *logging - 1 {
				break
			}
			g, _ = makeHypergraph(*input, *u, *f, *n, *m, *prefAttach, *prefAttachMod, *er, *evr)
			c = make(map[int32]bool)
		}
	} else {
		execs = alg.ThreeHS_F3ApprPoly(g, c, prio)
		fmt.Printf("Found a 3-Hitting-Set of size %d\n", len(c))
		fmt.Printf("Estimated Approximation Factor: %.2f\n", alg.GetRatio(execs))
		fmt.Println(execs)
	}
	pprof.StopCPUProfile()
}
