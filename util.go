package bbn

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConflictingEvidenceError struct{}

func (m *ConflictingEvidenceError) Error() string {
	return "conflicting evidence / all samples rejected"
}

type networkYaml struct {
	Name      string
	Variables []*Node
}

func FromFile(path string) (*Network, []*Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	ext := filepath.Ext(path)

	switch strings.ToLower(ext) {
	case ".yml":
		return FromYAML(data)
	case ".xml", ".bifxml":
		return FromBIFXML(data)
	default:
		return nil, nil, fmt.Errorf("unsupported file format '%s'", ext)
	}
}

func FromYAML(content []byte) (*Network, []*Node, error) {
	reader := bytes.NewReader(content)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	net := networkYaml{}
	err := decoder.Decode(&net)
	if err != nil {
		return nil, nil, err
	}

	n, err := New(net.Name, net.Variables...)
	if err != nil {
		return nil, nil, err
	}

	return n, net.Variables, nil
}

func ToYAML(net *Network) ([]byte, error) {
	def := networkYaml{
		Name:      net.name,
		Variables: make([]*Node, len(net.nodes)),
	}

	for _, node := range net.nodes {
		def.Variables[node.ID] = &Node{
			Variable: node.Variable,
			Given:    node.GivenNames,
			Outcomes: node.Outcomes,
			Table:    node.Table,
			Position: node.Position,
		}
	}

	writer := bytes.Buffer{}
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(2)

	err := encoder.Encode(def)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

type bifWrapper struct {
	Network networkXml `xml:"NETWORK"`
}

type networkXml struct {
	Name        string          `xml:"NAME"`
	Variables   []variableXml   `xml:"VARIABLE"`
	Definitions []definitionXml `xml:"DEFINITION"`
}

type variableXml struct {
	Name       string   `xml:"NAME"`
	Type       string   `xml:"TYPE,attr"`
	Outcomes   []string `xml:"OUTCOME"`
	Properties []string `xml:"PROPERTY"`
}

type definitionXml struct {
	For   string   `xml:"FOR"`
	Given []string `xml:"GIVEN"`
	Table string   `xml:"TABLE"`
}

func FromBIFXML(content []byte) (*Network, []*Node, error) {
	reader := bytes.NewReader(content)
	decoder := xml.NewDecoder(reader)

	net := bifWrapper{}

	err := decoder.Decode(&net)
	if err != nil {
		return nil, nil, err
	}

	defs := map[string]*definitionXml{}
	for i := range net.Network.Definitions {
		def := &net.Network.Definitions[i]
		defs[def.For] = def
	}

	nodes := make([]*Node, len(net.Network.Variables))
	for i, variable := range net.Network.Variables {
		def := defs[variable.Name]

		columns := len(variable.Outcomes)
		rows := 1
		tableValues := strings.Fields(def.Table)
		if columns > 0 {
			if len(tableValues)%columns != 0 {
				return nil, nil, fmt.Errorf("number of values in table for node '%s' does not match expected number", variable.Name)
			}
			rows = len(tableValues) / columns
		}
		table := make([][]float64, rows)

		for i := range table {
			row := make([]float64, columns)
			for j := 0; j < columns; j++ {
				v, err := strconv.ParseFloat(tableValues[i*columns+j], 64)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing table value in node '%s' to float", variable.Name)
				}
				row[j] = v
			}
			table[i] = row
		}
		position, err := parsePosition(&variable)
		if err != nil {
			return nil, nil, err
		}

		node := Node{
			Variable: variable.Name,
			Given:    def.Given,
			Outcomes: variable.Outcomes,
			Table:    table,
			Position: position,
		}
		nodes[i] = &node
	}

	n, err := New(net.Network.Name, nodes...)
	if err != nil {
		return nil, nil, err
	}

	return n, nodes, nil
}

