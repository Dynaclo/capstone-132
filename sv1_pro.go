package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

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

func (algo *SV1) CheckReachability(src string, dst string, resolutionCounts map[string]int, debugPrint bool) (bool, error) {
	svLabel := algo.SV.Label()

	//if src is support vertex
	if svLabel == src {
		resolutionCounts["src"] += 1
		if debugPrint {
			fmt.Println("[CheckReachability][Resolved] Src vertex is SV")
		}
		return algo.R_Plus[dst], nil
	}

	//if dest is support vertex
	if svLabel == dst {
		resolutionCounts["dst"] += 1
		if debugPrint {
			fmt.Println("[CheckReachability][Resolved] Dst vertex is SV")
		}
		return algo.R_Minus[dst], nil
	}

	//try to apply O1
	if algo.R_Minus[src] == true && algo.R_Plus[dst] == true {
		resolutionCounts["o1"] += 1
		if debugPrint {
			fmt.Println("[CheckReachability][Resolved] Using O1")
		}
		return true, nil
	}

	//try to apply O2
	if algo.R_Plus[src] == true && algo.R_Plus[dst] == false {
		resolutionCounts["o2"] += 1
		if debugPrint {
			fmt.Println("[CheckReachability][Resolved] Using O2")
		}
		return false, nil
	}

	//try to apply O3
	if algo.R_Minus[src] == false && algo.R_Minus[dst] == true {
		resolutionCounts["o3"] += 1
		if debugPrint {
			fmt.Println("[CheckReachability][Resolved] Using O3")
		}
		return false, nil
	}

	//if all else fails, fallback to BFS
	if debugPrint {
		fmt.Println("[CheckReachability][Resolved] Fallback to BFS")
	}
	resolutionCounts["bfs"] += 1
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

func main() {

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create a new directed graph
	graph := gograph.New[string](gograph.Directed())

	resolutionCounts := map[string]int{
		"src": 0,
		"dst": 0,
		"o1":  0,
		"o2":  0,
		"o3":  0,
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

	edgeCount := 1000
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

	testcases := generateTestCases(&graph, 1000)

	sv1Times := []DataPoint{}
	bfsTimes := []DataPoint{}

	for i, test := range testcases {
		src, dst := test[0], test[1]

		startTime := time.Now()
		isReachableSV1, err := index.CheckReachability(src, dst, resolutionCounts, false)
		if err != nil {
			fmt.Printf("Error checking reachability from %s to %s: %v\n", src, dst, err)
		}
		endTime := time.Now()
		timeTaken := endTime.Sub(startTime)
		sv1Times = append(sv1Times, DataPoint{TimeOfEntry: startTime, Duration: timeTaken, IsReachable: isReachableSV1})

		startTime = time.Now()
		isReachableBFS := bfs(&graph, src, dst)
		endTime = time.Now()
		timeTaken = endTime.Sub(startTime)
		bfsTimes = append(bfsTimes, DataPoint{TimeOfEntry: startTime, Duration: timeTaken, IsReachable: isReachableBFS})

		fmt.Printf("Testcase %d: %t\n", i, isReachableSV1 == isReachableBFS)
	}

	// Sort the data points by Duration
	sort.Slice(sv1Times, func(i, j int) bool {
		return sv1Times[i].Duration < sv1Times[j].Duration
	})
	fmt.Println("\nReachability Tests Using the SV1 Index:")
	printStatistics(sv1Times)

	// Sort the data points by Duration
	sort.Slice(bfsTimes, func(i, j int) bool {
		return bfsTimes[i].Duration < bfsTimes[j].Duration
	})
	fmt.Println("\nReachability Tests Using the BFS:")
	printStatistics(bfsTimes)

	plotResults(bfsTimes, sv1Times)
	makePieChart(resolutionCounts)

	fmt.Print(resolutionCounts)
}
