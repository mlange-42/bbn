package tui

import (
	"fmt"
	"math"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/ve"
)

func Solve(network *bbn.Network, evidence map[string]string, nodes []Node, ignorePolicies bool) (map[string][]float64, error) {
	queries := []string{}

	for _, n := range nodes {
		if n.Node().Type == ve.UtilityNode {
			continue
		}
		if _, ok := evidence[n.Node().Name]; ok {
			continue
		}
		queries = append(queries, n.Node().Name)

	}

	result := map[string][]float64{}
	err := solveEvidence(network, evidence, result)
	if err != nil {
		return nil, err
	}

	totalProb, err := solveQueries(network, evidence, queries, ignorePolicies, result)
	if err != nil {
		return nil, err
	}

	err = solveUtility(network, nodes, evidence, totalProb, ignorePolicies, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func solveEvidence(network *bbn.Network, evidence map[string]string, result map[string][]float64) error {
	for variable, value := range evidence {
		p, err := network.ToEvidence(variable, value)
		if err != nil {
			return err
		}
		result[variable] = p
	}
	return nil
}

func solveQueries(network *bbn.Network, evidence map[string]string, queries []string, ignorePolicies bool, result map[string][]float64) (float64, error) {
	_, f, err := network.SolveQuery(evidence, []string{}, ignorePolicies)
	if err != nil {
		return 0, err
	}
	totalProb := math.NaN()
	if len(f.Data()) == 1 {
		totalProb = f.Data()[0]
	}

	for _, q := range queries {
		r, _, err := network.SolveQuery(evidence, []string{q}, ignorePolicies)
		if err != nil {
			return 0, err
		}
		var ok bool
		result[q], ok = r[q]
		if !ok {
			panic(fmt.Sprintf("query variable %s not in result", q))
		}
	}
	return totalProb, nil
}

func solveUtility(network *bbn.Network, nodes []Node, evidence map[string]string, totalProb float64, ignorePolicies bool, result map[string][]float64) error {
	utilities := []string{}
	var totalUtilityNode *bbn.Variable

	for i, n := range nodes {
		if i == network.TotalUtilityIndex() {
			totalUtilityNode = n.Node()
			continue
		}
		if n.Node().Type == ve.UtilityNode {
			utilities = append(utilities, n.Node().Name)
		}
	}

	f, err := network.SolveUtility(evidence, []string{}, "", ignorePolicies)
	if err != nil {
		return err
	}

	totalUtility := math.NaN()
	if len(f.Data()) == 1 {
		totalUtility = f.Data()[0]
	}

	totalUtility /= totalProb

	for _, n := range utilities {
		f, err = network.SolveUtility(evidence, []string{}, n, ignorePolicies)
		if err != nil {
			return err
		}
		nodeUtility := math.NaN()
		if len(f.Data()) == 1 {
			nodeUtility = f.Data()[0]
		}

		nodeUtility /= totalProb

		result[n] = []float64{nodeUtility, totalUtility}
	}

	if totalUtilityNode != nil {
		util := make([]float64, len(totalUtilityNode.Factor.Given)+1)
		for i, g := range totalUtilityNode.Factor.Given {
			f := result[g]
			util[i] = f[0] * totalUtilityNode.Factor.Table[i]
		}
		util[len(util)-1] = totalUtility
		result[totalUtilityNode.Name] = util
	}

	return nil
}
