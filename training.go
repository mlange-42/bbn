package bbn

import "fmt"

type Trainer struct {
	network *Network
	data    [][][]int
	counter [][]int
	indices []int
	sample  []int
}

func NewTrainer(net *Network) Trainer {
	data := make([][][]int, len(net.nodes))
	counter := make([][]int, len(net.nodes))
	indices := make([]int, len(net.nodes))

	for i, node := range net.nodes {
		d := make([][]int, len(node.Table))
		for j, row := range node.Table {
			d[j] = make([]int, len(row))
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
	}
}

func (t *Trainer) AddSamples(sample []int) {
	for i, s := range sample {
		t.sample[t.indices[i]] = s
	}

	for i, node := range t.network.nodes {
		idx := node.Index(t.sample)
		s := t.sample[i]
		t.data[i][idx][s]++
		t.counter[i][idx]++
	}
}

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
