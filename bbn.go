package bbn

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

type NodeType uint8

const (
	NatureNode NodeType = iota
	DecisionNode
	UtilityNode
)

const (
	NatureNodeType   = "nature"
	DecisionNodeType = "decision"
	UtilityNodeType  = "utility"
)

var nodeTypes = map[string]NodeType{
	"":               NatureNode,
	NatureNodeType:   NatureNode,
	DecisionNodeType: DecisionNode,
	UtilityNodeType:  UtilityNode,
}

var nodeTypeNames = map[NodeType]string{
	NatureNode:   "",
	DecisionNode: DecisionNodeType,
	UtilityNode:  UtilityNodeType,
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
	Type     string      `yaml:",omitempty"` // Type of the node [nature, decision, utility]
	Given    []string    `yaml:",flow"`      // Names of parent nodes.
	Outcomes []string    `yaml:",flow"`      // Names of the node's possible states.
	Table    [][]float64 `yaml:",flow"`      // Conditional probability table.
	Position [2]int      `yaml:",flow"`      // Coordinates for visualization, optional.
}

// node is the [Network]s internal node type.
type node struct {
	Variable   string
	Type       NodeType
	ID         int
	GivenNames []string
	Given      []int
	Stride     []int
	Outcomes   []string
	Table      [][]float64
	TableCum   [][]float64
	Position   [2]int
}

func (n *node) Index(sample []int) int {
	idx := 0
	switch len(n.Given) {
	case 0:
		// Root nodes use the evidence is available.
		return 0
	case 1:
		// Optimized index calculation for one parent.
		idx = sample[n.Given[0]]
	case 2:
		// Optimized index calculation for two parents.
		idx = sample[n.Given[0]]*n.Stride[0] +
			sample[n.Given[1]]*n.Stride[1]
	default:
		// Default for more than 2 parents.
		for j, parIdx := range n.Given {
			parSample := sample[parIdx]
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

	if !isAcyclic(nodeList) {
		return nil, fmt.Errorf("graph has cycles")
	}

	nodeList, err = sortTopological(nodeList)
	if err != nil {
		return nil, err
	}

	byName := map[string]int{}
	for i, n := range nodeList {
		byName[n.Variable] = i
	}

	network := Network{
		name:   name,
		nodes:  nodeList,
		byName: byName,
	}

	network.cumulateTables()

	return &network, nil
}

// Name of the network.
func (n *Network) Name() string {
	return n.name
}

// cumulates CPTs of all nodes
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
	probs, ok := n.sample(ev, count, rng)
	// Error on zero matches.
	if !ok {
		return nil, &ConflictingEvidenceError{}
	}

	// Normalize result and return it as map.
	result := map[string][]float64{}
	for i, node := range n.nodes {
		result[node.Variable] = probs[i]
	}

	return result, nil
}

// sample performs rejection sampling to calculate marginal probabilities of the network.
// Internal method working on prepared evidence and returning raw results.
func (n *Network) sample(ev []int, count int, rng *rand.Rand) ([][]float64, bool) {
	var savedCounts [][]float64
	var savedWeight float64
	maxUtilityIndex := -1
	maxUtility := math.Inf(-1)

	decisionNodes, decisionStride, decisionChoices := n.collectDecisionNodes(ev)
	for choice := 0; choice < decisionChoices; choice++ {
		decisions := make([]int, len(n.nodes))
		for i, idx := range decisionNodes {
			node := n.nodes[idx]
			selected := (choice / decisionStride[i]) % len(node.Outcomes)
			decisions[idx] = selected
		}

		// Prepare slices for counting.
		counts := n.prepareCounts()

		// Sampling.
		sample := make([]int, len(n.nodes))
		utilitySample := make([]float64, len(n.nodes))
		weight := 0.0
		for r := 0; r < count; r++ {
			// Sample nodes.
			w := n.sampleOnce(sample, utilitySample, decisions, ev, rng)
			// Count matching samples
			n.countMatchingSample(sample, utilitySample, ev, counts, w)
			weight += w
		}

		if weight == 0 {
			continue
		}

		sumUtility := 0.0
		for i, node := range n.nodes {
			if node.Type != UtilityNode {
				continue
			}
			sumUtility += counts[i][0] / weight
		}
		if sumUtility > maxUtility {
			maxUtility = sumUtility
			maxUtilityIndex = choice

			savedCounts = counts
			savedWeight = weight
		}
	}

	decisions := make([]int, len(n.nodes))
	for i, idx := range decisionNodes {
		node := n.nodes[idx]
		selected := (maxUtilityIndex / decisionStride[i]) % len(node.Outcomes)
		decisions[idx] = selected
	}

	// Error on zero matches.
	if savedWeight == 0 {
		return nil, false
	}

	for _, idx := range decisionNodes {
		savedCounts[idx][decisions[idx]] = savedWeight
	}

	// Normalize result.
	normalizeCounts(savedCounts, savedWeight)

	return savedCounts, true
}

func (n *Network) sampleOnce(
	sample []int,
	utilitySample []float64,
	decisions []int,
	evidence []int,
	rng *rand.Rand) float64 {

	weight := 1.0
	for i, node := range n.nodes {
		idx := node.Index(sample)
		e := evidence[i]

		if node.Type == UtilityNode {
			utilitySample[i] = node.Table[idx][0]
		} else if node.Type == DecisionNode {
			if e >= 0 {
				sample[i] = e
			} else {
				sample[i] = decisions[i]
			}
		} else {
			if len(node.Given) == 0 {
				// Don't sample for root nodes with given evidence
				if e >= 0 {
					sample[i] = e
				} else {
					sample[i] = Sample(node.TableCum[idx], rng)
				}
				continue
			}

			if e >= 0 {
				sample[i] = e
				cum := node.TableCum[idx]
				weight *= node.Table[idx][e] / cum[len(cum)-1]
			} else {
				sample[i] = Sample(node.TableCum[idx], rng)
			}
		}
	}
	return weight
}

func (n *Network) countMatchingSample(sample []int, utilitySample []float64, evidence []int, counts [][]float64, weight float64) {
	for i, s := range sample {
		node := n.nodes[i]
		switch node.Type {
		case UtilityNode:
			counts[i][0] += utilitySample[i] * weight
		case DecisionNode:
			if evidence[i] >= 0 {
				counts[i][s] += weight
			}
		case NatureNode:
			counts[i][s] += weight
		}
	}
}

func (n *Network) prepareCounts() [][]float64 {
	counts := make([][]float64, len(n.nodes))
	for i := range counts {
		counts[i] = make([]float64, len(n.nodes[i].Outcomes))
	}
	return counts
}

func normalizeCounts(counts [][]float64, weight float64) {
	for i := range counts {
		for j, cnt := range counts[i] {
			counts[i][j] = cnt / float64(weight)
		}
	}
}

func (n *Network) collectDecisionNodes(evidence []int) (nodes []int, stride []int, choices int) {
	choices = 1
	for i, node := range n.nodes {
		if node.Type == DecisionNode && evidence[i] < 0 {
			nodes = append(nodes, i)
			choices *= len(node.Outcomes)
		}
	}

	if len(nodes) > 0 {
		stride = make([]int, len(nodes))
		stride[len(stride)-1] = 1
		for j := len(stride) - 2; j >= 0; j-- {
			stride[j] = stride[j+1] * len(n.nodes[nodes[j+1]].Outcomes)
		}
	}

	return
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
