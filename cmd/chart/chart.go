package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBoxPlotItems(boxPlotData [][]float64) []opts.BoxPlotData {
	items := make([]opts.BoxPlotData, 0)
	for i := 0; i < len(boxPlotData); i++ {
		items = append(items, opts.BoxPlotData{Value: boxPlotData[i]})
	}
	return items
}

// works in theory, but the data needs to be prepared first

func main() {
	in := flag.String("d", "./data", "path to folder with csv files")
	start := flag.Int("start", 1, "if placeholder is used in path, first value to be used")
	end := flag.Int("end", 1, "if placeholder is used in path, last value to be used")

	flag.Parse()

	ratios := [][]float64{}

	box := charts.NewBoxPlot()
	box.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "est. approximation factor"}),
		charts.WithYAxisOpts(opts.YAxis{Min: 1.8}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Save",
				},
			},
		}),
	)

	for i := *start; i <= *end; i++ {
		var data []fs.DirEntry
		var err error

		inrepl := strings.ReplaceAll(*in, "%", strconv.Itoa(i))
		data, err = os.ReadDir(inrepl)

		if err != nil {
			log.Fatal(err)
		}

		reg := regexp.MustCompile("master")
		for j := 0; j < len(data); j++ {
			if !reg.MatchString(data[j].Name()) {
				continue
			}

			//spl := strings.Split(data[i].Name(), "_")
			//evr, _ := strconv.ParseInt(spl[2],0,0)

			file, err := os.Open(inrepl + "/" + data[j].Name())
			if err != nil {
				log.Fatalf("Could not open file %s", inrepl+"/"+data[j].Name())
			}

			reader := csv.NewReader(file)
			reader.Comma = ';'

			reader.Read()
			ratios = append(ratios, []float64{})

			for {
				record, err := reader.Read()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				ratio, _ := strconv.ParseFloat(record[0], 64)
				ratios[i-*start] = append(ratios[i-*start], ratio)
			}
		}

	}
	fmt.Println(ratios)

	labels := []int{}
	for i:=*start; i<=*end; i++ {
		labels = append(labels, i)
	}

	box.SetXAxis(labels)
	box.AddSeries("Ratio", generateBoxPlotItems(ratios))
	outfile, err := os.OpenFile("./out/out.html", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Could not open ./out/out.html")
	}
	box.Render(outfile)
}
