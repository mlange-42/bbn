package net

import "fmt"

func sortTopological(nodes []*node) ([]*node, error) {
	visited := make([]bool, len(nodes))
	stack := []int{}

	for i := range nodes {
		var err error
		stack, err = sortTopologicalDFS(nodes, i, i, visited, stack)
		if err != nil {
			return nil, err
		}
	}

	newIndex := make([]int, len(nodes))
	for i, idx := range stack {
		newIndex[idx] = i
	}

	result := make([]*node, len(nodes))
	for i, idx := range stack {
		n := nodes[idx]
		for i, par := range n.Parents {
			n.Parents[i] = newIndex[par]
		}
		result[i] = n
	}

	return result, nil
}

func sortTopologicalDFS(nodes []*node, index int, start int, visited []bool, stack []int) ([]int, error) {
	if visited[index] {
		return stack, nil
	}

	visited[index] = true
	n := nodes[index]

	for _, parent := range n.Parents {
		if parent == start {
			return nil, fmt.Errorf("graph has cycles")
		}
		var err error
		stack, err = sortTopologicalDFS(nodes, parent, start, visited, stack)
		if err != nil {
			return nil, err
		}
	}

	stack = append(stack, index)

	return stack, nil
}
