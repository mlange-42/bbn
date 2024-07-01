package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromYaml(t *testing.T) {
	yml := `name: Umbrella Decision Network
variables:

- variable: Weather
  position: [16, 0]
  outcomes: [Sunny, Rainy]

- variable: Forecast
  position: [1, 8]
  outcomes: [Sunny, Cloudy, Rainy]

- variable: Umbrella
  position: [16, 16]
  type: decision
  outcomes: [Take, Leave]

- variable: Utility
  position: [31, 8]
  type: utility
  outcomes: [Expected]

factors:

- for: Weather
  table: 
  - [70, 30]
  
- for: Forecast
  given: [Weather]
  table:
  - [70, 20, 10] # Sunny
  - [15, 25, 60] # Rainy

- for: Umbrella
  given: [Forecast]

- for: Utility
  given: [Weather, Umbrella]
  table: 
  - [ 20] # Sunny, Take
  - [100] # Sunny, Leave
  - [ 70] # Rainy, Take
  - [  0] # Rainy, Leave
`
	n, err := FromYAML([]byte(yml))
	assert.Nil(t, err)

	err = n.SolvePolicies(true)
	assert.Nil(t, err)
}
