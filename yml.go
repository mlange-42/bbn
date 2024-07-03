package bbn

import (
	"bytes"
	"fmt"

	"github.com/mlange-42/bbn/internal/ve"
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

var nodeTypeNames = map[ve.NodeType]string{
	ve.ChanceNode:   "",
	ve.DecisionNode: DecisionNodeType,
	ve.UtilityNode:  UtilityNodeType,
}

type variableYaml struct {
	Variable string      // Name of the node.
	Given    []string    `yaml:",flow,omitempty"`
	Type     string      `yaml:",omitempty"` // Type of the node [nature, decision, utility]
	Outcomes []string    `yaml:",flow"`      // Names of the node's possible states.
	Position [2]int      `yaml:",flow"`      // Coordinates for visualization, optional.
	Color    string      `yaml:",omitempty"` // Node color, optional.
	Table    [][]float64 `yaml:",flow,omitempty"`
}

type networkYaml struct {
	Name      string
	Variables []variableYaml
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
	factors := []Factor{}
	for i, v := range net.Variables {
		tp, ok := nodeTypes[v.Type]
		if !ok {
			return nil, fmt.Errorf("unknown node type %s", v.Type)
		}
		variables[i] = Variable{
			Name:     v.Variable,
			Type:     tp,
			Outcomes: v.Outcomes,
			Position: v.Position,
			Color:    v.Color,
		}

		var table []float64
		if len(v.Table) > 0 {
			table = make([]float64, 0, len(v.Table)*len(v.Table[0]))
			for _, row := range v.Table {
				table = append(table, row...)
			}
		}
		factors = append(factors, Factor{
			For:   v.Variable,
			Given: v.Given,
			Table: table,
		})
	}

	n := New(net.Name, variables, factors)

	return n, nil
}

func ToYAML(network *Network) ([]byte, error) {
	variables := make([]variableYaml, len(network.variables))
	for i, v := range network.variables {
		cols := len(v.Outcomes)
		table := make([][]float64, len(v.Factor.Table)/cols)

		for i := range table {
			table[i] = v.Factor.Table[i*cols : (i+1)*cols]
		}

		variables[i] = variableYaml{
			Variable: v.Name,
			Given:    v.Factor.Given,
			Type:     nodeTypeNames[v.Type],
			Outcomes: v.Outcomes,
			Position: v.Position,
			Table:    table,
		}
	}

	net := networkYaml{
		Name:      network.Name(),
		Variables: variables,
	}

	writer := bytes.Buffer{}
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(2)

	err := encoder.Encode(net)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
