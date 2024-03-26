package hypergraph

import (
	"bufio"
	_ "embed"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
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

	if len(strings.Trim(filename, " \n")) != 0 {
		outFilename = fmt.Sprintf("./graphs/%s.graphml", filename)
	} else {
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

func WriteToFileSimple(g *HyperGraph, filepath string) bool {
	f, err := os.Create(filepath)

	if err != nil {
		log.Default().Println(err)
		log.Default().Printf("Could not create file %s\n", filepath)
		return false
	}

	bf := bufio.NewWriterSize(f, 8192)

	defer f.Close()

	for _, e := range g.Edges {
		line := ""
		i := 0
		for v := range e.V {
			if i == len(e.V)-1 {
				line += fmt.Sprintf("%d\n", v)
				break
			}
			line += fmt.Sprintf("%d ", v)
			i++
		}
		bf.Write([]byte(line))
	}
	bf.Flush()

	return true
}

func ReadFromFileSimpleCallback(filename string, callback func(line string)) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Could not open file '%s'", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		callback(scanner.Text())
	}
}

func ReadFromFileSimple(filename string) *HyperGraph {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Could not open file '%s'", filename)
	}
	defer file.Close()

	g := NewHyperGraph()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line[0] == '#' {
			continue
		}
		lines = append(lines, line)
	}

	for _, line := range lines {
		spl := strings.Fields(line)
		splInt32 := make([]int32, len(spl))
		for i, v := range spl {
			id, _ := strconv.Atoi(v)
			splInt32[i] = int32(id)
			g.AddVertex(int32(id), 0)
		}
		g.AddEdge(splInt32...)
	}

	return g
}

func ReadFromFile(filename string) *HyperGraph {
	extSpl := strings.Split(filename, ".")
	ext := extSpl[len(extSpl)-1]
	if ext == "txt" {
		return ReadFromFileSimple(filename)
	}

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

	for _, e := range graph.Graph.HyperEdges {
		edges := make([]int32, len(e.Endpoints))
		for i, ep := range e.Endpoints {
			edges[i] = ep.Node
		}
		g.AddEdge(edges...)
	}

	return g
}

func ReadFromFileRome(filename string) *HyperGraph {
	extSpl := strings.Split(filename, ".")
	ext := extSpl[len(extSpl)-1]
	if ext == "txt" {
		return ReadFromFileSimple(filename)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read from file '%s'", filename)
	}

	var graph StdGraphMl

	err = xml.Unmarshal(file, &graph)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Could not unmarshal graph from file '%s'", filename)
	}

	g := NewHyperGraph()

	for _, v := range graph.Graph.Nodes {
		vInt, _ := strconv.ParseInt(v.Id[1:], 10, 32)
		g.AddVertex(int32(vInt), 0)
	}

	for _, e := range graph.Graph.Edges {
		source, _ := strconv.ParseInt(e.Source[1:], 10, 32)
		target, _ := strconv.ParseInt(e.Target[1:], 10, 32)
		g.AddEdge(int32(source), int32(target))
	}

	return g
}

func WriteToFileSetCover(g *HyperGraph, filename string) {
	out, _ := os.Create(filename)
	bufWriter := bufio.NewWriterSize(out, 8192)
	hashes := make(map[string]bool)

	for v, inc := range g.IncMap{
		l := fmt.Sprintf("%d: ", v)

		ep := []int32{}

		for e := range inc{
			ep = append(ep, e)
			l += strconv.Itoa(int(e)) + " "
		}

		hash := GetHash(ep...)
		
		if hashes[hash]{
			continue
		}else{
			hashes[hash] = true
		}

		l = l[0:len(l)-1] + "\n"
		bufWriter.WriteString(l)
	}
	bufWriter.Flush()
}
