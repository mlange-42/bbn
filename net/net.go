package net

import "fmt"

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

func (n *Network) Sample(evidence map[string]int) error {
	ev := make([]int, len(n.nodes))
	for i := range ev {
		ev[i] = -1
	}
	for k, v := range evidence {
		idx, ok := n.byName[k]
		if !ok {
			return fmt.Errorf("node '%s' not found", k)
		}
		ev[idx] = v
	}

	//sample := make([]int, len(n.nodes))
	counts := make([][]int, len(n.nodes))

	for i := range counts {
		counts[i] = make([]int, len(n.nodes[i].States))
	}

	return nil
}
