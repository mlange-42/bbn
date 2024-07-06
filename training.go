package bbn

import (
	"fmt"
	"math"

	"github.com/mlange-42/bbn/ve"
)

// Trainer is a utility type to train a [Network].
type Trainer struct {
	network *Network
	data    [][][]float64
	counter [][]int
	indices [][]int
	sample  []int
	utility []float64
}

// NewTrainer creates a new [Trainer] for the given [Network].
func NewTrainer(net *Network) Trainer {
	nodes := net.Variables()

	data := make([][][]float64, len(nodes))
	counter := make([][]int, len(nodes))
	indices := make([][]int, len(nodes))

	nodeIndices := make(map[string]int, len(nodes))
	maxColumns := 0
	for i, node := range nodes {
		nodeIndices[node.Name] = i
		if len(node.Factor.Given) > maxColumns {
			maxColumns = len(node.Factor.Given)
		}
	}

	for i, node := range nodes {
		columns := len(node.Outcomes)
		rows := len(node.Factor.Table) / columns
		d := make([][]float64, rows)
		for j := 0; j < rows; j++ {
			d[j] = make([]float64, columns)
		}
		data[i] = d
		counter[i] = make([]int, rows)

		idx := make([]int, len(node.Factor.Given))
		for i, n := range node.Factor.Given {
			var ok bool
			idx[i], ok = nodeIndices[n]
			if !ok {
				panic(fmt.Sprintf("parent node %s for %s not found", n, node.Name))
			}
		}
		indices[i] = idx
	}

	return Trainer{
		network: net,
		data:    data,
		counter: counter,
		indices: indices,
		sample:  make([]int, 0, maxColumns),
		utility: make([]float64, 0, maxColumns),
	}
}

// AddSample adds a training sample.
// Order of values in the sample is the same as the order in which nodes were passed into the [Network] constructor.
func (t *Trainer) AddSample(sample []int, utility []float64) {
	nodes := t.network.Variables()

	for i, node := range nodes {
		if node.Type == ve.DecisionNode {
			continue
		}

		indices := t.indices[i]
		t.sample = t.sample[:0]
		for _, idx := range indices {
			t.sample = append(t.sample, sample[idx])
		}
		if utility != nil {
			t.utility = t.utility[:0]
			for _, idx := range indices {
				t.utility = append(t.utility, utility[idx])
			}
		}

		idx, ok := node.Factor.RowIndex(t.sample)
		if !ok {
			continue
		}

		if node.Type == ve.UtilityNode {
			u := utility[i]
			if math.IsNaN(u) {
				continue
			}
			t.data[i][idx][0] += u
			t.counter[i][idx]++
		} else {
			s := sample[i]
			if s < 0 {
				continue
			}
			t.data[i][idx][s]++
		}

		t.counter[i][idx]++
	}
}

// UpdateNetwork applies the training to the network, and returns a pointer to the original network.
func (t *Trainer) UpdateNetwork() (*Network, error) {
	nodes := t.network.Variables()

	for i, node := range nodes {
		data := t.data[i]

		cols := node.Factor.columns
		rows := len(node.Factor.Table) / cols
		for j := 0; j < rows; j++ {
			cnt := t.counter[i][j]
			if cnt == 0 {
				return nil, fmt.Errorf("no samples for node '%s', table row %d", node.Name, j)
			}
			for k := 0; k < cols; k++ {
				node.Factor.Table[j*cols+k] = float64(data[j][k]) / float64(cnt)
			}
		}
	}

	return t.network, nil
}
