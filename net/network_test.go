package net

import (
	"fmt"
	"math"
	"testing"

	"github.com/mlange-42/bbn/ve"
	"github.com/stretchr/testify/assert"
)

func TestNetworkToVE(t *testing.T) {
	vars := []Variable{
		{Name: "weather", Type: ve.ChanceNode, Outcomes: []string{"rainy", "sunny"}},
		{Name: "forecast", Type: ve.ChanceNode, Outcomes: []string{"sunny", "cloudy", "rainy"}},
		{Name: "umbrella", Type: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "utility", Type: ve.UtilityNode, Outcomes: []string{"utility"}},
	}

	factors := []Factor{
		{For: "weather", Table: []float64{
			// rain+, rain-
			0.3, 0.7,
		}},
		{For: "forecast", Given: []string{"weather"}, Table: []float64{
			// sunny, cloudy, rainy
			0.15, 0.25, 0.6, // rain+
			0.7, 0.2, 0.1, // rain-
		}},
		{For: "umbrella", Given: []string{"forecast"}},
		{For: "utility", Given: []string{"weather", "umbrella"}, Table: []float64{
			70,  // rain+, umbrella+
			0,   // rain+, umbrella-
			20,  // rain-, umbrella+
			100, // rain-, umbrella-
		}},
	}

	n := New(vars, factors)

	v, variables, err := n.ToVE()
	assert.Nil(t, err)

	result1 := v.SolveUtility(nil, nil, true)

	fmt.Println("Summarize")
	fmt.Println(result1)

	result := v.Variables.Rearrange(result1, []ve.Variable{variables["forecast"].VeVariable, variables["umbrella"].VeVariable})
	expected := []float64{
		12.95, 49, // sunny
		8.05, 14, // cloudy
		14, 7, // rainy
	}

	assert.Equal(t, len(expected), len(result.Data))

	for i := range expected {
		assert.Less(t, math.Abs(expected[i]-result.Data[i]), 0.0001)
	}
}

func TestNetworkSolveUmbrella(t *testing.T) {
	vars := []Variable{
		{Name: "weather", Type: ve.ChanceNode, Outcomes: []string{"rainy", "sunny"}},
		{Name: "forecast", Type: ve.ChanceNode, Outcomes: []string{"sunny", "cloudy", "rainy"}},
		{Name: "umbrella", Type: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "utility", Type: ve.UtilityNode, Outcomes: []string{"utility"}},
	}

	factors := []Factor{
		{For: "weather", Table: []float64{
			// rain+, rain-
			0.3, 0.7,
		}},
		{For: "forecast", Given: []string{"weather"}, Table: []float64{
			// sunny, cloudy, rainy
			0.15, 0.25, 0.6, // rain+
			0.7, 0.2, 0.1, // rain-
		}},
		{For: "umbrella", Given: []string{"forecast"}},
		{For: "utility", Given: []string{"weather", "umbrella"}, Table: []float64{
			70,  // rain+, umbrella+
			0,   // rain+, umbrella-
			20,  // rain-, umbrella+
			100, // rain-, umbrella-
		}},
	}

	n := New(vars, factors)
	err := n.SolvePolicies(true)
	assert.Nil(t, err)

	result, err := n.SolveQuery(map[string]string{}, []string{"umbrella"}, false, true)
	assert.Nil(t, err)

	fmt.Println("--> Query", n.Normalize(result))

	result, err = n.SolveQuery(map[string]string{}, []string{"weather", "umbrella"}, true, true)
	assert.Nil(t, err)

	fmt.Println("--> Utility", result)
	fmt.Println("--> Utility", n.Marginal(result, "weather"))
	fmt.Println("--> Utility", n.Marginal(result, "umbrella"))
}