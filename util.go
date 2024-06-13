package bbn

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Error type for conflicting evidence.
type ConflictingEvidenceError struct{}

func (m *ConflictingEvidenceError) Error() string {
	return "conflicting evidence / all samples rejected"
}

// helper for YAML serialization of a network.
type networkYaml struct {
	Name      string
	Variables []*Node
}

// FromFile reads a [Network] from an YAML or XML file.
func FromFile(path string) (*Network, []*Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	ext := filepath.Ext(path)

	switch strings.ToLower(ext) {
	case ".yml":
		return FromYAML(data)
	case ".xml", ".bifxml":
		return FromBIFXML(data)
	default:
		return nil, nil, fmt.Errorf("unsupported file format '%s'", ext)
	}
}

// toInternalNodes transforms nodes to their internal representation.
func toInternalNodes(nodes []*Node) ([]*node, error) {
	nodeMap := map[string]int{}
	for i, n := range nodes {
		nodeMap[n.Variable] = i
	}

	nodeList := make([]*node, len(nodes))
	for i, n := range nodes {
		nodeMap[n.Variable] = i

		parents := make([]int, len(n.Given))
		for j, p := range n.Given {
			par, ok := nodeMap[p]
			if !ok {
				return nil, fmt.Errorf("parent node '%s' not found", p)
			}
			parents[j] = par
		}

		var stride []int
		tableRows := 1
		if len(parents) > 0 {
			stride = make([]int, len(n.Given))
			stride[len(stride)-1] = 1
			for j := len(stride) - 2; j >= 0; j-- {
				parIdx := parents[j+1]
				stride[j] = stride[j+1] * len(nodes[parIdx].Outcomes)
			}
			tableRows = stride[0] * len(nodes[parents[0]].Outcomes)
		}

		if len(n.Table) != tableRows {
			return nil, fmt.Errorf("wrong number of table rows in node '%s'; got %d, expected %d", n.Variable, len(n.Table), tableRows)
		}

		tp, ok := nodeTypes[n.Type]
		if !ok {
			return nil, fmt.Errorf("unknown node type '%s' for '%s'", n.Type, n.Variable)
		}

		tableCols := len(n.Outcomes)
		if tp == UtilityNode {
			if tableCols != 1 {
				return nil, fmt.Errorf("utility node '%s' must have a single table column, got %d", n.Variable, tableCols)
			}
		} else {
			if tableCols < 2 {
				return nil, fmt.Errorf("node '%s' requires at least two outcomes, got %d", n.Variable, tableCols)
			}
		}

		for j, probs := range n.Table {
			if len(probs) != tableCols {
				return nil, fmt.Errorf("wrong number of table columns in node '%s', row %d; got %d, expected %d", n.Variable, j, len(probs), tableCols)
			}
		}

		nd := node{
			Variable:   n.Variable,
			Type:       tp,
			ID:         i,
			GivenNames: n.Given,
			Given:      parents,
			Stride:     stride,
			Outcomes:   n.Outcomes,
			Table:      n.Table,
			TableCum:   nil,
			Position:   n.Position,
		}

		nodeList[i] = &nd
	}

	if err := checkNodes(nodeList); err != nil {
		return nil, err
	}

	return nodeList, nil
}

func checkNodes(nodes []*node) error {
	for _, n := range nodes {
		if n.Type == DecisionNode && len(n.Given) > 0 {
			return fmt.Errorf("decision node '%s' can't have any parent nodes", n.Variable)
		}
		for _, parIdx := range n.Given {
			par := nodes[parIdx]
			if par.Type == UtilityNode {
				return fmt.Errorf("utility node '%s' can't be a parent of any other node", par.Variable)
			}
		}
	}
	return nil
}

func isAcyclic(nodes []*node) bool {
	for i := range nodes {
		if !isAcyclicRecursive(nodes, i, i) {
			return false
		}
	}
	return true
}

func isAcyclicRecursive(nodes []*node, index int, start int) bool {
	n := nodes[index]
	for _, parent := range n.Given {
		if parent == start || !isAcyclicRecursive(nodes, parent, start) {
			return false
		}
	}
	return true
}

// sortTopological sorts nodes in topological order.
func sortTopological(nodes []*node) ([]*node, error) {
	visited := make([]bool, len(nodes))
	stack := []int{}

	for i := range nodes {
		var err error
		stack, err = sortTopologicalRecursive(nodes, i, visited, stack)
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
		for i, par := range n.Given {
			n.Given[i] = newIndex[par]
		}
		result[i] = n
	}

	return result, nil
}

// sortTopologicalRecursive performs the recursion used in sortTopological.
func sortTopologicalRecursive(nodes []*node, index int, visited []bool, stack []int) ([]int, error) {
	if visited[index] {
		return stack, nil
	}

	visited[index] = true
	n := nodes[index]

	for _, parent := range n.Given {
		var err error
		stack, err = sortTopologicalRecursive(nodes, parent, visited, stack)
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

// Sample from cumulative (relative) probabilities.
func Sample(cum []float64, rng *rand.Rand) int {
	ln := len(cum)
	r := rng.Float64() * cum[ln-1]

	for i, v := range cum {
		if v >= r {
			return i
		}
	}
	return -1
}
