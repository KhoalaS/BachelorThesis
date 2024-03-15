package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func main() {
	o := flag.String("o", "./data/rome_master.csv", "path to output csv masterfile")
	n := flag.Int("n", 1, "number of algorithm runs per graph")
	i := flag.String("i", "", "input graphfile")
	hs := flag.Bool("hs", false, "enable graph history")
	flag.Parse()

	if *hs {
		hypergraph.HistoryEnabled = true
	}

	dir, _ := os.ReadDir("./graphs/rome")
	masterfile, _ := os.Create(*o)
	masterfile.WriteString("File;Ratio;")
	masterfile.WriteString(strings.Join(alg.Labels, ";"))
	masterfile.WriteString(";OVertices;OEdges;Vertices;Edges;HittingSet\n")

	bufWriter := bufio.NewWriterSize(masterfile, 65536)

	if len(*i) > 0 {
		g := hypergraph.ReadFromFileRome(*i)
		h := hypergraph.P3Detection(g)
		eOSize := len(g.Edges)
		vOSize := len(g.Vertices)
		eSize := len(h.Edges)
		vSize := len(h.Vertices)

		c := make(map[int32]bool)

		g.CurrentRule = "Tiny"
		hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)
		
		g.CurrentRule = "VDom"
		hypergraph.VertexDominationRule(g, c)
		
		g.CurrentRule = "Tiny"
		hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)
		
		g.CurrentRule = "Edom"
		hypergraph.EdgeDominationRule(g)
		
		g.CurrentRule = "AVD"
		hypergraph.ApproxVertexDominationRule(g, c)
		
		g.CurrentRule = "ADVD"
		hypergraph.ApproxDoubleVertexDominationRule(g, c)
		
		g.CurrentRule = "SED2"
		hypergraph.SmallEdgeDegreeTwoRule(g, c)
		
		g.CurrentRule = "Tri"
		hypergraph.SmallTriangleRule(g, c)
		
		g.CurrentRule = "ETri"
		hypergraph.ExtendedTriangleRule(g, c)
		
		g.CurrentRule = "Small"
		hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)

		execs := alg.ThreeHS_F3ApprPoly(h, c, 0)

		rules := ""
		for _, label := range alg.Labels {
			rules += strconv.Itoa(execs[label])
			rules += ";"
		}
		bufWriter.WriteString((fmt.Sprintf("%s;%f;%s%d;%d;%d;%d;%d\n", *i, alg.GetRatio(execs), rules, vOSize, eOSize, vSize, eSize, len(c))))
		bufWriter.Flush()
		return
	}

	for _, file := range dir {
		if file.Name() == "Graph.log" {
			continue
		}

		for i := 0; i < *n; i++ {
			g := hypergraph.ReadFromFileRome("./graphs/rome/" + file.Name())
			h := hypergraph.P3Detection(g)
			eOSize := len(g.Edges)
			vOSize := len(g.Vertices)
			eSize := len(h.Edges)
			vSize := len(h.Vertices)

			c := make(map[int32]bool)

			execs := alg.ThreeHS_F3ApprPolyFrontier(h, c)

			rules := ""
			for _, label := range alg.Labels {
				rules += strconv.Itoa(execs[label])
				rules += ";"
			}
			bufWriter.WriteString((fmt.Sprintf("%s;%f;%s%d;%d;%d;%d;%d\n", file.Name(), alg.GetRatio(execs), rules, vOSize, eOSize, vSize, eSize, len(c))))
		}
	}
	bufWriter.Flush()
}
