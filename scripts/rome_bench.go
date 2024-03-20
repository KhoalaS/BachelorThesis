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
	flag.Parse()

	dir, _ := os.ReadDir("./graphs/rome")
	masterfile, _ := os.Create(*o)
	masterfile.WriteString("File;Ratio;")
	masterfile.WriteString(strings.Join(alg.Labels, ";"))
	masterfile.WriteString(";OVertices;OEdges;Vertices;Edges;HittingSet;Opt\n")

	bufWriter := bufio.NewWriterSize(masterfile, 65536)

	if len(*i) > 0 {
		g := hypergraph.ReadFromFileRome(*i)
		h := hypergraph.P3Detection(g)
		eOSize := len(g.Edges)
		vOSize := len(g.Vertices)
		eSize := len(h.Edges)
		vSize := len(h.Vertices)

		c := make(map[int32]bool)
		execs := alg.MakeExecs()

		for len(h.Edges) > 0 {
			execs["kTiny"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)
			execs["kVertDom"] += hypergraph.VertexDominationRule(h, c)
			execs["kTiny"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)
			hypergraph.WriteToFileSimple(h, "./out/rome_single_pre_edom.txt")
			execs["kEdgeDom"] += hypergraph.EdgeDominationRule(h)
			execs["kApVertDom"] += hypergraph.ApproxVertexDominationRule(h, c)
			execs["kApDoubleVertDom"] += hypergraph.ApproxDoubleVertexDominationRule(h, c)
			k0, k1 := hypergraph.SmallEdgeDegreeTwoRule(h, c)
			execs["kSmallEdgeDegTwo"] += k0
			execs["kSmallEdgeDegTwo2"] += k1
			execs["kTri"] += hypergraph.SmallTriangleRule(h, c)
			execs["kExtTri"] += hypergraph.ExtendedTriangleRule(h, c)
			execs["kSmall"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.SMALL)
		}

		rules := ""
		for _, label := range alg.Labels {
			rules += strconv.Itoa(execs[label])
			rules += ";"
		}
		bufWriter.WriteString((fmt.Sprintf("%s;%f;%s%d;%d;%d;%d;%d;%d\n", *i, alg.GetRatio(execs), rules, vOSize, eOSize, vSize, eSize, len(c), alg.GetEstOpt(execs))))
		bufWriter.Flush()
		return
	}

	nFiles := len(dir)

	for idx, file := range dir {
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
			bufWriter.WriteString((fmt.Sprintf("%s;%f;%s%d;%d;%d;%d;%d;%d\n", file.Name(), alg.GetRatio(execs), rules, vOSize, eOSize, vSize, eSize, len(c), alg.GetEstOpt(execs))))
		}
		fmt.Printf("(%d/%d) finished for file:%s\n", idx+1, nFiles, file.Name())
	}
	bufWriter.Flush()
}
