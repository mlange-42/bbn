package net

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
