package tui

import (
	"fmt"

	"github.com/mlange-42/bbn/net"
	"github.com/mlange-42/bbn/ve"
)

func Solve(network *net.Network, evidence map[string]string, nodes []Node) (map[string][]float64, error) {
	queries := []string{}
	utilities := []string{}

	for _, n := range nodes {
		if n.Node().Type == ve.UtilityNode {
			utilities = append(utilities, n.Node().Name)
			continue
		}
		if _, ok := evidence[n.Node().Name]; ok {
			continue
		}
		queries = append(queries, n.Node().Name)

	}

	result := map[string][]float64{}
	for variable, value := range evidence {
		p, err := network.ToEvidence(variable, value)
		if err != nil {
			return nil, err
		}
		result[variable] = p
	}

	_, f, err := network.SolveQuery(evidence, []string{}, false)
	if err != nil {
		return nil, err
	}
	totalProb := f.Data[0]

	for _, q := range queries {
		r, _, err := network.SolveQuery(evidence, []string{q}, false)
		if err != nil {
			return nil, err
		}
		var ok bool
		result[q], ok = r[q]
		if !ok {
			panic(fmt.Sprintf("query variable %s not in result", q))
		}
	}

	f, err = network.SolveUtility(evidence, []string{}, false)
	if err != nil {
		return nil, err
	}
	totalUtility := f.Data[0]
	if totalProb != 0 {
		totalUtility /= totalProb
	}

	for _, n := range utilities {
		result[n] = []float64{totalUtility}
	}

	return result, nil
}
