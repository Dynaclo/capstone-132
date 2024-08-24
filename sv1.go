package main

import (
	"fmt"
	"math/rand"
	"time"
	"container/list"
)


type Graph struct {
	edges map[int]map[int]bool
}

 
func NewGraph() *Graph {
	return &Graph{edges: make(map[int]map[int]bool)}
}

 
func (g *Graph) AddEdge(u, v int) {
	if _, exists := g.edges[u]; !exists {
		g.edges[u] = make(map[int]bool)
	}
	g.edges[u][v] = true
}


func (g *Graph) DFS(start int, visited map[int]bool, reachable *map[int]bool) {
	visited[start] = true
	(*reachable)[start] = true
	for neighbor := range g.edges[start] {
		if !visited[neighbor] {
			g.DFS(neighbor, visited, reachable)
		}
	}
}

// ComputeReachability computes R+ (reachable from each node) and R- (reaches each node)
func (g *Graph) ComputeReachability(numNodes int) (map[int]map[int]bool, map[int]map[int]bool) {
	rPlus := make(map[int]map[int]bool)
	rMinus := make(map[int]map[int]bool)

	for i := 0; i < numNodes; i++ {
		rPlus[i] = make(map[int]bool)
		rMinus[i] = make(map[int]bool)
	}

	for node := range g.edges {
		visited := make(map[int]bool)
		reachable := make(map[int]bool)
		g.DFS(node, visited, &reachable)
		for k := range reachable {
			rPlus[node][k] = true
			rMinus[k][node] = true
		}
	}

	return rPlus, rMinus
}

// BFS performs Breadth-First Search to find if a path exists
func (g *Graph) BFS(start, end int) bool {
	if start == end {
		return true
	}
	visited := make(map[int]bool)
	queue := list.New()
	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(int)
		for neighbor := range g.edges[current] {
			if neighbor == end {
				return true
			}
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushBack(neighbor)
			}
		}
	}
	return false
}

// CheckPath determines if a path exists between two nodes using various methods
func (g *Graph) CheckPath(start, end int, rPlus, rMinus map[int]map[int]bool) bool {
	// Check using R+ and R-
	if rPlus[start][end] || rMinus[end][start] {
		return true
	}

	// Apply additional policies if R+ and R- do not provide an answer
	// Placeholder for three policies (e.g., heuristics)
	// if applyPolicy1(g, start, end) || applyPolicy2(g, start, end) || applyPolicy3(g, start, end) {
	// 	return true
	// }

	// Fallback to BFS if no policies provide an answer
	return g.BFS(start, end)
}

// Placeholder policy functions
func applyPolicy1(g *Graph, start, end int) bool {
	// Implement policy 1
	return false
}

func applyPolicy2(g *Graph, start, end int) bool {
	// Implement policy 2
	return false
}

func applyPolicy3(g *Graph, start, end int) bool {
	// Implement policy 3
	return false
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numNodes := 100
	numEdges := 500 // Adjust the number of edges for sparser or denser graph
	graph := NewGraph()

	// Create a sparser graph with fewer edges
	for i := 0; i < numEdges; i++ {
		// u := rand.Intn(numNodes)
		// v := rand.Intn(numNodes)
		u:=3
		v:=4

		if u != v {
			graph.AddEdge(3, 4)
			graph.AddEdge(5,45)
			graph.AddEdge(23,33)

		}

	}

	// Compute R+ and R-
	rPlus, rMinus := graph.ComputeReachability(numNodes)

	// Example query
	startNode := 0
	endNode := 3
	canReach := graph.CheckPath(startNode, endNode, rPlus, rMinus)
	fmt.Printf("Can node %d reach node %d? %v\n", startNode, endNode, canReach)
}
