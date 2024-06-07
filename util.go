package bbn

import (
	"fmt"
	"math/rand"
)

// sortTopological sorts nodes in topological order.
func sortTopological(nodes []*node) ([]*node, error) {
	visited := make([]bool, len(nodes))
	stack := []int{}

	for i := range nodes {
		var err error
		stack, err = sortTopologicalRecursive(nodes, i, i, visited, stack)
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

// sortTopologicalRecursive performs the recursion used in sortTopological.
func sortTopologicalRecursive(nodes []*node, index int, start int, visited []bool, stack []int) ([]int, error) {
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
		stack, err = sortTopologicalRecursive(nodes, parent, start, visited, stack)
		if err != nil {
			return nil, err
		}
	}

	stack = append(stack, index)

	return stack, nil
}

func cumulate(values []float64) []float64 {
	c := make([]float64, len(values))
	c[0] = values[0]
	for k := 1; k < len(values); k++ {
		c[k] = c[k-1] + values[k]
	}
	return c
}

// sample from cumulative (relative) probabilities.
func sample(cum []float64, rng *rand.Rand) int {
	ln := len(cum)
	r := rng.Float64() * cum[ln-1]

	for i, v := range cum {
		if v >= r {
			return i
		}
	}
	return -1
}
