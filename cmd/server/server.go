package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

var g *hypergraph.HyperGraph
var gf *hypergraph.HyperGraph

var c map[int32]bool
var execs map[string]int
var currFilename string
var initFrontier bool

var logging bool
var graphDir string
const adjLog = "temp/adj.log"
const incLog = "temp/inc.log"

var short_names = map[string]string{
	"kTiny":             "Tiny",
	"kVertDom":          "VD",
	"kEdgeDom":          "ED",
	"kSmall":            "Small",
	"kTri":              "Tri",
	"kExtTri":           "ETri",
	"kApVertDom":        "AVD",
	"kApDoubleVertDom":  "ADVD",
	"kSmallEdgeDegTwo":  "SED2",
	"kSmallEdgeDegTwo2": "SED2*",
	"kFallback":         "F3",
}

func getGraphs(w http.ResponseWriter, r *http.Request) {
	dir, err := os.ReadDir(graphDir)
	if err != nil {
		log.Default().Println(err)
		w.WriteHeader(400)
		return
	}

	files := make([]string, len(dir))
	for idx, file := range dir {
		if file.IsDir() {
			continue
		}
		files[idx] = file.Name()
	}

	body, _ := json.Marshal(files)
	w.Header().Add("content-type", "application/json")
	w.Write(body)
}

func setGraph(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("g")
	if len(filename) == 0 {
		w.WriteHeader(404)
		return
	}

	currFilename = filename

	filepath := fmt.Sprintf("./%s/%s",graphDir, filename)
	_, err := os.Stat(filepath)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	resetGraph()
	g = hypergraph.ReadFromFile(filepath)

	payload := getGraphPres(false)
	body, _ := json.Marshal(payload)
	w.Write(body)
}

