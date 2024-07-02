package tui

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/ve"
)

func TrainNetwork(net *bbn.Network, nodes []bbn.Variable, dataFile, noData string, delimiter rune) (*bbn.Network, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.ReuseRecord = true
	r.Comma = delimiter

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	indices, isUtility, outcomes, err := prepare(nodes, header, noData)
	if err != nil {
		return nil, err
	}

	train := bbn.NewTrainer(net)
	sample := make([]int, len(nodes))
	utility := make([]float64, len(nodes))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for i, idx := range indices {
			if isUtility[i] {
				var err error
				if record[idx] == noData {
					utility[i] = math.NaN()
				}
				utility[i], err = strconv.ParseFloat(record[idx], 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse utility value '%s' to integer in node '%s'", record[idx], nodes[i].Name)
				}

			} else {
				var ok bool
				sample[i], ok = outcomes[i][record[idx]]
				if !ok {
					return nil, fmt.Errorf("outcome '%s' not available in node '%s'", record[idx], nodes[i].Name)
				}
			}
		}
		train.AddSample(sample, utility)
	}

	return train.UpdateNetwork()
}

func prepare(nodes []bbn.Variable, header []string, noData string) (indices []int, isUtility []bool, outcomes []map[string]int, err error) {
	indices = make([]int, len(nodes))
	isUtility = make([]bool, len(nodes))
	outcomes = make([]map[string]int, len(nodes))
	for i, node := range nodes {
		idx := slices.Index(header, node.Name)
		if idx < 0 {
			err = fmt.Errorf("no column '%s' in training data file", node.Name)
			return
		}
		indices[i] = idx

		outcomes[i] = make(map[string]int, len(node.Outcomes))
		for j, o := range node.Outcomes {
			if o == noData {
				err = fmt.Errorf("no-data value '%s' appears as outcomes of node '%s'", noData, node.Name)
				return
			}
			outcomes[i][o] = j
		}
		outcomes[i][noData] = -1

		isUtility[i] = node.Type == ve.UtilityNode
	}
	return
}
