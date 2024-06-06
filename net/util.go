package net

func sortTopological(nodes []*node) ([]*node, error) {
	visited := make([]bool, len(nodes))
	stack := []int{}

	for i := range nodes {
		stack = sortTopologicalDFS(nodes, i, visited, stack)
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

func sortTopologicalDFS(nodes []*node, index int, visited []bool, stack []int) []int {
	if visited[index] {
		return stack
	}

	visited[index] = true
	n := nodes[index]

	for _, parent := range n.Parents {
		stack = sortTopologicalDFS(nodes, parent, visited, stack)
	}

	stack = append(stack, index)

	return stack
}
