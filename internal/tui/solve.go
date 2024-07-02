package tui

import (
	"fmt"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/ve"
)

func Solve(network *bbn.Network, evidence map[string]string, nodes []Node, ignorePolicies bool) (map[string][]float64, error) {
	queries := []string{}
	utilities := []string{}

	totalUtilityName := ""
	var totalUtilityNode *bbn.Variable

	for i, n := range nodes {
		if i == network.TotalUtilityIndex() {
			totalUtilityName = n.Node().Name
			totalUtilityNode = n.Node()
			continue
		}
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

	_, f, err := network.SolveQuery(evidence, []string{}, ignorePolicies)
	if err != nil {
		return nil, err
	}
	totalProb := f.Data[0]

	for _, q := range queries {
		r, _, err := network.SolveQuery(evidence, []string{q}, ignorePolicies)
		if err != nil {
			return nil, err
		}
		var ok bool
		result[q], ok = r[q]
		if !ok {
			panic(fmt.Sprintf("query variable %s not in result", q))
		}
	}

	f, err = network.SolveUtility(evidence, []string{}, "", ignorePolicies)
	if err != nil {
		return nil, err
	}
	totalUtility := f.Data[0]
	if totalProb != 0 {
		totalUtility /= totalProb
	}

	for _, n := range utilities {
		f, err = network.SolveUtility(evidence, []string{}, n, ignorePolicies)
		if err != nil {
			return nil, err
		}
		nodeUtility := f.Data[0]
		if totalProb != 0 {
			nodeUtility /= totalProb
		}

		result[n] = []float64{nodeUtility, totalUtility}
	}

	if totalUtilityName != "" {
		util := make([]float64, len(totalUtilityNode.Factor.Given)+1)
		for i, g := range totalUtilityNode.Factor.Given {
			f := result[g]
			util[i] = f[0] * totalUtilityNode.Factor.Table[i]
		}
		util[len(util)-1] = totalUtility
		result[totalUtilityName] = util
	}

	return result, nil
}
