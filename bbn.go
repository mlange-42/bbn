package bbn

import (
	"fmt"
	"math/rand"
	"slices"
)

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

// New creates a new network by sorting nodes topologically.
func New(nodes ...*Node) (*Network, error) {
	nodeMap := map[string]int{}
	nodeList := make([]*node, len(nodes))
	for i, n := range nodes {
		nodeMap[n.Name] = i

		cum := make([][]float64, len(n.CPT))
		for j, probs := range n.CPT {
			c := make([]float64, len(probs))
			c[0] = probs[0]
			for k := 1; k < len(probs); k++ {
				c[k] = c[k-1] + probs[k]
			}
			cum[j] = c
		}

		nodeList[i] = &node{
			Name:   n.Name,
			ID:     i,
			States: n.States,
			CPT:    n.CPT,
			CPTCum: cum,
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

	byName := map[string]int{}
	for i, n := range nodeList {
		byName[n.Name] = i
		if len(n.Parents) == 0 {
			continue
		}

		stride := make([]int, len(n.Parents))
		stride[len(stride)-1] = 1
		for j := len(stride) - 2; j >= 0; j-- {
			parIdx := n.Parents[j+1]
			stride[j] = stride[j+1] * len(nodeList[parIdx].States)
		}

		n.Stride = stride
	}

	return &Network{
		nodes:  nodeList,
		byName: byName,
	}, nil
}

func (n *Network) Sample(evidence map[string]string, count int, rng *rand.Rand) (map[string][]float64, error) {
	ev, err := n.prepareEvidence(evidence)
	if err != nil {
		return nil, err
	}
	anyEvidence := len(evidence) > 0

	counts := make([][]int, len(n.nodes))
	for i := range counts {
		counts[i] = make([]int, len(n.nodes[i].States))
	}

	samples := make([]int, len(n.nodes))
	matches := 0
	for r := 0; r < count; r++ {
		for i, node := range n.nodes {
			idx := 0
			for j, parIdx := range node.Parents {
				parSample := samples[parIdx]
				idx += parSample * node.Stride[j]
			}
			samples[i] = sample(node.CPTCum[idx], rng)
		}
		match := true
		if anyEvidence {
			for j, e := range ev {
				if e >= 0 && e != samples[j] {
					match = false
					break
				}
			}
		}
		if match {
			for i, s := range samples {
				counts[i][s]++
			}
			matches++
		}
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
