package bbn

import (
	"fmt"
	"math"
)

// Trainer is a utility type to train a [Network].
type Trainer struct {
	network *Network
	data    [][][]float64
	counter [][]int
	indices []int
	sample  []int
	utility []float64
}

// NewTrainer creates a new [Trainer] for the given [Network].
func NewTrainer(net *Network) Trainer {
	data := make([][][]float64, len(net.nodes))
	counter := make([][]int, len(net.nodes))
	indices := make([]int, len(net.nodes))

	for i, node := range net.nodes {
		d := make([][]float64, len(node.Table))
		for j, row := range node.Table {
			d[j] = make([]float64, len(row))
		}
		data[i] = d
		counter[i] = make([]int, len(node.Table))
		indices[node.ID] = i
	}

	return Trainer{
		network: net,
		data:    data,
		counter: counter,
		indices: indices,
		sample:  make([]int, len(net.nodes)),
		utility: make([]float64, len(net.nodes)),
	}
}

// AddSample adds a training sample.
// Order of values in the sample is the same as the order in which nodes were passed into the [Network] constructor.
func (t *Trainer) AddSample(sample []int, utility []float64) {
	for i, s := range sample {
		t.sample[t.indices[i]] = s
	}
	for i, u := range utility {
		t.utility[t.indices[i]] = u
	}

	for i, node := range t.network.nodes {
		idx, ok := node.IndexWithNoData(t.sample)
		if !ok {
			continue
		}

		if node.Type == UtilityNode {
			u := t.utility[i]
			if math.IsNaN(u) {
				continue
			}
			t.data[i][idx][0] += u
		} else {
			s := t.sample[i]
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
	for i, node := range t.network.nodes {
		data := t.data[i]
		for j, row := range node.Table {
			cnt := t.counter[i][j]
			if cnt == 0 {
				return nil, fmt.Errorf("no samples for node '%s', table row %d", node.Variable, j)
			}
			for k := range row {
				row[k] = float64(data[j][k]) / float64(cnt)
			}
		}
	}

	t.network.cumulateTables()
	return t.network, nil
}