func parsePosition(variable *variableXml) ([2]int, error) {
	position := [2]int{}
	for _, prob := range variable.Properties {
		parts := strings.Split(prob, "=")
		if len(parts) != 2 || strings.TrimSpace(parts[0]) != "position" {
			continue
		}
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "()")
		parts = strings.Split(value, ",")
		if len(parts) != 2 {
			return position, fmt.Errorf("syntax error in property 'position' of node '%s'", variable.Name)
		}
		var err error
		position[0], err = strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return position, fmt.Errorf("error parsing '%s' to integer in property 'position' of node '%s'", parts[0], variable.Name)
		}
		position[1], err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return position, fmt.Errorf("error parsing '%s' to integer in property 'position' of node '%s'", parts[1], variable.Name)
		}
		position[0] /= 2
		position[1] /= 12
	}

	return position, nil
}

// toInternalNodes transforms nodes to their internal representation.
func toInternalNodes(nodes []*Node) ([]*node, error) {
	nodeMap := map[string]int{}
	for i, n := range nodes {
		nodeMap[n.Variable] = i
	}

	nodeList := make([]*node, len(nodes))
	for i, n := range nodes {
		nodeMap[n.Variable] = i

		parents := make([]int, len(n.Given))
		for j, p := range n.Given {
			par, ok := nodeMap[p]
			if !ok {
				return nil, fmt.Errorf("parent node '%s' not found", p)
			}
			parents[j] = par
		}

		var stride []int
		tableRows := 1
		if len(parents) > 0 {
			stride = make([]int, len(n.Given))
			stride[len(stride)-1] = 1
			for j := len(stride) - 2; j >= 0; j-- {
				parIdx := parents[j+1]
				stride[j] = stride[j+1] * len(nodes[parIdx].Outcomes)
			}
			tableRows = stride[0] * len(nodes[parents[0]].Outcomes)
		}

		if len(n.Table) != tableRows {
			return nil, fmt.Errorf("wrong number of table rows in node '%s'; got %d, expected %d", n.Variable, len(n.Table), tableRows)
		}

		tableCols := len(n.Outcomes)
		if tableCols < 2 {
			return nil, fmt.Errorf("node '%s' requires at least two outcomes, got %d", n.Variable, tableCols)
		}

		for j, probs := range n.Table {
			if len(probs) != tableCols {
				return nil, fmt.Errorf("wrong number of table columns in node '%s', row %d; got %d, expected %d", n.Variable, j, len(probs), tableCols)
			}
		}

		nd := node{
			Variable:   n.Variable,
			ID:         i,
			GivenNames: n.Given,
			Given:      parents,
			Stride:     stride,
			Outcomes:   n.Outcomes,
			Table:      n.Table,
			TableCum:   nil,
			Position:   n.Position,
		}

		nodeList[i] = &nd
	}

	return nodeList, nil
}

// sortTopological sorts nodes in topological order.
func sortTopological(nodes []*node) ([]*node, error) {
	visited := make([]bool, len(nodes))
	stack := []int{}

	for i := range nodes {
		var err error
		stack, err = sortTopologicalRecursive(nodes, i, i, visited, stack)
		if err != nil {
			return nil, err
		}
	}

	newIndex := make([]int, len(nodes))
	for i, idx := range stack {
		newIndex[idx] = i
	}

	result := make([]*node, len(nodes))
	for i, idx := range stack {
		n := nodes[idx]
		for i, par := range n.Given {
			n.Given[i] = newIndex[par]
		}
		result[i] = n
	}

	return result, nil
}

// sortTopologicalRecursive performs the recursion used in sortTopological.
func sortTopologicalRecursive(nodes []*node, index int, start int, visited []bool, stack []int) ([]int, error) {
	if visited[index] {
		return stack, nil
	}

	visited[index] = true
	n := nodes[index]

	for _, parent := range n.Given {
		if parent == start {
			return nil, fmt.Errorf("graph has cycles")
		}
		var err error
		stack, err = sortTopologicalRecursive(nodes, parent, start, visited, stack)
		if err != nil {
			return nil, err
		}
	}

	stack = append(stack, index)

	return stack, nil
}

func cumulate(values []float64) []float64 {
	c := make([]float64, len(values))
	c[0] = values[0]
	for k := 1; k < len(values); k++ {
		c[k] = c[k-1] + values[k]
	}
	return c
}

// sample from cumulative (relative) probabilities.
func sample(cum []float64, rng *rand.Rand) int {
	ln := len(cum)
	r := rng.Float64() * cum[ln-1]

	for i, v := range cum {
		if v >= r {
			return i
		}
	}
	return -1
}
