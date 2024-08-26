package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	//"github.com/Smuzzy-waiii/capstone-132"
	"github.com/hmdsefi/gograph"
	//"github.com/hmdsefi/gograph/traverse"
)

// Fully Dynamic Transitive Closure Index
type FullDynTCIndex interface {
	NewIndex(graph gograph.Graph[string])

	InsertEdge(src string, dst string) error
	DeleteEdge(src string, dst string) error

	CheckReachability(src string, dst string) (bool, error)
}

type SV1 struct {
	Graph        gograph.Graph[string]
	ReverseGraph gograph.Graph[string]
	SV           *gograph.Vertex[string]
	R_Plus       map[string]bool //store all vertices reachable from SV
	R_Minus      map[string]bool //store all vertices that can reach SV
}

// maybe have a Init() fn return pointer to Graph object to which
// vertices are added instead of taking in graph as param which casues huge copy
// ok since it is a inti step tho ig
func (algo *SV1) NewIndex(graph gograph.Graph[string]) {
	algo.Graph = graph
	fmt.Printf("hiiii")
	print(algo.Graph)
	//make reverse DiGraph
	algo.ReverseGraph = gograph.New[string](gograph.Directed())
	for _, e := range algo.Graph.AllEdges() {
		algo.ReverseGraph.AddEdge(e.Destination(), e.Source())
	}

	//select support vertex
	//TODO: implement getting random vertex in the library itself
	vertices := algo.Graph.GetAllVertices()
	randomIndex := rand.Intn(len(vertices))
	algo.SV = vertices[randomIndex]

	//make sure this is not a isolated vertex and repick if it is
	for algo.SV.Degree() == 0 {
		randomIndex = rand.Intn(len(vertices))
		algo.SV = vertices[randomIndex]
	}
	fmt.Println(algo.SV.Label(), " chosen as SV")

	//initialize R_Plus
	algo.R_Plus = make(map[string]bool)
	for _, v := range vertices {
		algo.R_Plus[v.Label()] = false
	}
	algo.recomputeRPlus()

	//initialize R_Minus
	algo.R_Minus = make(map[string]bool)
	for _, v := range vertices {
		algo.R_Minus[v.Label()] = false
	}
	algo.recomputeRMinus()
}

// func (algo *SV1) recomputeRPlus() {
// 	bfs, err := traverse.NewBreadthFirstIterator(algo.Graph, algo.SV.Label())
// 	if err != nil {
// 		panic(err)
// 	}
// 	bfs.Iterate(func(v *gograph.Vertex[string]) error {
// 		algo.R_Plus[v.Label()] = true
// 		return nil
// 	})
// }

// func (algo *SV1) recomputeRMinus() {
// 	bfs_rev, err := traverse.NewBreadthFirstIterator(algo.ReverseGraph, algo.SV.Label())
// 	if err != nil {
// 		panic(err)
// 	}
// 	bfs_rev.Iterate(func(v *gograph.Vertex[string]) error {
// 		algo.R_Minus[v.Label()] = true
// 		return nil
// 	})
// }

func (algo *SV1) recomputeRPlus() {
	// Initialize a queue for BFS
	queue := []*gograph.Vertex[string]{algo.SV}

	// Reset R_Plus to mark all vertices as not reachable
	for key := range algo.R_Plus {
		algo.R_Plus[key] = false
	}

	// Start BFS
	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]

		algo.R_Plus[current.Label()] = true

		// Enqueue all neighbors (vertices connected by an outgoing edge)
		for _, edge := range algo.Graph.AllEdges() {
			if edge.Source().Label() == current.Label() {
				destVertex := edge.Destination()
				if !algo.R_Plus[destVertex.Label()] {
					queue = append(queue, destVertex)
				}
			}
		}
	}
}

func (algo *SV1) recomputeRMinus() {

	queue := []*gograph.Vertex[string]{algo.SV}

	for key := range algo.R_Minus {
		algo.R_Minus[key] = false
	}

	// Start BFS
	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]

		algo.R_Minus[current.Label()] = true

		for _, edge := range algo.ReverseGraph.AllEdges() {
			if edge.Source().Label() == current.Label() {
				destVertex := edge.Destination()
				if !algo.R_Minus[destVertex.Label()] {
					queue = append(queue, destVertex)
				}
			}
		}
	}
}

func (algo *SV1) InsertEdge(src string, dst string) error {
	srcVertex := algo.Graph.GetVertexByID(src)
	if srcVertex == nil {
		srcVertex = gograph.NewVertex(src)
		algo.Graph.AddVertex(srcVertex)
		algo.R_Plus[src] = false
		algo.R_Minus[src] = false
	}

	dstVertex := algo.Graph.GetVertexByID(dst)
	if dstVertex == nil {
		dstVertex = gograph.NewVertex(dst)
		algo.Graph.AddVertex(dstVertex)
		algo.R_Plus[dst] = false
		algo.R_Minus[dst] = false
	}

	algo.Graph.AddEdge(srcVertex, dstVertex)
	algo.ReverseGraph.AddEdge(dstVertex, srcVertex)

	//update R+ and R-
	//TODO: Make this not be a full recompute using an SSR data structure
	algo.recomputeRPlus()
	algo.recomputeRMinus()
	return nil
}

