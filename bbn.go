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

func (n *node) Index(samples []int) int {
	idx := 0
	switch len(n.Parents) {
	case 0:
		// Root nodes use the evidence is available.
		return 0
	case 1:
		// Optimized index calculation for one parent.
		idx = samples[n.Parents[0]]
	case 2:
		// Optimized index calculation for two parents.
		idx = samples[n.Parents[0]]*n.Stride[0] +
			samples[n.Parents[1]]*n.Stride[1]
	default:
		// Default for more than 2 parents.
		for j, parIdx := range n.Parents {
			parSample := samples[parIdx]
			idx += parSample * n.Stride[j]
		}
	}
	return idx
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
	// Evidence map to int slice.
	ev, err := n.prepareEvidence(evidence)
	if err != nil {
		return nil, err
	}

	// Do the actual sampling.
	counts, matches := n.sample(ev, count, rng)

	// Error on zero matches.
	if matches == 0 {
		return nil, &ConflictingEvidenceError{}
	}

	// Normalize result and return it as map.
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

// sample performs rejection sampling to calculate marginal probabilities of the network.
// Internal method working on prepared evidence and returning raw results.
func (n *Network) sample(ev []int, count int, rng *rand.Rand) ([][]int, int) {
	// Prepare slices for counting.
	counts := make([][]int, len(n.nodes))
	for i := range counts {
		counts[i] = make([]int, len(n.nodes[i].States))
	}

	// Sampling.
	samples := make([]int, len(n.nodes))
	matches := 0
	for r := 0; r < count; r++ {
		// Sample nodes.
		match := true
		for i, node := range n.nodes {
			idx := node.Index(samples)

			e := ev[i]
			// Don't sample for root nodes with given evidence
			if len(node.Parents) == 0 && e >= 0 {
				samples[i] = e
				continue
			}

			// Sample from cumulative probabilities.
			s := sample(node.CPTCum[idx], rng)

			// Reject if sample is not equal to evidence
			if e >= 0 && e != s {
				match = false
				break
			}

			// Otherwise, fill in the sample.
			samples[i] = s
		}

		// Count matching samples
		if match {
			for i, s := range samples {
				counts[i][s]++
			}
			matches++
		}
	}

	return counts, matches
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
