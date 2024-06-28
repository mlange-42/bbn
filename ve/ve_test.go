package ve

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEliminate(t *testing.T) {
	vars := NewVariables()

	rain := vars.Add(ChanceNode, 2)
	sprinkler := vars.Add(ChanceNode, 2)
	grass := vars.Add(ChanceNode, 2)

	fRain := vars.CreateFactor([]Variable{rain}, []float64{
		0.2, 0.8,
	})

	fSprinkler := vars.CreateFactor([]Variable{rain, sprinkler}, []float64{
		0.01, 0.99, // rain+
		0.2, 0.8, // rain-
	})

	fGrass := vars.CreateFactor([]Variable{rain, sprinkler, grass}, []float64{
		0.99, 0.01, // rain+ sprinkler+
		0.8, 0.2, // rain+ sprinkler-
		0.9, 0.1, // rain- sprinkler+
		0.0, 1.0, // rain- sprinkler-
	})

	evidence := []Evidence{{Variable: sprinkler, Value: 1}, {Variable: grass, Value: 0}}
	query := []Variable{rain}
	ve := New(vars,
		[]Factor{fRain, fSprinkler, fGrass},
		nil,
		evidence, query)
	result := ve.Eliminate()

	for _, q := range query {
		fmt.Println(vars.Marginal(result, q))
	}

	pRain := vars.Marginal(result, rain)
	pRain.Normalize()
	assert.Equal(t, []float64{1, 0}, pRain.data)
}

func TestEliminateDecision(t *testing.T) {
	v := NewVariables()

	weather := v.Add(ChanceNode, 2)
	forecast := v.Add(ChanceNode, 3)
	umbrella := v.Add(DecisionNode, 2)
	utility := v.Add(UtilityNode, 1)
	_ = utility

	fWeather := v.CreateFactor([]Variable{weather}, []float64{
		// rain+, rain-
		0.3, 0.7,
	})

	fForecast := v.CreateFactor([]Variable{weather, forecast}, []float64{
		// sunny, cloudy, rainy
		0.15, 0.25, 0.6, // rain+
		0.7, 0.2, 0.1, // rain-
	})

	fUtility := v.CreateFactor([]Variable{weather, umbrella}, []float64{
		70,  // rain+, umbrella+
		0,   // rain+, umbrella-
		20,  // rain-, umbrella+
		100, // rain-, umbrella-
	})

	evidence := []Evidence{}
	query := []Variable{}
	ve := New(v,
		[]Factor{fWeather, fForecast, fUtility},
		[]Dependencies{{Decision: umbrella, Parents: []Variable{forecast}}},
		evidence, query)

	ve.eliminateEvidence()
	fmt.Println("Eliminate evidence")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.eliminateHidden()

	fmt.Println("Eliminate hidden")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	result := ve.summarize()

	fmt.Println("Summarize")
	fmt.Println(result)

	fmt.Println("Marginalize")
	for _, q := range query {
		marg := v.Marginal(result, q)
		if q.nodeType == ChanceNode {
			marg.Normalize()
		}
		fmt.Println(marg)
	}

	var expected []float64
	if result.variables[0].id == 1 {
		expected = []float64{
			12.95, 49, // sunny
			8.05, 14, // cloudy
			14, 7, // rainy
		}
	} else if result.variables[0].id == 2 {
		expected = []float64{
			// sunny, cloudy, rainy
			12.95, 8.05, 14, // umbrella+
			49, 14, 7, // umbrella-
		}
	} else {
		panic("unexpected variable order")
	}

	assert.Equal(t, len(expected), len(result.data))

	for i := range expected {
		assert.Less(t, math.Abs(expected[i]-result.data[i]), 0.0001)
	}
}
