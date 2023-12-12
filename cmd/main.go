package main

import (
	"fmt"
	"os"

	"github.com/KhoalaS/BachelorThesis/pkg/alg"
	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)




func makeChart() {
	var baseSize int32 = 10
	baseSizes := []int32{}
	var g *hypergraph.HyperGraph
	var maxVert int32 = 10000

	labels := make([]int, 20)
	lineSeries := make(map[int32][]opts.LineData)
	barLabels := []string{"kTiny", "kEdgeDom", "kVertDom", "kTri", "kApVertDom", "kSmall", "kApDoubleVertDom"}
	barLabelsShort := []string{"Tiny", "EDom", "VDom", "Tri", "ApVDom", "Small", "ApDVDom"}

	barSeries1 := make(map[int32][]opts.BarData)
	barSeries10 := make(map[int32][]opts.BarData)

	for baseSize <= maxVert {
		baseSizes = append(baseSizes, baseSize)
		lineSeries[baseSize] = []opts.LineData{}
		for i := 1; i <= 20; i++ {
			labels[i-1] = i
			g = hypergraph.GenerateTestGraph(baseSize, int32(i)*baseSize, true)
			c := make(map[int32]bool)

			execs, _ := alg.ApplyRules(g, c, 1000000)

			var nom float64 = 0
			var denom float64 = 0

			for key, val := range execs {
				if key == "kEdgeDom" || key == "kVertDom"{
					continue
				}
				nom += float64(alg.Ratios[key].A * val)
				denom += float64(alg.Ratios[key].B * val)
			}
			fmt.Println("Nom: " ,nom)
			fmt.Println("Denom: " ,denom)



			lineSeries[baseSize] = append(lineSeries[baseSize], opts.LineData{Value: fmt.Sprintf("%.2f",(nom / denom))})

			if _, ex := barSeries10[baseSize]; !ex {
				barSeries10[baseSize] = []opts.BarData{}
				barSeries1[baseSize] = []opts.BarData{}
			}
			if i == 10 {
				for _, v := range barLabels {
					barSeries10[baseSize] = append(barSeries10[baseSize], opts.BarData{Value: execs[v]})
				}
			} else if i == 1 {
				for _, v := range barLabels {
					barSeries1[baseSize] = append(barSeries1[baseSize], opts.BarData{Value: execs[v]})
				}
			}
			fmt.Println("Edges/Vertices Factor:", i, "|", "Approximation Factor:", nom/denom)
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
		charts.WithYAxisOpts(opts.YAxis{Min: 1, Max: 2, Name: "est. Approximation Factor"}),
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
	bar10 := charts.NewBar()

	bar1.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "#Rule Executions",
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

	bar10.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "#Rule Executions",
			Subtitle: "[#Edges/#Vertices = 10]",
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

	bar10.SetXAxis(barLabelsShort).
		AddSeries("10 Vertices", barSeries10[10]).
		AddSeries("100 Vertices", barSeries10[100]).
		AddSeries("1K Vertices", barSeries10[1000]).
		AddSeries("10K Vertices", barSeries10[10000]).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	page.AddCharts(line, bar1, bar10)
	f, _ := os.Create("approx_factor_chart.html")
	page.Render(f)
}

func main() {
	g := hypergraph.NewHyperGraph()
	g.AddEdge(1,2,3)
	g.AddEdge(3,4)

	sol := alg.MinEdgeCover(g)
	fmt.Println(sol)

}
