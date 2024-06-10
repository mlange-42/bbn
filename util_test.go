package bbn_test

import (
	"os"
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/stretchr/testify/assert"
)

func TestToYAML(t *testing.T) {
	rain := bbn.Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := bbn.Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := bbn.Node{
		Variable: "GrassWet",
		Given:    []string{"Rain", "Sprinkler"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.99, 0.01}, // rain yes, sprikler yes
			{0.8, 0.2},   // rain yes, sprikler no
			{0.9, 0.1},   // rain no, sprikler yes
			{0.0, 1.0},   // rain no, sprikler no
		},
	}

	net, err := bbn.New("Sprinkler", &sprinkler, &grassWet, &rain)
	assert.Nil(t, err)

	yml, err := bbn.ToYAML(net)
	assert.Nil(t, err)

	expected := `name: Sprinkler
variables:
  - variable: Sprinkler
    given: [Rain]
    outcomes: ["yes", "no"]
    table: [[0.01, 0.99], [0.2, 0.8]]
    position: [0, 0]
  - variable: GrassWet
    given: [Rain, Sprinkler]
    outcomes: ["yes", "no"]
    table: [[0.99, 0.01], [0.8, 0.2], [0.9, 0.1], [0, 1]]
    position: [0, 0]
  - variable: Rain
    given: []
    outcomes: ["yes", "no"]
    table: [[0.2, 0.8]]
    position: [0, 0]
`
	assert.Equal(t, expected, string(yml))
}

func TestFromBIFXML(t *testing.T) {
	xmlData, err := os.ReadFile("_examples/dog-problem.xml")
	assert.Nil(t, err)

	net, nodes, err := bbn.FromBIFXML(xmlData)
	assert.Nil(t, err)

	_, _ = net, nodes
}
