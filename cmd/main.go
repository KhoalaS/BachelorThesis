package main

import (
	"fmt"
	"os"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool) map[float64]int{

	execs := make(map[float64]int)

	for{
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kEdom := hypergraph.EdgeDominationRule(g, c)	
		kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		kApVertDom :=  hypergraph.ApproxVertexDominationRule3(g, c)
		kApDoubleVertDom :=  hypergraph.ApproxDoubleVertexDominationRule(g, c)
		//kApDoubleVertDom := 0

		execs[1] += kTiny
		execs[1.5] += kTri
		execs[2] += kSmall
		execs[2] += kApVertDom
		execs[2] += kApDoubleVertDom

		if kTiny + kTri + kSmall + kApVertDom + kApDoubleVertDom + kEdom == 0 {
			break
		}
	}

	return execs
}

func main(){
	var baseSize int32 = 1000
	var g *hypergraph.HyperGraph

	labels := make([]int, 20)
	lineSeries := make(map[int32][]opts.LineData)

	//g.RemoveDuplicate()
	for baseSize <= 16000{
		lineSeries[baseSize] = []opts.LineData{}
		for i:=1 ; i<=20; i++ {
			labels[i-1] = i
			g = hypergraph.GenerateTestGraph(baseSize, int32(i)*baseSize, false)
			c := make(map[int32]bool)
	
			execs := ApplyRules(g, c)
	
			var nom float64 = 0
			var denom float64 = 0
	
			for key, val := range execs {
				nom += key * float64(val)
				denom += float64(val)
			}
	
			lineSeries[baseSize] = append(lineSeries[baseSize], opts.LineData{Value: nom/denom})
	
			fmt.Println("Edges/Vertices Factor:",i,"|", "Approximation Factor:", nom/denom)
		}
		baseSize = baseSize*2
	}

	
	
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithYAxisOpts(opts.YAxis{Min: 1, Max: 2, Name: "est. Approximation Factor"}),
			charts.WithXAxisOpts(opts.XAxis{Name: "#Edges\\#Vertices"}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
			charts.WithLegendOpts(opts.Legend{Show: true, Right: "80px"}),
		)
		line.SetXAxis(labels)
	
		for key, val := range lineSeries{
			line.AddSeries(fmt.Sprintf("%dK Vertices", key/1000), val).SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					ShowSymbol: true,
				}),
				
			)
		}
		f, _ := os.Create("approx_factor_chart.html")
		line.Render(f)
}