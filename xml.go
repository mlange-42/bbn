package bbn

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type bifXmlWrapper struct {
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

// FromBIFXML creates a [Network] from XML. See also [FromFile].
func FromBIFXML(content []byte) (*Network, []*Node, error) {
	reader := bytes.NewReader(content)
	decoder := xml.NewDecoder(reader)

	net := bifXmlWrapper{}

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
