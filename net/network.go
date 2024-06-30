package net

import (
	"fmt"
	"slices"

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

type variable struct {
	Variable   Variable
	VeVariable ve.Variable
}

type Network struct {
	variables     []Variable
	factors       []Factor
	policies      map[string]ve.Factor
	ve            *ve.VE
	variableNames map[string]variable
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

	var err error
	n.ve, n.variableNames, err = n.ToVE()
	if err != nil {
		return err
	}

	policies := n.ve.SolvePolicies(verbose)
	for name, v := range n.variableNames {
		if v.VeVariable.NodeType != ve.DecisionNode {
			continue
		}
		if p, ok := policies[v.VeVariable]; ok {
			n.policies[name] = *p[1]
		}
	}

	return nil
}

func (n *Network) SolveQuery(evidence map[string]string, query []string, utility bool, verbose bool) (*ve.Factor, error) {
	var err error
	n.ve, n.variableNames, err = n.ToVE()
	if err != nil {
		return nil, err
	}

	ev := []ve.Evidence{}
	for name, value := range evidence {
		vv, ok := n.variableNames[name]
		if !ok {
			return nil, fmt.Errorf("evidence variable %s not found", name)
		}
		idx := slices.Index(vv.Variable.Outcomes, value)
		if idx < 0 {
			return nil, fmt.Errorf("outcome %s for evidence variable %s not found", value, name)
		}
		ev = append(ev, ve.Evidence{Variable: vv.VeVariable, Value: idx})
	}

	q := make([]ve.Variable, len(query))
	for i, name := range query {
		vv, ok := n.variableNames[name]
		if !ok {
			return nil, fmt.Errorf("query variable %s not found", name)
		}
		q[i] = vv.VeVariable
	}

	if utility {
		return n.ve.SolveUtility(ev, q, verbose), nil
	} else {
		return n.ve.SolveQuery(ev, q, verbose), nil
	}
}

func (n *Network) ToVE() (*ve.VE, map[string]variable, error) {
	vars := ve.NewVariables()
	varNames := map[string]variable{}
	dependencies := map[ve.Variable][]ve.Variable{}

	for _, v := range n.variables {
		if v.Type == ve.DecisionNode {
			if _, ok := n.policies[v.Name]; ok {
				varNames[v.Name] = variable{
					Variable:   v,
					VeVariable: vars.Add(ve.ChanceNode, uint16(len(v.Outcomes))),
				}
				continue
			}
		}
		varNames[v.Name] = variable{
			Variable:   v,
			VeVariable: vars.Add(v.Type, uint16(len(v.Outcomes))),
		}
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
			variables[j] = vv.VeVariable
		}

		if forVar.VeVariable.NodeType == ve.DecisionNode {
			dependencies[forVar.VeVariable] = variables
			continue
		}
		if forVar.Variable.Type == ve.DecisionNode {
			continue
		}

		variables = append(variables, forVar.VeVariable)
		factors = append(factors, vars.CreateFactor(variables, f.Table))
	}

	fmt.Println("Solved policies", n.policies)
	for k, f := range n.policies {
		variable := varNames[k]
		variables := make([]ve.Variable, len(f.Variables))
		for i, v := range f.Variables {
			if v.Id == variable.VeVariable.Id {
				v.NodeType = ve.ChanceNode
			}
			variables[i] = v
		}
		factors = append(factors, vars.CreateFactor(variables, f.Data))
	}

	return ve.New(vars, factors, dependencies), varNames, nil
}

func (n *Network) Normalize(f *ve.Factor) ve.Factor {
	return n.ve.Variables.Normalize(f)
}

func (n *Network) Marginal(f *ve.Factor, v string) ve.Factor {
	vv, ok := n.variableNames[v]
	if !ok {
		panic(fmt.Sprintf("marginal: variable %s not found", v))
	}
	return n.ve.Variables.Marginal(f, vv.VeVariable)
}
