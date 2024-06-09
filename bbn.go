package bbn

import (
	"fmt"
	"math/rand"
	"slices"
)

type networkDef struct {
	Name      string
	Variables []*Node
}

// Node definition.
//
// CPT is the conditional probability table.
// Each row represents the probabilities of the node's states for a certain
// combination of states of the nodes parents.
// Values in each row are relative, i.e. they do not necessarily sum up to 1.0.
// See the package examples.
type Node struct {
	Variable string      // Name of the node.
	Given    []string    `yaml:",flow"` // Names of parent nodes.
	Outcomes []string    `yaml:",flow"` // Names of the node's possible states.
	Table    [][]float64 `yaml:",flow"` // Conditional probability table.
	Position [2]int      `yaml:",flow"` // Coordinates for visualization, optional.
}

// node is the [Network]s internal node type.
type node struct {
	Variable   string
	ID         int
	GivenNames []string
	Given      []int
	Stride     []int
	Outcomes   []string
	Table      [][]float64
	TableCum   [][]float64
	Position   [2]int
}

func (n *node) Index(samples []int) int {
	idx := 0
	switch len(n.Given) {
	case 0:
		// Root nodes use the evidence is available.
		return 0
	case 1:
		// Optimized index calculation for one parent.
		idx = samples[n.Given[0]]
	case 2:
		// Optimized index calculation for two parents.
		idx = samples[n.Given[0]]*n.Stride[0] +
			samples[n.Given[1]]*n.Stride[1]
	default:
		// Default for more than 2 parents.
		for j, parIdx := range n.Given {
			parSample := samples[parIdx]
			idx += parSample * n.Stride[j]
		}
	}
	return idx
}

func (n *node) IndexWithNoData(samples []int) (int, bool) {
	idx := 0
	switch len(n.Given) {
	case 0:
		// Root nodes use the evidence is available.
		return 0, true
	case 1:
		// Optimized index calculation for one parent.
		g := n.Given[0]
		if samples[g] < 0 {
			return -1, false
		}
		return samples[g], true
	case 2:
		// Optimized index calculation for two parents.
		g1, g2 := n.Given[0], n.Given[1]
		if samples[g1] < 0 || samples[g2] < 0 {
			return -1, false
		}
		return samples[g1]*n.Stride[0] +
			samples[g2]*n.Stride[1], true
	default:
		// Default for more than 2 parents.
		for j, parIdx := range n.Given {
			parSample := samples[parIdx]
			if parSample < 0 {
				return -1, false
			}
			idx += parSample * n.Stride[j]
		}
		return idx, true
	}
}

// Network definition.
type Network struct {
	name   string
	nodes  []*node
	byName map[string]int
}

// New creates a new network. Sorts nodes topologically.
func New(name string, nodes ...*Node) (*Network, error) {
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
		byName[n.Variable] = i
		if len(n.Given) == 0 {
			continue
		}
	}

	network := Network{
		name:   name,
		nodes:  nodeList,
		byName: byName,
	}

	network.cumulateTables()

	return &network, nil
}

func (n *Network) cumulateTables() {
	for _, n := range n.nodes {
		cum := make([][]float64, len(n.Table))
		for j, probs := range n.Table {
			cum[j] = cumulate(probs)
		}
		n.TableCum = cum
	}
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
		result[node.Variable] = probs
	}

	return result, nil
}

// sample performs rejection sampling to calculate marginal probabilities of the network.
// Internal method working on prepared evidence and returning raw results.
func (n *Network) sample(ev []int, count int, rng *rand.Rand) ([][]int, int) {
	// Prepare slices for counting.
	counts := make([][]int, len(n.nodes))
	for i := range counts {
		counts[i] = make([]int, len(n.nodes[i].Outcomes))
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
			if len(node.Given) == 0 && e >= 0 {
				samples[i] = e
				continue
			}

			// Sample from cumulative probabilities.
			s := sample(node.TableCum[idx], rng)

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
		vIdx := slices.Index(n.nodes[idx].Outcomes, v)
		if vIdx < 0 {
			return nil, fmt.Errorf("value '%s' not found for node '%s' (has %s)", v, n.nodes[idx].Variable, n.nodes[idx].Outcomes)
		}
		ev[idx] = vIdx
	}
	return ev, nil
}
