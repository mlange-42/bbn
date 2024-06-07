package bbn

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/yaml.v3"
)

func NodesFromYAML(path string) ([]*Node, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	nodes := []*Node{}
	err = decoder.Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// toInternalNodes transforms nodes to their internal representation.
func toInternalNodes(nodes []*Node) ([]*node, error) {
	nodeMap := map[string]int{}
	for i, n := range nodes {
		nodeMap[n.Name] = i
	}

	nodeList := make([]*node, len(nodes))
	for i, n := range nodes {
		nodeMap[n.Name] = i

		cum := make([][]float64, len(n.CPT))
		for j, probs := range n.CPT {
			cum[j] = cumulate(probs)
		}

		parents := make([]int, len(n.Parents))
		for j, p := range n.Parents {
			par, ok := nodeMap[p]
			if !ok {
				return nil, fmt.Errorf("parent node '%s' not found", p)
			}
			parents[j] = par
		}

		var stride []int
		if len(parents) > 0 {
			stride = make([]int, len(n.Parents))
			stride[len(stride)-1] = 1
			for j := len(stride) - 2; j >= 0; j-- {
				parIdx := parents[j+1]
				stride[j] = stride[j+1] * len(nodeList[parIdx].States)
			}
		}

		nd := node{
			Name:    n.Name,
			ID:      i,
			Parents: parents,
			Stride:  stride,
			States:  n.States,
			CPT:     n.CPT,
			CPTCum:  cum,
		}

		nodeList[i] = &nd
	}

	return nodeList, nil
}

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
