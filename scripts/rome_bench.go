package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

func main() {
	dir, _ := os.ReadDir("./graphs/rome")
	masterfile, _ := os.Create("./data/rome_master.csv")
	masterfile.WriteString("File;Ratio;")
	masterfile.WriteString(strings.Join(alg.Labels, ";"))
	masterfile.WriteString(";OVertices;OEdges;Vertices;Edges;HittingSet\n")

	bufWriter := bufio.NewWriterSize(masterfile, 8192)

	for _, file := range dir {
		if file.Name() == "Graph.log" {
			continue
		}
		for i := 0; i < 100; i++ {
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
