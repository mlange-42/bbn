package bbn

import (
	"fmt"
	"math/rand"
	"slices"
)

// Node definition.
//
// CPT is the conditional probability table.
// Each row represents the probabilities of the node's states for a certain
// combination of states of the nodes parents.
// Values in each row are relative, i.e. they do not necessarily sum up to 1.0.
// See the package examples.
type Node struct {
	Name    string      // Name of the node.
	Parents []string    // Names of parent nodes.
	States  []string    // Names of the node's possible states.
	CPT     [][]float64 // Conditional probability table.
	Coords  [2]int      // Coordinates for visualization, optional.
}

// node is the [Network]s internal node type.
type node struct {
	Name    string
	ID      int
	Parents []int
	Stride  []int
	States  []string
	CPT     [][]float64
	CPTCum  [][]float64
}

// Network definition.
type Network struct {
	nodes  []*node
	byName map[string]int
}

// New creates a new network. Sorts nodes topologically.
func New(nodes ...*Node) (*Network, error) {
	nodeList, err := toInternalNodes(nodes)
	if err != nil {
		return nil, err
	}

	nodeList, err = sortTopological(nodeList)
	if err != nil {
		return nil, err
	}

	byName := map[string]int{}
	for i, n := range nodeList {
		byName[n.Name] = i
		if len(n.Parents) == 0 {
			continue
		}
	}

	return &Network{
		nodes:  nodeList,
		byName: byName,
	}, nil
}

// Sample performs rejection sampling to calculate marginal probabilities of the network.
func (n *Network) Sample(evidence map[string]string, count int, rng *rand.Rand) (map[string][]float64, error) {
	ev, err := n.prepareEvidence(evidence)
	if err != nil {
		return nil, err
	}

	counts := make([][]int, len(n.nodes))
	for i := range counts {
		counts[i] = make([]int, len(n.nodes[i].States))
	}

	samples := make([]int, len(n.nodes))
	matches := 0
	for r := 0; r < count; r++ {
		match := true
		for i, node := range n.nodes {
			idx := 0
			for j, parIdx := range node.Parents {
				parSample := samples[parIdx]
				idx += parSample * node.Stride[j]
			}
			s := sample(node.CPTCum[idx], rng)
			e := ev[i]
			if e >= 0 && e != s {
				match = false
				break
			}
			samples[i] = s
		}
		if match {
			for i, s := range samples {
				counts[i][s]++
			}
			matches++
		}
	}

	if matches == 0 {
		return nil, &ConflictingEvidenceError{}
	}

	result := map[string][]float64{}
	for i, node := range n.nodes {
		probs := make([]float64, len(counts[i]))
		for j, cnt := range counts[i] {
			probs[j] = float64(cnt) / float64(matches)
		}
		result[node.Name] = probs
	}

	return result, nil
}

// transforms the evidence map into an array with one entry per node.
// Missing evidence is indicated by -1.
func (n *Network) prepareEvidence(evidence map[string]string) ([]int, error) {
	ev := make([]int, len(n.nodes))
	for i := range ev {
		ev[i] = -1
	}

	for k, v := range evidence {
		idx, ok := n.byName[k]
		if !ok {
			return nil, fmt.Errorf("node '%s' not found", k)
		}
		vIdx := slices.Index(n.nodes[idx].States, v)
		if vIdx < 0 {
			return nil, fmt.Errorf("value '%s' not found for node '%s' (has %s)", v, n.nodes[idx].Name, n.nodes[idx].States)
		}
		ev[idx] = vIdx
	}
	return ev, nil
}
