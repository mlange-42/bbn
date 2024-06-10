package bbn

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// FromBIFXML creates a [Network] from YAML. See also [FromFile].
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

// ToYAML serializes a [Network] to YAML.
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
