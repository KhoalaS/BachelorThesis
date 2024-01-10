package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func makeChart(pa float64, u int, evr int, maxv int, checkpoint int, fixRatio string) {
	var baseSize int32 = 10
	baseSizes := []int32{}
	var g *hypergraph.HyperGraph
	var maxest float64 = 2

	var maxVert int32 = 10000
	if maxv > 0 {
		maxVert = int32(maxv)
	}

	maxratio := 20
	if evr > 0 {
		maxratio = evr
	}

	if pa > 0 {
		maxratio = 1
	}

	if checkpoint > maxratio {
		checkpoint = maxratio
	}

	labels := make([]int, maxratio)
	lineSeries := make(map[int32][]opts.LineData)
	barLabels := []string{"kTiny", "kEdgeDom", "kVertDom", "kTri", "kExtTri", "kApVertDom", "kSmall", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kFallback"}
	barLabelsShort := []string{"Tiny", "EDom", "VDom", "Tri", "ETri", "ApVDom", "Small", "ApDVDom", "SmED2", "F3"}

	barSeries1 := make(map[int32][]opts.BarData)
	barSeries2 := make(map[int32][]opts.BarData)

	for baseSize <= maxVert {
		baseSizes = append(baseSizes, baseSize)
		lineSeries[baseSize] = []opts.LineData{}
		for i := 1; i <= maxratio; i++ {
			labels[i-1] = i
			if u > 0 {
				g = hypergraph.GenerateUniformTestGraph(baseSize, int32(i)*baseSize, u)
			} else if len(fixRatio) > 0 {
				spl := strings.Split(fixRatio, ",")
				ratios := make([]int, len(spl))
				for i, val := range spl {
					valInt, _ := strconv.Atoi(val)
					ratios[i] = valInt
				}
				g = hypergraph.GenerateFixDistTestGraph(baseSize, int32(i)*baseSize, ratios)
			} else if pa > 0 {
				g = hypergraph.GeneratePrefAttachmentGraph(baseSize, pa, 3)
			} else {
				g = hypergraph.GenerateTestGraph(baseSize, int32(i)*baseSize, true)
			}
			c := make(map[int32]bool)
			execs := make(map[string]int)

			alg.ThreeHS_F3ApprPoly(g, c, execs, 0, alg.Linear)
			var nom float64 = 0
			var denom float64 = 0

			for key, val := range execs {
				if key == "kEdgeDom" || key == "kVertDom" {
					continue
				}
				nom += float64(alg.Ratios[key].A * val)
				denom += float64(alg.Ratios[key].B * val)
			}

			//fmt.Println("Nom: " ,nom)
			//fmt.Println("Denom: " ,denom)
			if nom/denom > maxest {
				maxest = 3
			}

			lineSeries[baseSize] = append(lineSeries[baseSize], opts.LineData{Value: fmt.Sprintf("%.2f", (nom / denom))})

			if _, ex := barSeries2[baseSize]; !ex {
				barSeries2[baseSize] = []opts.BarData{}
				barSeries1[baseSize] = []opts.BarData{}
			}
			if i == checkpoint {
				for _, v := range barLabels {
					barSeries2[baseSize] = append(barSeries2[baseSize], opts.BarData{Value: execs[v]})
				}
			} else if i == 1 {
				for _, v := range barLabels {
					barSeries1[baseSize] = append(barSeries1[baseSize], opts.BarData{Value: execs[v]})
				}
			}
			fmt.Println(len(g.Edges), "Edges/Vertices Factor:", i, "|", "Approximation Factor:", nom/denom)
		}
		baseSize = baseSize * 10
	}

	page := components.NewPage()
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Save",
				},
			}},
		),
		charts.WithYAxisOpts(opts.YAxis{Min: 1, Max: maxest, Name: "est. Approximation Factor"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "#Edges\\#Vertices"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "80px"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	)
	line.SetXAxis(labels)

	for _, val := range baseSizes {
		l := fmt.Sprintf("%dK Vertices", val/1000)
		if val < 1000 {
			l = fmt.Sprintf("%d Vertices", val)
		}
		line.AddSeries(l, lineSeries[val]).SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: true,
			}),
		)
	}

	bar1 := charts.NewBar()
	bar2 := charts.NewBar()

	bar1.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{AxisLabel: &opts.AxisLabel{Rotate: 45, Show: true, ShowMinLabel: true, ShowMaxLabel: true}}),
		charts.WithTitleOpts(opts.Title{
			Title:    "#Rule Executions",
			Subtitle: "[#Edges/#Vertices = 1]",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "save",
				},
			}},
		),
	)
	bar2.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{AxisLabel: &opts.AxisLabel{Rotate: 45, Show: true, ShowMinLabel: true, ShowMaxLabel: true}}),
		charts.WithTitleOpts(opts.Title{
			Title:    "#Rule Executions",
			Subtitle: fmt.Sprintf("[#Edges/#Vertices = %d]", checkpoint),
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "save",
				},
			}},
		),
	)

	bar1.SetXAxis(barLabelsShort).
		AddSeries("10 Vertices", barSeries1[10]).
		AddSeries("100 Vertices", barSeries1[100]).
		AddSeries("1K Vertices", barSeries1[1000]).
		AddSeries("10K Vertices", barSeries1[10000]).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	bar2.SetXAxis(barLabelsShort).
		AddSeries("10 Vertices", barSeries2[10]).
		AddSeries("100 Vertices", barSeries2[100]).
		AddSeries("1K Vertices", barSeries2[1000]).
		AddSeries("10K Vertices", barSeries2[10000]).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	page.AddCharts(line, bar1, bar2)
	f, _ := os.Create("approx_factor_chart.html")
	page.Render(f)
}

