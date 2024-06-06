package net

import "fmt"

// Node definition.
type Node struct {
	Name    string
	Parents []string
	States  []string
	CPT     [][]float64
}

type node struct {
	Name    string
	ID      int
	Parents []int
	States  []string
	CPT     [][]float64
}

// Network definition.
type Network struct {
	nodes []*node
}

// New creates a new network by sorting nodes topologically.
func New(nodes []*Node) (*Network, error) {
	nodeMap := map[string]int{}
	nodeList := make([]*node, len(nodes))
	for i, n := range nodes {
		nodeMap[n.Name] = i
		nodeList[i] = &node{
			Name:   n.Name,
			ID:     i,
			States: n.States,
			CPT:    n.CPT,
		}
	}
	for i, n := range nodes {
		nn := nodeList[i]
		nn.Parents = make([]int, len(n.Parents))
		for j, p := range n.Parents {
			par, ok := nodeMap[p]
			if !ok {
				return nil, fmt.Errorf("parent node '%s' not found", p)
			}
			nn.Parents[j] = par
		}
	}

	nodeList, err := sortTopological(nodeList)
	if err != nil {
		return nil, err
	}

	return &Network{
		nodes: nodeList,
	}, nil
}
