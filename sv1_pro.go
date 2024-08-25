package main

import (
	"fmt"
	"math/rand"
	"os"
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


func (algo *SV1) CheckReachability(src string, dst string) (bool, error) {
	svLabel := algo.SV.Label()

	//if src is support vertex
	if svLabel == src {
		fmt.Println("[CheckReachability][Resolved] Src vertex is SV")
		return algo.R_Plus[dst], nil
	}

	//if dest is support vertex
	if svLabel == dst {
		fmt.Println("[CheckReachability][Resolved] Dst vertex is SV")
		return algo.R_Minus[dst], nil
	}

	//try to apply O1
	if algo.R_Minus[src] == true && algo.R_Plus[dst] == true {
		fmt.Println("[CheckReachability][Resolved] Using O1")
		return true, nil
	}

	//try to apply O2
	if algo.R_Plus[src] == true && algo.R_Plus[dst] == false {
		fmt.Println("[CheckReachability][Resolved] Using O2")
		return false, nil
	}

	//try to apply O3
	if algo.R_Minus[src] == false && algo.R_Minus[dst] == true {
		fmt.Println("[CheckReachability][Resolved] Using O3")
		return false, nil
	}

	//if all else fails, fallback to BFS
	fmt.Println("[CheckReachability][Resolved] Fallback to BFS")
	bfs, err := directedBFS(algo.Graph, src,dst)
	if err != nil {
		return false, err
	}
	// for bfs.HasNext() {
	// 	v := bfs.Next()
	// 	if v.Label() == dst {
	// 		return true, nil
	// 	}
	// }
	if bfs == true{
		return true,nil }
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
 
func main() {
 

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create a new directed graph
	graph := gograph.New[string](gograph.Directed())

	// Create vertices 
	vertexCount := 50
	vertices := make([]*gograph.Vertex[string], vertexCount)
	for i := 0; i < vertexCount; i++ {
		label := "v" + strconv.Itoa(i)
		vertices[i] = gograph.NewVertex(label)
		graph.AddVertex(vertices[i])
	}

	
	edgeCount := 100
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

	// Perform some reachability tests
	testCases := [][2]string{
		{"v0", "v10"},
		{"v10", "v20"},
		{"v20", "v30"},
		{"v30", "v40"},
		{"v40", "v0"},
		{"v5", "v25"},
		{"v25", "v5"},
		{"v48","v50"},
		{"v75","v32"},
		
	}

	fmt.Println("\nReachability Tests:")
	for _, test := range testCases {
		src, dst := test[0], test[1]
		reachable, err := index.CheckReachability(src, dst)
		if err != nil {
			fmt.Printf("Error checking reachability from %s to %s: %v\n", src, dst, err)
		} else {
			fmt.Printf("Is %s reachable from %s? %v\n", dst, src, reachable)
		}
	}
}