func (algo *SV1) DeleteEdge(src string, dst string) error {
	//TODO: Check if either src or dst are isolated after edgedelete and delete the node if they are not schema nodes
	//TODO: IF deleted node is SV or if SV gets isolated repick SV

	srcVertex := algo.Graph.GetVertexByID(src)
	dstVertex := algo.Graph.GetVertexByID(dst)

	edge := algo.Graph.GetEdge(srcVertex, dstVertex)
	algo.Graph.RemoveEdges(edge)

	rev_edge := algo.ReverseGraph.GetEdge(dstVertex, srcVertex)
	algo.ReverseGraph.RemoveEdges(rev_edge)
	//TODO: Add error handling here for if vertex or edge does not exist

	//TODO: Make this not be a full recompute using an SSR data structure
	algo.recomputeRPlus()
	algo.recomputeRMinus()
	return nil
}

// Directed BFS implementation
func directedBFS(graph gograph.Graph[string], src string, dst string) (bool, error) {

	queue := []*gograph.Vertex[string]{}

	visited := make(map[string]bool)

	// Start BFS from the source vertex
	startVertex := graph.GetVertexByID(src)
	if startVertex == nil {
		return false, fmt.Errorf("source vertex %s not found", src)
	}
	queue = append(queue, startVertex)
	visited[src] = true

	for len(queue) > 0 {
		// Dequeue the front of the queue
		currentVertex := queue[0]
		queue = queue[1:]

		// If we reach the destination vertex return true
		if currentVertex.Label() == dst {
			return true, nil
		}

		// Get all edges from the current vertex
		for _, edge := range graph.AllEdges() {
			// Check if the edge starts from the current vertex (directed edge)
			if edge.Source().Label() == currentVertex.Label() {
				nextVertex := edge.Destination()
				if !visited[nextVertex.Label()] {
					visited[nextVertex.Label()] = true
					queue = append(queue, nextVertex)
				}
			}
		}
	}

	// If we exhaust the queue without finding the destination, return false
	return false, nil
}
func bfs(graph *gograph.Graph[string], src, dst string) bool {
	visited := make(map[string]bool)
	queue := []*gograph.Vertex[string]{(*graph).GetVertexByID(src)}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Label() == dst {
			return true
		}

		if visited[current.Label()] {
			continue
		}
		visited[current.Label()] = true

		for _, edge := range (*graph).AllEdges() {
			if edge.Source().Label() == current.Label() {
				destVertex := edge.Destination()
				if !visited[destVertex.Label()] {
					queue = append(queue, destVertex)
				}
			}
		}
	}
	return false
}

// Directed BFS implementation without a specific destination
func directedBFSWithoutDestination(graph gograph.Graph[string], src string) ([]string, error) {
	// Queue for BFS
	queue := []*gograph.Vertex[string]{}

	visited := make(map[string]bool)

	reachableVertices := []string{}

	// Start BFS from the source vertex
	startVertex := graph.GetVertexByID(src)
	if startVertex == nil {
		return nil, fmt.Errorf("source vertex %s not found", src)
	}
	queue = append(queue, startVertex)
	visited[src] = true
	reachableVertices = append(reachableVertices, src)

	for len(queue) > 0 {
		// Dequeue the front of the queue
		currentVertex := queue[0]
		queue = queue[1:]

		// Get all edges from the current vertex
		for _, edge := range graph.AllEdges() {
			// Check if the edge starts from the current vertex (directed edge)
			if edge.Source().Label() == currentVertex.Label() {
				nextVertex := edge.Destination()
				if !visited[nextVertex.Label()] {
					visited[nextVertex.Label()] = true
					reachableVertices = append(reachableVertices, nextVertex.Label())
					queue = append(queue, nextVertex)
				}
			}
		}
	}

	return reachableVertices, nil
}