func getRatio(execs map[string]int) float64 {
	var nom float64 = 0
	var denom float64 = 0

	for key, val := range execs {
		nom += float64(alg.Ratios[key].A * val)
		denom += float64(alg.Ratios[key].B * val)
	}
	return nom / denom
}

func main() {
	input := flag.String("i", "", "Filepath to graphml input file.")
	n := flag.Int("n", 10000, "Number of vertices if no graph file supplied.")
	m := flag.Int("m", 20000, "Number of edges if no graph file supplied.")
	K := flag.Int("k", 0, "The parameter k.")
	chart := flag.Bool("c", false, "Make charts.")
	u := flag.Int("u", 0, "Generate a u-uniform graph.")
	f := flag.String("f", "", "Generate a random hypergraph with fixed ratios for the edge sizes.")
	fr := flag.Int("fr", 0, "Preprocess the graph with fr many Factor-3 rule executions.")
	evr := flag.Int("evr", 0, "Maximum ratio |E|/|V| to compute for random graphs.")
	maxv := flag.Int("maxv", 0, "Maximum vertices for random graphs used in charts.")
	profile := flag.Bool("prof", false, "Make CPU profile")
	export := flag.String("o", "", "Export the generated graph with the given string as filename. The will create a 'graphs' folder where the file is located.")
	exportSimple := flag.String("os", "", "Export the generated graph to the given filepath.")
	prefAttach := flag.Float64("pa", 0.0, "Generate a random preferential attachment hypergraph with given float as probablity to add a new vertex.")
	prefAttachMod := flag.Bool("pamod", false, "")

	preset := flag.String("p", "", "Use a preconfigured chart preset. For available presets run with 'list -p'.")
	list := flag.NewFlagSet("list", flag.ExitOnError)
	printPreset := list.Bool("p", false, "")

	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "list" {
		list.Parse(os.Args[2:])
		if *printPreset {
			fmt.Println("u3\t 3-uniform graphs, E\\V ratio of 5, 1K maximum vertices")
			fmt.Println("u2\t 2-uniform graphs, E\\V ratio 10, 10K maximum vertices")
			fmt.Println("pa05\t preferential attachment graphs, probabilty 0.5 of attaching a new vertex, 10K maximum vertices")
			return
		}
	}

	if len(strings.Trim(*preset, " ")) > 0 {
		checkpoint := 10
		switch strings.Trim(*preset, " ") {
		case "u3":
			*u = 3
			*evr = 5
			*maxv = 1000
			checkpoint = 5
			*prefAttach = 0
		case "u2":
			*u = 2
			*evr = 10
			*maxv = 10000
			*prefAttach = 0
		case "pa05":
			*prefAttach = 0.5
			checkpoint = -1
		}
		makeChart(*prefAttach, *u, *evr, *maxv, checkpoint, *f)
		return
	}

	if *chart {
		makeChart(*prefAttach, *u, *evr, *maxv, 10, *f)
		return
	}

	if *K == 0 {
		*K = *n
	}

	k := *K

	var g *hypergraph.HyperGraph
	if len(strings.Trim(*input, " ")) > 0 {
		g = hypergraph.ReadFromFile(strings.Trim(*input, " "))
	} else if *u > 0 {
		g = hypergraph.GenerateUniformTestGraph(int32(*n), int32(*m), *u)
	} else if len(*f) > 0 {
		spl := strings.Split(*f, ",")
		ratios := make([]int, len(spl))
		for i, val := range spl {
			valInt, _ := strconv.Atoi(val)
			ratios[i] = valInt
		}
		g = hypergraph.GenerateFixDistTestGraph(int32(*n), int32(*m), ratios)
	} else if *prefAttach > 0 {
		g = hypergraph.GeneratePrefAttachmentGraph(int32(*n), 0.5, 3)
	} else if *prefAttachMod {
		g = hypergraph.GenerateModPrefAttachmentGraph(int(*n), 5, 0.5, 0.21)
	} else {
		g = hypergraph.GenerateTestGraph(int32(*n), int32(*m), true)
	}

	if len(*export) > 0 {
		hypergraph.WriteToFile(g, *export)
		return
	}

	if len(*exportSimple) > 0 {
		hypergraph.WriteToFileSimple(g, *exportSimple)
		return
	}

	c := make(map[int32]bool)
	execs := make(map[string]int)
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

	if *fr > 0 {
		// this might put all edges into the hitting set
		kFront := hypergraph.F3Prepocess(g, c, *fr)
		execs["kFallback"] = kFront
		k -= kFront
	}

	ex, hs, execs := alg.ThreeHS_F3ApprPoly(g, c, execs, prio, alg.Linear)

	pprof.StopCPUProfile()
	if ex {
		fmt.Printf("Found a 3-Hitting-Set of size %d\n", len(hs))
		fmt.Printf("Estimated Approximation Factor: %.2f\n", getRatio(execs))
		fmt.Println(execs)
	} else {
		fmt.Println("Did not find a 3-Hitting-Set")
	}
}
