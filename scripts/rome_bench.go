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
	flag.Parse()

	dir, _ := os.ReadDir("./graphs/rome")
	masterfile, _ := os.Create(*o)
	masterfile.WriteString("File;Ratio;")
	masterfile.WriteString(strings.Join(alg.Labels, ";"))
	masterfile.WriteString(";OVertices;OEdges;Vertices;Edges;HittingSet;Opt\n")

	bufWriter := bufio.NewWriterSize(masterfile, 65536)

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

			execs := alg.ThreeHS_F3ApprPolyFrontierSingle(h, c)

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
