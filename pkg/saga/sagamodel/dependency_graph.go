package sagamodel

import (
	"fmt"
)

type dependencyGraph struct {
	vertices map[string]struct{}
	edges    map[string]map[string]struct{}
}

func buildDependencyGraph(defs []*TaskDefinition) (*dependencyGraph, error) {
	g := &dependencyGraph{
		vertices: make(map[string]struct{}),
		edges:    make(map[string]map[string]struct{}),
	}

	for _, def := range defs {
		if err := g.addVertex(def.Id); err != nil {
			return nil, fmt.Errorf("add vertex: %s: %w", def.Id, err)
		}
	}

	for _, def := range defs {
		for _, dep := range def.Dependencies {
			if err := g.addEdge(def.Id, dep); err != nil {
				return nil, fmt.Errorf("add edge: %s->%s: %w",
					def.Id, dep, err)
			}
		}
	}

	return g, nil
}

func (g *dependencyGraph) addVertex(id string) error {
	if _, ok := g.vertices[id]; ok {
		return ErrVertexAlreadyAdded
	}

	g.vertices[id] = struct{}{}
	g.edges[id] = make(map[string]struct{})

	return nil
}

func (g *dependencyGraph) addEdge(from, to string) error {
	if _, ok := g.vertices[from]; !ok {
		return fmt.Errorf("%w: %s", ErrVertexNotFound, from)
	}
	if _, ok := g.vertices[to]; !ok {
		return fmt.Errorf("%w: %s", ErrVertexNotFound, to)
	}

	g.edges[from][to] = struct{}{}

	return nil
}

func (g *dependencyGraph) hasCycles() bool {
	checked := make(map[string]bool, len(g.vertices))
	inStack := make(map[string]bool, len(checked))

	var dfs func(string) bool
	dfs = func(from string) bool {
		if _, ok := checked[from]; ok {
			return false
		}
		defer func() { checked[from] = true }()

		inStack[from] = true
		defer func() { inStack[from] = false }()

		for to := range g.edges[from] {
			if inStack[to] {
				return true
			}
			if dfs(to) {
				return true
			}
		}

		return false
	}

	for id := range g.vertices {
		if dfs(id) {
			return true
		}
	}

	return false
}
