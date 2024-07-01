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

	n := New("umbrella", vars, factors)

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

	n := New("umbrella", vars, factors)
	policy, err := n.SolvePolicies(false)
	assert.Nil(t, err)

	assert.Equal(t,
		map[string]Factor{
			"umbrella": {
				For:   "umbrella",
				Given: []string{"forecast"},
				Table: []float64{0, 1, 0, 1, 1, 0},
			}},
		policy)

	query := []string{"weather"}
	evidence := map[string]string{"forecast": "rainy"}

	result, f, err := n.SolveQuery(evidence, query, false)
	assert.Nil(t, err)

	fmt.Println("--> Query", f)
	for q, v := range result {
		fmt.Println("--> Query", q, v)
	}

	utility, err := n.SolveUtility(evidence, query, false)
	assert.Nil(t, err)

	fmt.Println("--> Utility", utility)
	for _, v := range query {
		fmt.Println("--> Utility", n.Marginal(utility, v))
	}

	normUtil := n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility", normUtil)
}

func TestNetworkSolveOil(t *testing.T) {
	vars := []Variable{
		{Name: "oil", Type: ve.ChanceNode, Outcomes: []string{"dry", "wet", "soaking"}},
		{Name: "test", Type: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "test-result", Type: ve.ChanceNode, Outcomes: []string{"closed", "open", "diffuse"}},
		{Name: "drill", Type: ve.DecisionNode, Outcomes: []string{"yes", "no"}},

		{Name: "drill-utility", Type: ve.UtilityNode, Outcomes: []string{"utility"}},
		{Name: "test-utility", Type: ve.UtilityNode, Outcomes: []string{"utility"}},
	}

	factors := []Factor{
		{For: "oil", Table: []float64{
			// dry, wet, soaking
			0.5, 0.3, 0.2,
		}},
		{For: "test-result", Given: []string{"oil", "test"}, Table: []float64{
			// closed, open, diffuse
			0.1, 0.3, 0.6, // dry, test+
			0.333, 0.333, 0.333, // dry, test-
			0.3, 0.4, 0.3, // wet, test+
			0.333, 0.333, 0.333, // wet, test-
			0.5, 0.4, 0.1, // soaking, test+
			0.333, 0.333, 0.333, // soaking, test-
		}},

		{For: "test", Given: []string{}},
		{For: "drill", Given: []string{"test", "test-result"}},

		{For: "drill-utility", Given: []string{"oil", "drill"}, Table: []float64{
			-70, // dry, drill+
			0,   // dry, drill-
			50,  // wet, test+
			0,   // wet, drill-
			200, // soaking, test+
			0,   // soaking, drill-
		}},
		{For: "test-utility", Given: []string{"test"}, Table: []float64{
			-10, // test+
			0,   // test-
		}},
	}

	n := New("oil", vars, factors)
	policy, err := n.SolvePolicies(false)
	assert.Nil(t, err)

	assert.Equal(t,
		map[string]Factor{
			"drill": {
				For:   "drill",
				Given: []string{"test", "test-result"},
				Table: []float64{1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0},
			},
			"test": {
				For:   "test",
				Given: []string{},
				Table: []float64{1, 0},
			},
		}, policy)

	query := []string{"test-result"}
	evidence := map[string]string{}

	result, f, err := n.SolveQuery(evidence, query, false)
	assert.Nil(t, err)

	fmt.Println("--> Query", f)
	for q, v := range result {
		fmt.Println("--> Query", q, v)
	}

	utility, err := n.SolveUtility(evidence, query, false)
	assert.Nil(t, err)

	fmt.Println("--> Utility", utility)
	for _, v := range query {
		fmt.Println("--> Utility", n.Marginal(utility, v))
	}

	normUtil := n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility", normUtil)
}