func (algo *SV1) CheckReachability(src string, dst string,resolutionCounts map[string]int) (bool, error) {
	svLabel := algo.SV.Label()

	//if src is support vertex
	if svLabel == src {
		//fmt.Println("[CheckReachability][Resolved] Src vertex is SV")
		
		return algo.R_Plus[dst], nil
	}

	//if dest is support vertex
	if svLabel == dst {
		//fmt.Println("[CheckReachability][Resolved] Dst vertex is SV")
		return algo.R_Minus[dst], nil
	}

	//try to apply O1
	if algo.R_Minus[src] == true && algo.R_Plus[dst] == true {
		resolutionCounts["o1"]+=1
		//fmt.Println("[CheckReachability][Resolved] Using O1")
		return true, nil
	}

	//try to apply O2
	if algo.R_Plus[src] == true && algo.R_Plus[dst] == false {
		resolutionCounts["o2"]+=1
		//fmt.Println("[CheckReachability][Resolved] Using O2")
		return false, nil
	}

	//try to apply O3
	if algo.R_Minus[src] == false && algo.R_Minus[dst] == true {
		resolutionCounts["o3"]+=1
		//fmt.Println("[CheckReachability][Resolved] Using O3")
		return false, nil
	}

	//if all else fails, fallback to BFS
	//fmt.Println("[CheckReachability][Resolved] Fallback to BFS")
	resolutionCounts["bfs"]+=1
	bfs, err := directedBFS(algo.Graph, src, dst)
	if err != nil {
		return false, err
	}
	// for bfs.HasNext() {
	// 	v := bfs.Next()
	// 	if v.Label() == dst {
	// 		return true, nil
	// 	}
	// }
	if bfs == true {
		return true, nil
	}
	return false, nil
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

func main() {

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create a new directed graph
	graph := gograph.New[string](gograph.Directed())
	
	resolutionCounts := map[string]int{
		"o1":  0,
		"o2":  0,
		"o3":   0, 
		"bfs": 0, 
	}
	// Create vertices
	vertexCount := 300
	vertices := make([]*gograph.Vertex[string], vertexCount)
	for i := 0; i < vertexCount; i++ {
		label := "v" + strconv.Itoa(i)
		vertices[i] = gograph.NewVertex(label)
		graph.AddVertex(vertices[i])
	}

	edgeCount := 2000
	for i := 0; i < edgeCount; i++ {
		srcIndex := rand.Intn(vertexCount)
		dstIndex := rand.Intn(vertexCount)

		// Ensure no self-loops
		for srcIndex == dstIndex {
			dstIndex = rand.Intn(vertexCount)
		}

		graph.AddEdge(vertices[srcIndex], vertices[dstIndex])
		fmt.Printf("Edge added: %s -> %s\n", vertices[srcIndex].Label(), vertices[dstIndex].Label())
	}

	//for visualization
	err := generateDotFile(graph, "graph_visualization.dot")
	if err != nil {
		fmt.Printf("Error generating dot file: %v\n", err)
		return
	}
	fmt.Println("Dot file generated successfully: graph_visualization.dot")

	// Initialize the SV1 transitive closure index
	index := SV1{}
	index.NewIndex(graph)

	testcases := generateTestCases(&graph, 11200)

	fmt.Println("\nReachability Tests Using the SV1 Index:")
	sv1Times := []DataPoint{}
	for _, test := range testcases {
		src, dst := test[0], test[1]

		startTime := time.Now()
		_, err := index.CheckReachability(src, dst,resolutionCounts)
		if err != nil {
			fmt.Printf("Error checking reachability from %s to %s: %v\n", src, dst, err)
		}
		endTime := time.Now()
		timeTaken := endTime.Sub(startTime)

		sv1Times = append(sv1Times, DataPoint{TimeOfEntry: startTime, Duration: timeTaken})
	}
	// Sort the data points by Duration
	sort.Slice(sv1Times, func(i, j int) bool {
		return sv1Times[i].Duration < sv1Times[j].Duration
	})
	printStatistics(sv1Times)

	fmt.Println("\nReachability Tests Using the BFS:")
	bfsTimes := []DataPoint{}
	for _, test := range testcases {
		src, dst := test[0], test[1]

		startTime := time.Now()
		_ = bfs(&graph, src, dst)
		endTime := time.Now()
		timeTaken := endTime.Sub(startTime)

		bfsTimes = append(bfsTimes, DataPoint{TimeOfEntry: startTime, Duration: timeTaken})
	}
	// Sort the data points by Duration
	sort.Slice(bfsTimes, func(i, j int) bool {
		return bfsTimes[i].Duration < bfsTimes[j].Duration
	})
	printStatistics(bfsTimes)

	plotResults(bfsTimes, sv1Times)
	makePieChart(resolutionCounts)
}

type DataPoint struct {
	TimeOfEntry time.Time
	Duration    time.Duration
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
// func createPieChart(resolutionCounts map[string]int) {
//     // create a new pie instance
//     pie := charts.NewPie()
//     pie.SetGlobalOptions(
//         charts.WithTitleOpts(
//             opts.Title{
//                 Title:    "Which policies were used",
//                 Subtitle: "O1,O2,O3,BFS",
//             },
//         ),
//     )
//     pie.SetSeriesOptions()
//     pie.AddSeries("Monthly revenue",
//         resolutionCounts.).
//         SetSeriesOptions(
//             charts.WithPieChartOpts(
//                 opts.PieChart{
//                     Radius: 200,
//                 },
//             ),
//             charts.WithLabelOpts(
//                 opts.Label{
//                     Show:      true,
//                     Formatter: "{b}: {c}",
//                 },
//             ),
//         )
//     f, _ := os.Create("pie.html")
//     _ = pie.Render(f)
// }
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


