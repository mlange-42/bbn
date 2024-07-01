package net

import (
	"bytes"
	"fmt"

	"github.com/mlange-42/bbn/ve"
	"gopkg.in/yaml.v3"
)

const (
	ChanceNodeType   = "nature"
	DecisionNodeType = "decision"
	UtilityNodeType  = "utility"
)

var nodeTypes = map[string]ve.NodeType{
	"":               ve.ChanceNode,
	ChanceNodeType:   ve.ChanceNode,
	DecisionNodeType: ve.DecisionNode,
	UtilityNodeType:  ve.UtilityNode,
}

type variableYaml struct {
	Variable string   // Name of the node.
	Type     string   `yaml:",omitempty"` // Type of the node [nature, decision, utility]
	Outcomes []string `yaml:",flow"`      // Names of the node's possible states.
	Position [2]int   `yaml:",flow"`      // Coordinates for visualization, optional.
}

type factorYaml struct {
	For   string
	Given []string    `yaml:",omitempty"`
	Table [][]float64 `yaml:",omitempty"`
}

type networkYaml struct {
	Name      string
	Variables []variableYaml
	Factors   []factorYaml
}

// FromBIFXML creates a [Network] from YAML. See also [FromFile].
func FromYAML(content []byte) (*Network, error) {
	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	net := networkYaml{}
	err := decoder.Decode(&net)
	if err != nil {
		return nil, err
	}

	variables := make([]Variable, len(net.Variables))
	factors := make([]Factor, len(net.Factors))
	for i, v := range net.Variables {
		tp, ok := nodeTypes[v.Type]
		if !ok {
			return nil, fmt.Errorf("unknown node type %s", v.Type)
		}
		variables[i] = Variable{
			Name:     v.Variable,
			Type:     tp,
			Outcomes: v.Outcomes,
		}
	}
	for i, f := range net.Factors {
		var table []float64
		if len(f.Table) > 0 {
			table = make([]float64, 0, len(f.Table)*len(f.Table[0]))
			for _, row := range f.Table {
				table = append(table, row...)
			}
		}
		factors[i] = Factor{
			For:   f.For,
			Given: f.Given,
			Table: table,
		}
	}

	n := New(variables, factors)

	return n, nil
}
