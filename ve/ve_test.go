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

func TestDecisionUmbrella(t *testing.T) {
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

	ve.eliminateDecisions()

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

func TestDecisionEvacuate(t *testing.T) {
	v := NewVariables()

	earthquake := v.Add(ChanceNode, 3)
	sensor := v.Add(ChanceNode, 3)
	maintenance := v.Add(ChanceNode, 2)
	evacuate := v.Add(DecisionNode, 2)

	materialDamage := v.Add(UtilityNode, 1)
	humanDamage := v.Add(UtilityNode, 1)
	evacCost := v.Add(UtilityNode, 1)

	_, _, _ = materialDamage, humanDamage, evacCost

	fEarthquake := v.CreateFactor([]Variable{earthquake}, []float64{
		0.01, 0.05, 0.94,
	})
	fSensor := v.CreateFactor([]Variable{maintenance, earthquake, sensor}, []float64{
		0.9, 0.1, 0.0, // strong, good
		0.05, 0.9, 0.05, // slight, good
		0.0, 0.1, 0.9, // none,   good
		0.6, 0.4, 0.0, // strong, poor
		0.2, 0.6, 0.2, // slight, poor
		0.0, 0.4, 0.6, // none,   poor
	})
	fMaintenance := v.CreateFactor([]Variable{maintenance}, []float64{
		0.5, 0.5,
	})
	fMaterialDamage := v.CreateFactor([]Variable{earthquake}, []float64{
		1000, 250, 0,
	})
	fHumanDamage := v.CreateFactor([]Variable{evacuate, earthquake}, []float64{
		-100,  // strong, e+
		-20,   // slight, e+
		0,     // none, e+
		-5000, // strong, e-
		-250,  // slight, e-
		0,     // none, e-
	})
	fEvacCost := v.CreateFactor([]Variable{evacuate}, []float64{
		-100, 0,
	})

	query := []Variable{}
	ve := New(v,
		[]Factor{fEarthquake, fSensor, fMaintenance, fMaterialDamage, fHumanDamage, fEvacCost},
		[]Dependencies{{Decision: evacuate, Parents: []Variable{sensor}}},
		[]Evidence{}, query)

	ve.eliminateEvidence()
	fmt.Println("Eliminate evidence")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.eliminateDecisions()

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
}

func TestSortDecisions(t *testing.T) {
	v := NewVariables()

	d3 := v.Add(DecisionNode, 2)
	d2 := v.Add(DecisionNode, 2)
	d1 := v.Add(DecisionNode, 2)

	deps := []Dependencies{
		{Decision: d2, Parents: []Variable{d1}},
		{Decision: d3, Parents: []Variable{d2}},
	}

	ve := New(v,
		[]Factor{},
		deps,
		[]Evidence{}, []Variable{})

	assert.Equal(t, []Variable{d1, d2, d3}, ve.decisions)
}
