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
func FromBIFXML(content []byte) (*Network, error) {
	reader := bytes.NewReader(content)
	decoder := xml.NewDecoder(reader)

	bifNet := bifXmlWrapper{}

	err := decoder.Decode(&bifNet)
	if err != nil {
		return nil, err
	}

	defs := map[string]*definitionXml{}
	for i := range bifNet.Network.Definitions {
		def := &bifNet.Network.Definitions[i]
		defs[def.For] = def
	}

	variables := make([]Variable, len(bifNet.Network.Variables))
	factors := []Factor{}
	for i, variable := range bifNet.Network.Variables {
		def := defs[variable.Name]

		columns := len(variable.Outcomes)
		tableValues := strings.Fields(def.Table)
		if columns > 0 {
			if len(tableValues)%columns != 0 {
				return nil, fmt.Errorf("number of values in table for node '%s' does not match expected number", variable.Name)
			}
		}

		position, err := parsePosition(&variable)
		if err != nil {
			return nil, err
		}
		tp, ok := nodeTypes[variable.Type]
		if !ok {
			return nil, fmt.Errorf("unknown node type %s", variable.Type)
		}

		variables[i] = Variable{
			Name:     variable.Name,
			Type:     tp,
			Outcomes: variable.Outcomes,
			Position: position,
		}

		var table []float64
		if len(tableValues) > 0 {
			table = make([]float64, len(tableValues))
			for i := range table {
				v, err := strconv.ParseFloat(tableValues[i], 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing table value in node '%s' to float", variable.Name)
				}
				table[i] = v
			}
		}
		factors = append(factors, Factor{
			For:   variable.Name,
			Given: def.Given,
			Table: table,
		})
	}

	n := New(bifNet.Network.Name, variables, factors)
	return n, nil
}

// Search and parse position property in format `position = (x, y)`
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
