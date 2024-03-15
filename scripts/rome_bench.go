package main

import (
	"bufio"
	"encoding/json"
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
			h.CurrentRule = "Tiny"
			execs["kTiny"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)

			h.CurrentRule = "VDom"
			execs["kVertDom"] += hypergraph.VertexDominationRule(h, c)

			h.CurrentRule = "Tiny"
			execs["kTiny"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.TINY)

			h.CurrentRule = "EDom"
			execs["kEdgeDom"] += hypergraph.EdgeDominationRule(h)

			h.CurrentRule = "AVD"
			execs["kApVertDom"] += hypergraph.ApproxVertexDominationRule(h, c)

			h.CurrentRule = "ADVD"
			execs["kApDoubleVertDom"] += hypergraph.ApproxDoubleVertexDominationRule(h, c)

			h.CurrentRule = "SED2"
			k0, k1 := hypergraph.SmallEdgeDegreeTwoRule(h, c)
			execs["kSmallEdgeDegTwo"] += k0
			execs["kSmallEdgeDegTwo2"] += k1

			h.CurrentRule = "Tri"
			execs["kTri"] += hypergraph.SmallTriangleRule(h, c)

			h.CurrentRule = "ETri"
			execs["kExtTri"] += hypergraph.ExtendedTriangleRule(h, c)

			h.CurrentRule = "Small"
			execs["kSmall"] += hypergraph.RemoveEdgeRule(h, c, hypergraph.SMALL)
		}

		if *hs {
			hsBytes, err := json.Marshal(h.History)
			if err != nil {
				panic(err)
			}
			hsFile, err := os.Create("./out/rome_single.json")
			if err != nil {
				panic(err)
			}
			hsFile.Write(hsBytes)
			defer hsFile.Close()
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
			bufWriter.WriteString((fmt.Sprintf("%s;%f;%s%d;%d;%d;%d;%d;%d\n", file.Name(), alg.GetRatio(execs), rules, vOSize, eOSize, vSize, eSize, len(c), alg.GetEstOpt(execs))))
		}
	}
	bufWriter.Flush()
}
