package bbn

import (
	"fmt"
	"math"
	"testing"

	"github.com/mlange-42/bbn/ve"
	"github.com/stretchr/testify/assert"
)

func TestFactorRow(t *testing.T) {
	a := Variable{Name: "a", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no", "maybe"}}
	b := Variable{Name: "b", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no"}}
	c := Variable{Name: "c", NodeType: ve.DecisionNode, Outcomes: []string{"yes", "no", "maybe"}}

	factor := Factor{
		For:   "c",
		Given: []string{"a", "b"},
		Table: []float64{
			0, 1, 2,
			3, 4, 5,
			6, 7, 8,
			9, 10, 11,
			12, 13, 14,
			15, 16, 17,
		},
	}

	net, err := New("test", "", []Variable{a, b, c}, []Factor{factor})
	assert.Nil(t, err)

	variable := net.Variables()[2]

	s, ok := variable.Factor.rowIndex([]int{0, 0})
	assert.True(t, ok)
	assert.Equal(t, 0, s)

	s, ok = variable.Factor.rowIndex([]int{2, 0})
	assert.True(t, ok)
	assert.Equal(t, 4, s)

	s, ok = variable.Factor.rowIndex([]int{-1, 0})
	assert.False(t, ok)
	_ = s
}

func TestNetworkToVE(t *testing.T) {
	vars := []Variable{
		{Name: "weather", NodeType: ve.ChanceNode, Outcomes: []string{"rainy", "sunny"}},
		{Name: "forecast", NodeType: ve.ChanceNode, Outcomes: []string{"sunny", "cloudy", "rainy"}},
		{Name: "umbrella", NodeType: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "utility", NodeType: ve.UtilityNode, Outcomes: []string{"utility"}},
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

	n, err := New("umbrella", "", vars, factors)
	assert.Nil(t, err)

	v, variables, err := n.toVE(nil)
	assert.Nil(t, err)

	result1 := v.SolveUtility(nil, nil, nil)

	fmt.Println("Summarize")
	fmt.Println(result1)

	result := v.Variables().Rearrange(result1, []ve.Variable{variables["forecast"].VeVariable, variables["umbrella"].VeVariable})
	expected := []float64{
		12.95, 49, // sunny
		8.05, 14, // cloudy
		14, 7, // rainy
	}

	assert.Equal(t, len(expected), len(result.Data()))

	for i := range expected {
		assert.Less(t, math.Abs(expected[i]-result.Data()[i]), 0.0001)
	}
}

func TestNetworkSolveUmbrella(t *testing.T) {
	vars := []Variable{
		{Name: "weather", NodeType: ve.ChanceNode, Outcomes: []string{"rainy", "sunny"}},
		{Name: "forecast", NodeType: ve.ChanceNode, Outcomes: []string{"sunny", "cloudy", "rainy"}},
		{Name: "umbrella", NodeType: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "utility", NodeType: ve.UtilityNode, Outcomes: []string{"utility"}},
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

	n, err := New("umbrella", "", vars, factors)
	assert.Nil(t, err)
	policy, err := n.SolvePolicies(true)
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

	utility, err := n.SolveUtility(evidence, query, "", false)
	assert.Nil(t, err)

	fmt.Println("--> Utility", utility)
	for _, v := range query {
		fmt.Println("--> Utility", n.Marginal(utility, v))
	}

	normUtil := n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility", normUtil)

	expected := []float64{70, 20}
	assert.Equal(t, len(expected), len(normUtil.Data()))

	for i := range expected {
		assert.Less(t, math.Abs(expected[i]-normUtil.Data()[i]), 0.0001)
	}
}

func TestNetworkSolveOil(t *testing.T) {
	vars := []Variable{
		{Name: "oil", NodeType: ve.ChanceNode, Outcomes: []string{"dry", "wet", "soaking"}},
		{Name: "test", NodeType: ve.DecisionNode, Outcomes: []string{"yes", "no"}},
		{Name: "test-result", NodeType: ve.ChanceNode, Outcomes: []string{"closed", "open", "diffuse"}},
		{Name: "drill", NodeType: ve.DecisionNode, Outcomes: []string{"yes", "no"}},

		{Name: "drill-utility", NodeType: ve.UtilityNode, Outcomes: []string{"utility"}},
		{Name: "test-utility", NodeType: ve.UtilityNode, Outcomes: []string{"utility"}},
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

	n, err := New("oil", "", vars, factors)
	assert.Nil(t, err)
	policy, err := n.SolvePolicies(true)
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

	utility, err := n.SolveUtility(evidence, query, "", false)
	assert.Nil(t, err)

	fmt.Println("--> Utility", utility)
	for _, v := range query {
		fmt.Println("--> Utility", n.Marginal(utility, v))
	}

	normUtil := n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility", normUtil)

	utility, err = n.SolveUtility(evidence, query, "test-utility", false)
	assert.Nil(t, err)
	normUtil = n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility test-utility", normUtil)

	utility, err = n.SolveUtility(evidence, query, "drill-utility", false)
	assert.Nil(t, err)
	normUtil = n.NormalizeUtility(utility, f)
	fmt.Println("--> NormalizeUtility drill-utility", normUtil)
}

func TestNetworkRearrange(t *testing.T) {
	variables := []Variable{
		{Name: "a", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no"}},
		{Name: "b", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no"}},
		{Name: "c", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no"}},
		{Name: "d", NodeType: ve.ChanceNode, Outcomes: []string{"yes", "no"}},
	}

	fac := Factor{
		For:   "d",
		Given: []string{"a", "b", "c"},
	}

	net, err := New("test", "", variables, []Factor{fac})
	assert.Nil(t, err)

	_, f, err := net.SolveQuery(map[string]string{}, []string{"a", "b", "c", "d"}, true)
	assert.Nil(t, err)

	vars := []ve.Variable(f.Variables())
	a, b, c, d := vars[0], vars[1], vars[2], vars[3]

	result := net.rearrangeVariables(f, []string{"a", "b", "c", "d"})
	assert.Equal(t, vars, result)

	result = net.rearrangeVariables(f, []string{"d", "c", "b", "a"})
	assert.Equal(t, []ve.Variable{d, c, b, a}, result)

	result = net.rearrangeVariables(f, []string{"d", "a"})
	assert.Equal(t, []ve.Variable{d, b, c, a}, result)
}
