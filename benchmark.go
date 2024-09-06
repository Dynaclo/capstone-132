package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/hmdsefi/gograph"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"os"
	"time"
)

type DataPoint struct {
	TimeOfEntry time.Time
	Duration    time.Duration
}

func makePieChart(data map[string]int) {
	// Create a new pie chart instance
	pie := charts.NewPie()

	// Convert data map to slices for plotting
	var items []opts.PieData
	for label, value := range data {
		items = append(items, opts.PieData{
			Name:  label,
			Value: float64(value),
		})
	}

	// Set pie chart options
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Pie Chart",
		}),
	)
	pie.AddSeries("Pie Chart", items)

	// Save the chart to a file
	f, err := os.Create("piechart.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pie.Render(f); err != nil {
		panic(err)
	}

	fmt.Println("Pie chart saved as piechart.html")
}

func generateDotFile(graph gograph.Graph[string], filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("digraph G {\n")
	if err != nil {
		return err
	}

	for _, edge := range graph.AllEdges() {
		line := fmt.Sprintf("  \"%s\" -> \"%s\";\n", edge.Source().Label(), edge.Destination().Label())
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString("}\n")
	return err
}

func generateTestCases(graph *gograph.Graph[string], numCases int) [][2]string {
	rand.Seed(time.Now().UnixNano())
	vertices := (*graph).GetAllVertices()
	vertexCount := len(vertices)

	testCases := make([][2]string, numCases)
	for i := 0; i < numCases; i++ {
		srcIndex := rand.Intn(vertexCount)
		dstIndex := rand.Intn(vertexCount)

		testCases[i] = [2]string{
			vertices[srcIndex].Label(),
			vertices[dstIndex].Label(),
		}
	}

	return testCases
}

func plotResults(bfsTimes, sv1Times []DataPoint) {
	p := plot.New()

	p.Title.Text = "BFS vs SV1 Performance"
	p.X.Label.Text = "Test Cases"
	p.Y.Label.Text = "Time (ns)"

	// Set log scales for both axes
	p.Y.Scale = plot.LogScale{}

	pts1 := make(plotter.XYs, len(bfsTimes))
	pts2 := make(plotter.XYs, len(sv1Times))

	for i := range bfsTimes {
		pts1[i].X = float64(i + 1)
		pts1[i].Y = float64(bfsTimes[i].Duration)
	}

	for i := range sv1Times {
		pts2[i].X = float64(i + 1)
		pts2[i].Y = float64(sv1Times[i].Duration)
	}

	err := plotutil.AddLinePoints(p, "BFS", pts1, "SV1", pts2)
	if err != nil {
		panic(err)
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "performance.png"); err != nil {
		panic(err)
	}
}

func printStatistics(dataPoints []DataPoint) {
	n := len(dataPoints)
	if n == 0 {
		fmt.Println("No data points available")
		return
	}

	min_ := dataPoints[0].Duration
	max_ := dataPoints[n-1].Duration
	mean := calculateMean(dataPoints)
	median := calculateMedian(dataPoints)
	q1, q3 := calculateQuartiles(dataPoints)

	fmt.Println("Statistics:")
	fmt.Printf("Min: %v\n", min_)
	fmt.Printf("Max: %v\n", max_)
	fmt.Printf("Mean: %v\n", mean)
	fmt.Printf("Median: %v\n", median)
	fmt.Printf("Q1 (25th percentile): %v\n", q1)
	fmt.Printf("Q3 (75th percentile): %v\n", q3)
	fmt.Printf("99th percentile: %v\n", calculatePercentile(dataPoints, 99))
	fmt.Printf("99.9th percentile: %v\n", calculatePercentile(dataPoints, 99.9))
}

func calculateMean(dataPoints []DataPoint) time.Duration {
	total := time.Duration(0)
	for _, dp := range dataPoints {
		total += dp.Duration
	}
	return total / time.Duration(len(dataPoints))
}

func calculateMedian(dataPoints []DataPoint) time.Duration {
	n := len(dataPoints)
	if n%2 == 0 {
		return (dataPoints[n/2-1].Duration + dataPoints[n/2].Duration) / 2
	}
	return dataPoints[n/2].Duration
}

func calculateQuartiles(dataPoints []DataPoint) (q1, q3 time.Duration) {
	q1 = calculatePercentile(dataPoints, 25)
	q3 = calculatePercentile(dataPoints, 75)
	return
}

func calculatePercentile(dataPoints []DataPoint, percentile float64) time.Duration {
	index := float64(len(dataPoints)) * percentile / 100.0
	if index == float64(int(index)) {
		i := int(index)
		return (dataPoints[i-1].Duration + dataPoints[i].Duration) / 2
	}
	return dataPoints[int(math.Ceil(index))-1].Duration
}