func nextRule(w http.ResponseWriter, r *http.Request) {
	frontierOnly := len(r.URL.Query().Get("fronly")) > 0

	if logging {
		adjFile, err := os.OpenFile(adjLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		for v, val := range g.AdjCount {
			adjFile.WriteString(fmt.Sprintf("%d:%v\n", v, val))
		}
		adjFile.WriteString("--------------\n")

		incFile, err := os.OpenFile(incLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		for v, val := range g.IncMap {
			incFile.WriteString(fmt.Sprintf("%d:%v\n", v, val))
		}
		incFile.WriteString("--------------\n")
	}
	expDepth := 1
	expand := make(map[int32]bool)

	if len(g.Edges) == 0 {
		w.WriteHeader(204)
		return
	}

	if !initFrontier {
		initFrontier = true
		alg.PreProcessOnly(g, c, execs, expand)
		gf = hypergraph.ExpandFrontier(g, 1, expand)
		log.Default().Println("preprocess")
		payload := getGraphPres(frontierOnly)
		body, _ := json.Marshal(payload)
		w.Write(body)
		return
	}

	ApplyRulesSingle(gf, g, c, execs, expand, true)

	if len(expand) > 0 {
		gf = hypergraph.ExpandFrontier(g, expDepth, expand)
		payload := getGraphPres(frontierOnly)
		body, _ := json.Marshal(payload)
		w.Write(body)
		return
	} else {
		ApplyRulesSingle(g, g, c, execs, expand, false)
		gf = hypergraph.ExpandFrontier(g, expDepth, expand)
		if len(expand) > 0 {
			payload := getGraphPres(frontierOnly)
			body, _ := json.Marshal(payload)
			w.Write(body)
			return
		}
	}

	e := hypergraph.F3TargetLowDegreeDetect(g)
	if e == -1 {
		fmt.Println("No size 3 edge")
		return
	}

	//log.Default().Println("F3")

	for v := range g.Edges[e].V {
		c[v] = true
	}

	gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
	execs["kFallback"] += 1

	payload := getGraphPres(frontierOnly)
	body, _ := json.Marshal(payload)
	w.Write(body)
}

func getRandom(w http.ResponseWriter, r *http.Request) {
	resetGraph()
	if logging {
		os.Create(adjLog)
		os.Create(incLog)
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	var model ModelPayload
	json.Unmarshal(bodyBytes, &model)

	if model.Model == "PA" {
		g = hypergraph.PrefAttachmentGraph(int32(model.N), model.P, 3)
	} else if model.Model == "ER" {
		g = hypergraph.UniformERGraph(model.N, model.P, model.Evr, 3)
	}

	hypergraph.WriteToFileSimple(g, "./temp/rand.txt")

	payload := getGraphPres(false)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Default().Println(err)
	}
	w.Write(payloadBytes)
}

func resetCurrent(w http.ResponseWriter, r *http.Request) {
	g = hypergraph.ReadFromFile(currFilename)
	gf = hypergraph.NewHyperGraph()
	initFrontier = false

	for k := range execs {
		execs[k] = 0
	}
	for k := range c {
		delete(c, k)
	}

	payload := getGraphPres(false)
	body, _ := json.Marshal(payload)
	w.Write(body)
}

func resetGraph() {
	g = hypergraph.NewHyperGraph()
	gf = hypergraph.NewHyperGraph()

	initFrontier = false
	currFilename = ""

	for k := range c {
		delete(c, k)
	}

	for k := range execs {
		execs[k] = 0
	}
}

func main() {
	logging = *flag.Bool("log", false, "log adjaceny map and incidence map")
	graphDir = *flag.String("g", "./graphs", "path to graph files directory")	

	if logging {
		os.Create(adjLog)
		os.Create(incLog)
	}

	c = make(map[int32]bool)
	execs = alg.MakeExecs()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./cmd/server")))
	mux.HandleFunc("/graphs", getGraphs)
	mux.HandleFunc("/setgraph", setGraph)
	mux.HandleFunc("/nextrule", nextRule)
	mux.HandleFunc("/resetcurrent", resetCurrent)
	mux.HandleFunc("/random", getRandom)

	mainHandler := CorsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8097", mainHandler))
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

type GraphPayload struct {
	Graph     map[int32][]int32 `json:"graph"`
	C         []int32           `json:"c"`
	ExecsData []*ExecsTableRow  `json:"execs_data"`
	Vertices  []int32           `json:"vertices"`
	Ratio     float64           `json:"ratio"`
	Opt       int               `json:"opt"`
	Frontier  []int32           `json:"frontier"`
}

type ModelPayload struct {
	Model string  `json:"model"`
	P     float64 `json:"p"`
	Evr   float64 `json:"evr"`
	N     int     `json:"n"`
}

type ExecsTableRow struct {
	Name string `json:"name"`
	K    int    `json:"k"`
}

func ApplyRulesSingle(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool, exact bool) {

	if exact {
		kTiny := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		if kTiny > 0 {
			execs["kTiny"] += kTiny
			return
		}

		kVertDom := hypergraph.S_VertexDominationRule(gf, g, c, expand)
		if kVertDom > 0 {
			execs["kVertDom"] += kVertDom
			return
		}

		kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		if kEdgeDom > 0 {
			execs["kEdgeDom"] += kEdgeDom
			return
		}
	}

	kApVertDom := hypergraph.FS_ApproxVertexDominationRule(gf, g, c, expand)
	if kApVertDom > 0 {
		execs["kApVertDom"] += kApVertDom
		return
	}
	kApDoubleVertDom := hypergraph.FS_ApproxDoubleVertexDominationRule(gf, g, c, expand)
	if kApDoubleVertDom > 0 {
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		return
	}
	kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.FS_SmallEdgeDegreeTwoRule(gf, g, c, expand)
	if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
		return
	}
	kTri := hypergraph.FS_SmallTriangleRule(gf, g, c, expand)
	if kTri > 0 {
		execs["kTri"] += kTri
		return
	}
	kExtTri := hypergraph.FS_ExtendedTriangleRule(gf, g, c, expand)
	if kExtTri > 0 {
		execs["kExtTri"] += kExtTri
		return
	}
	kSmall := hypergraph.FS_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
	if kSmall > 0 {
		execs["kSmall"] += kSmall
		return
	}
}

func getGraphPres(useFrontier bool) *GraphPayload {
	var pg *hypergraph.HyperGraph

	if useFrontier {
		pg = gf
	} else {
		pg = g
	}

	graphPres := make(map[int32][]int32)
	for eId, e := range pg.Edges {
		graphPres[eId] = make([]int32, len(e.V))
		i := 0
		for v := range e.V {
			graphPres[eId][i] = v
			i++
		}
	}

	i := 0
	cArr := make([]int32, len(c))
	for v := range c {
		cArr[i] = v
		i++
	}

	i = 0
	vertices := make([]int32, len(pg.Vertices))
	for v := range pg.Vertices {
		vertices[i] = v
		i++
	}

	execsTable := make([]*ExecsTableRow, len(execs))

	for idx, label := range alg.Labels {
		execsTable[idx] = &ExecsTableRow{Name: short_names[label], K: execs[label]}
	}

	frArr := make([]int32, len(gf.Vertices))
	i = 0
	for v := range gf.Vertices {
		frArr[i] = v
		i++
	}

	return &GraphPayload{Graph: graphPres, C: cArr, ExecsData: execsTable, Vertices: vertices, Ratio: alg.GetRatio(execs), Opt: alg.GetEstOpt(execs), Frontier: frArr}
}
