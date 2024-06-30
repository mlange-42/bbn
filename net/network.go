package net

import (
	"fmt"

	"github.com/mlange-42/bbn/ve"
)

type Variable struct {
	Name     string
	Type     ve.NodeType
	Outcomes []string
}

type Factor struct {
	For   string
	Given []string
	Table []float64
}

type Network struct {
	variables []Variable
	factors   []Factor
	policies  map[string]ve.Factor
}

func New(variables []Variable, factors []Factor) *Network {
	return &Network{
		variables: variables,
		factors:   factors,
		policies:  map[string]ve.Factor{},
	}
}

func (n *Network) SolvePolicies(verbose bool) error {
	clear(n.policies)

	varElim, variables, err := n.ToVE()
	if err != nil {
		return err
	}

	policies := varElim.SolvePolicies(verbose)
	for name, v := range variables {
		if v.NodeType != ve.DecisionNode {
			continue
		}
		if p, ok := policies[v]; ok {
			n.policies[name] = *p[1]
		}
	}

	return nil
}

func (n *Network) SolveQuery(evidence map[string]int, query []string, utility bool, verbose bool) (*ve.Factor, *ve.Variables, error) {
	varElim, variables, err := n.ToVE()
	if err != nil {
		return nil, nil, err
	}

	ev := []ve.Evidence{}
	for name, value := range evidence {
		vv, ok := variables[name]
		if !ok {
			return nil, nil, fmt.Errorf("evidence variable %s not found", name)
		}
		ev = append(ev, ve.Evidence{Variable: vv, Value: value})
	}

	q := make([]ve.Variable, len(query))
	for i, name := range query {
		vv, ok := variables[name]
		if !ok {
			return nil, nil, fmt.Errorf("query variable %s not found", name)
		}
		q[i] = vv
	}

	if utility {
		return varElim.SolveUtility(ev, q, verbose), varElim.Variables, nil
	} else {
		return varElim.SolveQuery(ev, q, verbose), varElim.Variables, nil
	}
}

func (n *Network) ToVE() (*ve.VE, map[string]ve.Variable, error) {
	vars := ve.NewVariables()
	varNames := map[string]ve.Variable{}
	dependencies := map[ve.Variable][]ve.Variable{}

	for _, v := range n.variables {
		if v.Type == ve.DecisionNode {
			if _, ok := n.policies[v.Name]; ok {
				varNames[v.Name] = vars.Add(ve.ChanceNode, uint16(len(v.Outcomes)))
				continue
			}
		}
		varNames[v.Name] = vars.Add(v.Type, uint16(len(v.Outcomes)))
	}

	factors := []ve.Factor{}
	for _, f := range n.factors {
		forVar, ok := varNames[f.For]
		if !ok {
			return nil, nil, fmt.Errorf("variable %s for factor not found", f.For)
		}

		variables := make([]ve.Variable, len(f.Given))
		for j, v := range f.Given {
			vv, ok := varNames[v]
			if !ok {
				return nil, nil, fmt.Errorf("variable %s in factor for %s not found", v, f.For)
			}
			variables[j] = vv
		}

		if forVar.NodeType == ve.DecisionNode {
			dependencies[forVar] = variables
			continue
		}
		if n.variables[forVar.Id].Type == ve.DecisionNode {
			continue
		}

		variables = append(variables, forVar)
		factors = append(factors, vars.CreateFactor(variables, f.Table))
	}

	fmt.Println("Solved policies", n.policies)
	for k, f := range n.policies {
		variable := varNames[k]
		variables := make([]ve.Variable, len(f.Variables))
		for i, v := range f.Variables {
			if v.Id == variable.Id {
				v.NodeType = ve.ChanceNode
			}
			variables[i] = v
		}
		factors = append(factors, vars.CreateFactor(variables, f.Data))
	}

	return ve.New(vars, factors, dependencies), varNames, nil
}
