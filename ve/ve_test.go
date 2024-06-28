package ve

import (
	"fmt"
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
	ve := New(vars, []Factor{fRain, fSprinkler, fGrass}, evidence, query)
	result := ve.Eliminate()

	for _, q := range query {
		fmt.Println(vars.Marginal(result, q))
	}

	pRain := vars.Marginal(result, rain)
	pRain.Normalize()
	assert.Equal(t, []float64{1, 0}, pRain.data)
}

func TestEliminateDecision(t *testing.T) {
	vars := NewVariables()

	weather := vars.Add(ChanceNode, 2)
	forecast := vars.Add(ChanceNode, 3)
	umbrella := vars.Add(DecisionNode, 2)
	utility := vars.Add(UtilityNode, 1)
	_ = utility

	fWeather := vars.CreateFactor([]Variable{weather}, []float64{
		0.3, 0.7,
	})

	fForecast := vars.CreateFactor([]Variable{weather, forecast}, []float64{
		0.7, 0.2, 0.1, // sunny
		0.15, 0.25, 0.6, // rainy
	})

	fUtility := vars.CreateFactor([]Variable{weather, umbrella}, []float64{
		20,  // sunny, umbrella+
		100, // sunny, umbrella-
		70,  // rainy, umbrella+
		0,   // rainy, umbrella-
	})

	evidence := []Evidence{}
	query := []Variable{umbrella}
	ve := New(vars, []Factor{fWeather, fForecast, fUtility}, evidence, query)

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
		marg := vars.Marginal(result, q)
		if q.nodeType == ChanceNode {
			marg.Normalize()
		}
		fmt.Println(marg)
	}
}
