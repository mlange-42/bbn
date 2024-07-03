package bbn

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mlange-42/bbn/internal/logic"
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
	Type     string      `yaml:",omitempty"`      // Type of the node [nature, decision, utility]
	Outcomes []string    `yaml:",flow"`           // Names of the node's possible states.
	Position [2]int      `yaml:",flow"`           // Coordinates for visualization, optional.
	Color    string      `yaml:",omitempty"`      // Node color, optional.
	Logic    string      `yaml:",omitempty"`      // Logic operations, alternative to a table
	Table    [][]float64 `yaml:",flow,omitempty"` // Table with the variable's factor
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

		table, err := toTable(&v)
		if err != nil {
			return nil, err
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

func toTable(v *variableYaml) ([]float64, error) {
	if len(v.Table) > 0 && v.Logic != "" {
		return nil, fmt.Errorf("node can only have one of 'table' or 'logic'")
	}

	if v.Logic != "" {
		l, ok := logic.Factors[strings.ToLower(v.Logic)]
		if !ok {
			return nil, fmt.Errorf("unknown logic operator %s; valid operators are e.g.: not, and, or, xor, if-then, if-not-then, if-then-not, if-not-then-not, not-and, etc", v.Logic)
		}
		if len(v.Given) != l.Operands() {
			return nil, fmt.Errorf("logic %s requires %d operands, but %d were given",
				v.Logic, l.Operands(), len(v.Given))
		}
		return l.Table(), nil
	}

	var table []float64
	if len(v.Table) > 0 {
		table = make([]float64, 0, len(v.Table)*len(v.Table[0]))
		for _, row := range v.Table {
			table = append(table, row...)
		}
	}
	return table, nil
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
