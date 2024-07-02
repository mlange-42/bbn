package bbn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromYaml(t *testing.T) {
	yml := `name: Umbrella Decision Network
variables:

- variable: Weather
  position: [16, 0]
  outcomes: [Sunny, Rainy]
  table: 
  - [70, 30]

- variable: Forecast
  position: [1, 8]
  given: [Weather]
  outcomes: [Sunny, Cloudy, Rainy]
  table: 
  - [70, 20, 10] # Sunny
  - [15, 25, 60] # Rainy

- variable: Umbrella
  position: [16, 16]
  given: [Forecast]
  type: decision
  outcomes: [Take, Leave]

- variable: Utility
  position: [31, 8]
  type: utility
  given: [Weather, Umbrella]
  outcomes: [Expected]
  table: 
  - [ 20] # Sunny, Take
  - [100] # Sunny, Leave
  - [ 70] # Rainy, Take
  - [  0] # Rainy, Leave
`
	n, err := FromYAML([]byte(yml))
	assert.Nil(t, err)

	policy, err := n.SolvePolicies(true)
	assert.Nil(t, err)

	fmt.Println(policy)
}

func TestToYAML(t *testing.T) {
	rain := Variable{
		Name:     "Rain",
		Outcomes: []string{"yes", "no"},
	}

	sprinkler := Variable{
		Name:     "Sprinkler",
		Outcomes: []string{"yes", "no"},
	}

	grassWet := Variable{
		Name:     "GrassWet",
		Outcomes: []string{"yes", "no"},
	}

	fRain := Factor{
		For:   "Rain",
		Table: []float64{0.2, 0.8},
	}

	fSprinkler := Factor{
		For:   "Sprinkler",
		Given: []string{"Rain"},
		Table: []float64{
			0.01, 0.99, // rain yes
			0.2, 0.8, // rain no
		},
	}

	fGrass := Factor{
		For:   "GrassWet",
		Given: []string{"Rain", "Sprinkler"},
		Table: []float64{
			0.99, 0.01, // rain yes, sprikler yes
			0.8, 0.2, // rain yes, sprikler no
			0.9, 0.1, // rain no, sprikler yes
			0.0, 1.0, // rain no, sprikler no
		},
	}

	net := New("Sprinkler", []Variable{rain, sprinkler, grassWet}, []Factor{fRain, fSprinkler, fGrass})

	yml, err := ToYAML(net)
	assert.Nil(t, err)

	expected := `name: Sprinkler
variables:
  - variable: Rain
    outcomes: ["yes", "no"]
    position: [0, 0]
    table: [[0.2, 0.8]]
  - variable: Sprinkler
    given: [Rain]
    outcomes: ["yes", "no"]
    position: [0, 0]
    table: [[0.01, 0.99], [0.2, 0.8]]
  - variable: GrassWet
    given: [Rain, Sprinkler]
    outcomes: ["yes", "no"]
    position: [0, 0]
    table: [[0.99, 0.01], [0.8, 0.2], [0.9, 0.1], [0, 1]]
`
	assert.Equal(t, expected, string(yml))

	net, err = FromYAML(yml)
	assert.Nil(t, err)

	evidence := map[string]string{
		"Rain":     "no",
		"GrassWet": "yes",
	}
	result, _, err := net.SolveQuery(evidence, []string{"Sprinkler"}, false)
	assert.Nil(t, err)

	assert.Equal(t, map[string][]float64{
		"Sprinkler": {1, 0},
	}, result)
}
