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
	p := flag.Float64("p", 0.05, "probability of adding an edge")
	evr := flag.Float64("evr", 0.0, "targetted edge/vertex ratio, takes priority over p")
	logging := flag.Int("log", 1, "log the number of rule executions, do log many runs")
	outdir := flag.String("d", "./data", "output directory")
	profile := flag.Bool("prof", false, "make pprof profile")
	frontier := flag.Bool("fr", false, "use frontier algorithm")
	cvd := flag.Bool("cvd", false, "reduce to cluster vertex deletion")
	debug := flag.Bool("dbg", false, "enable debug printing for rules")
	tvdi := flag.String("tvdi", "", "use graph file that is already a TVD instance")
	gr := flag.Bool("gr", false, "")
	gml := flag.Bool("gml", false, "")

	flag.Parse()

	hypergraph.Logging = *debug

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

		if len(*in) > 0 && !*gml {
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
		} else if flagPassed("tvdi") {
			// pass
		} else if *gml {
			// pass
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

		if *cvd {
			var j *hypergraph.HyperGraph
			if *gml {
				j = hypergraph.ReadFromFileRome(*in)
			} else {
				j = hypergraph.NewHyperGraph()
				for v, adj := range adjList {
					j.AddVertex(v, 0)
					for w := range adj {
						j.AddEdge(v, w)
					}
				}
			}

			fmt.Println("Start P3 detection and problem reduction...")
			g = hypergraph.P3Detection(j)
			fmt.Printf("Graph had %d many P3's\n", len(g.Edges))
		} else if flagPassed("tvdi") {
			fmt.Println("Read from TVD instance file...")
			g = hypergraph.ReadFromFileSimple(*tvdi)
			fmt.Printf("Graph had %d many triangles\n", len(g.Edges))
		} else {
			fmt.Println("Start Triangle detection and problem reduction...")
			g = hypergraph.TriangleDetection(adjList)
			fmt.Printf("Graph had %d many triangles\n", len(g.Edges))
		}

		c := make(map[int32]bool)

		var execs map[string]int

		if *profile {
			fmt.Println("Start CPU profile...")
			f, err := os.Create("./default.pgo")
			if err != nil {
				return
			}
			pprof.StartCPUProfile(f)
		}

		fmt.Printf("Problem instance has %d vertices and %d edges\n", len(g.Vertices), len(g.Edges))
		fmt.Println("Start 3-HS algorithm...")
		defer hypergraph.LogTime(time.Now(), "Main Algorithm")
		if *frontier {
			if flagPassed("log") {
				execs = alg.LoggingThreeHS_F3ApprPolyFrontier(g, c, graphtype, masterfilename, i, *outdir)
			} else {
				execs = alg.ThreeHS_F3ApprPolyFrontier(g, c)
			}
		} else {
			if flagPassed("log") {
				execs = alg.LoggingThreeHS_F3ApprPoly(g, c, graphtype, masterfilename, i, *outdir)
			} else if *gr {
				alg.GreedyHighDeg(g, c)
			} else {
				execs = alg.ThreeHS_F3ApprPoly(g, c, 0)
			}
		}

		pprof.StopCPUProfile()
		fmt.Println(execs)
		fmt.Printf("Found a hitting-set with size %d\n", len(c))
		fmt.Println("Est. Approximation Factor:", alg.GetRatio(execs))
	}
}
