package hypergraph

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

//go:embed samplegraph.html
var graphTemplate string

func WriteToFile(g *HyperGraph, filename string) bool {
	tmpl := template.Must(template.New("graph").Parse(graphTemplate))

	err := os.Mkdir("./graphs", 0777)
	if err != nil {
		if err == os.ErrExist {
			log.Default().Println(err)
			log.Default().Println("Could not create folder './graphs'")
			return false
		}
	}

	var outFilename string
	
	if len(strings.Trim(filename, " \n")) != 0{
		outFilename = fmt.Sprintf("./graphs/%s.graphml", filename)
	}else{
		outFilename = fmt.Sprintf("./graphs/g_%d.graphml", time.Now().Unix())
	}
	
	f, err := os.Create(outFilename)
	
	if err != nil {
		log.Default().Println(err)
		log.Default().Printf("Could not create file %s\n", outFilename)
		return false
	}

	defer f.Close()


	err = tmpl.Execute(f, g)
	if err != nil {
		log.Default().Println(err)
		log.Default().Printf("Could not write to file %s\n", outFilename)
		return false
	}
	
	return true
}

func ReadFromFile(filename string) *HyperGraph {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read from file '%s'", filename)
	}

	var graph GraphMl 

	err = xml.Unmarshal(file, &graph)
	if err != nil {
		log.Fatalf("Could not unmarshal graph from file '%s'", filename)
	}

	g := NewHyperGraph()

	for _, v := range graph.Graph.Nodes {
		g.AddVertex(v.Id, v.Data.Value)
	}
	
	for _, e := range graph.Graph.Edges {
		edges := make([]int32, len(e.Endpoints))
		for i, ep := range e.Endpoints {
			edges[i] = ep.Node			
		}
		g.AddEdgeArr(edges)
	}

	return g
}
