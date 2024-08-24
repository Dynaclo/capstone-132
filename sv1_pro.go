package main

import (
	"github.com/hmdsefi/gograph"
	"github.com/hmdsefi/gograph/traverse"
	"math/rand"
)

// Fully Dynamic Transitive Closure Index
type FullDynTCIndex interface {
	NewIndex(graph gograph.Graph[string])

	insertEdge(src string, dst string) error
	deleteEdge(src string, dst string) error

	checkReachability(src string, dst string) (bool, error)
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

	//initialize R_Plus
	for _, v := range vertices {
		algo.R_Plus[v.Label()] = false
	}

	bfs, err := traverse.NewBreadthFirstIterator(algo.Graph, algo.SV.Label())
	if err != nil {
		panic(err)
	}
	bfs.Iterate(func(v *gograph.Vertex[string]) error {
		algo.R_Plus[v.Label()] = true
		return nil
	})
	//------------------

	//initialize R_Minus
	for _, v := range vertices {
		algo.R_Minus[v.Label()] = false
	}

	bfs_rev, err := traverse.NewBreadthFirstIterator(algo.ReverseGraph, algo.SV.Label())
	if err != nil {
		panic(err)
	}
	bfs_rev.Iterate(func(v *gograph.Vertex[string]) error {
		algo.R_Minus[v.Label()] = true
		return nil
	})
}

func (algo *SV1) insertEdge(src string, dst string) error {
	srcVertex := algo.Graph.GetVertexByID(src)
	if srcVertex == nil {
		srcVertex = gograph.NewVertex(src)
		algo.Graph.AddVertex(srcVertex)
	}

	dstVertex := algo.Graph.GetVertexByID(dst)
	if dstVertex == nil {
		dstVertex = gograph.NewVertex(dst)
		algo.Graph.AddVertex(dstVertex)
	}

	algo.Graph.AddEdge(srcVertex, dstVertex)
	algo.ReverseGraph.AddEdge(dstVertex, srcVertex)

	//TODO: UPDATE R+ AND R-
	return nil
}

func (algo *SV1) deleteEdge(src string, dst string) error {
	srcVertex := algo.Graph.GetVertexByID(src)
	dstVertex := algo.Graph.GetVertexByID(dst)

	edge := algo.Graph.GetEdge(srcVertex, dstVertex)
	algo.Graph.RemoveEdges(edge)

	rev_edge := algo.ReverseGraph.GetEdge(dstVertex, srcVertex)
	algo.ReverseGraph.RemoveEdges(rev_edge)
	//TODO: Add error handling here for if vertex or edge does not exist

	//TODO: UPDATE R+ AND R-
	return nil
}

func main() {
	graph := gograph.New[int](gograph.Directed())

	//bfs, err := traverse.NewBreadthFirstIterator[int](graph, 1)
	//if err != nil {
	//	panic(err)
	//}

	graph.GetAllVerticesByID()

}
