package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/KhoalaS/BachelorThesis/pkg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var ratios = map[string]pkg.IntTuple{
	"kTiny": {A:1, B:1},
	"kSmall": {A:2, B:1},
	"kTri": {A:3, B:2},
	"kApVertDom": {A:2, B:1},
	"kApDoubleVertDom": {A:2, B:1},
}

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {

	execs := make(map[string]int)

	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kEdom := hypergraph.EdgeDominationRule(g, c)
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		kApVertDom := hypergraph.ApproxVertexDominationRule3(g, c)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		//kApDoubleVertDom := 0

		execs["kTiny"] += kTiny
		execs["kTri"] += kTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdom == 0 {
			break
		}
	}

	m, err := os.Create("mem_main.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer m.Close() // error handling omitted for example
	//runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(m); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	return execs
}

func makeChart() {
	var baseSize int32 = 10
	var g *hypergraph.HyperGraph

	labels := make([]int, 20)
	lineSeries := make(map[int32][]opts.LineData)

	for baseSize <= 10000 {
		lineSeries[baseSize] = []opts.LineData{}
		for i := 1; i <= 20; i++ {
			labels[i-1] = i
			g = hypergraph.GenerateTestGraph(baseSize, int32(i)*baseSize, false)
			g.RemoveDuplicate()
			c := make(map[int32]bool)

			execs := ApplyRules(g, c)

			var nom float64 = 0
			var denom float64 = 0

			for key, val := range execs {
				nom += float64(ratios[key].A * val)
				denom += float64(ratios[key].B * val)
			}

			lineSeries[baseSize] = append(lineSeries[baseSize], opts.LineData{Value: nom / denom})

			fmt.Println("Edges/Vertices Factor:", i, "|", "Approximation Factor:", nom/denom)
		}
		baseSize = baseSize * 10
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{Min: 1, Max: 2, Name: "est. Approximation Factor"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "#Edges\\#Vertices"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "80px"}),
	)
	line.SetXAxis(labels)

	for key, val := range lineSeries {
		l := fmt.Sprintf("%dK Vertices", key/1000)
		if key < 1000 {
			l = fmt.Sprintf("%d Vertices", key)
		}
		line.AddSeries(l, val).SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: true,
			}),
		)
	}
	f, _ := os.Create("approx_factor_chart.html")
	line.Render(f)
}

func main() {
	makeChart()
}
