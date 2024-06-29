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

	fUtility := v.CreateFactor([]Variable{weather, umbrella, utility}, []float64{
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
	fMaterialDamage := v.CreateFactor([]Variable{earthquake, materialDamage}, []float64{
		-1000, -250, 0,
	})
	fHumanDamage := v.CreateFactor([]Variable{evacuate, earthquake, humanDamage}, []float64{
		-100,  // strong, e+
		-20,   // slight, e+
		0,     // none, e+
		-5000, // strong, e-
		-250,  // slight, e-
		0,     // none, e-
	})
	fEvacCost := v.CreateFactor([]Variable{evacuate, evacCost}, []float64{
		-100, 0,
	})

	query := []Variable{evacuate}
	evidence := []Evidence{}
	ve := New(v,
		[]Factor{fEarthquake, fSensor, fMaintenance, fMaterialDamage, fHumanDamage, fEvacCost},
		[]Dependencies{{Decision: evacuate, Parents: []Variable{sensor}}},
		evidence, query)

	ve.eliminateEvidence()
	fmt.Println("Eliminate evidence")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.sumUtilities()
	fmt.Println("Sum utilities")
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

	expectedUtility := v.Marginal(result, evacuate)
	assert.Equal(t, []float64{-124.5, -85.0}, expectedUtility.data)
}

func TestDecisionOil(t *testing.T) {
	v := NewVariables()

	oil := v.Add(ChanceNode, 3)
	test := v.Add(DecisionNode, 2)
	testResult := v.Add(ChanceNode, 3)
	drill := v.Add(DecisionNode, 2)

	utilityDrill := v.Add(UtilityNode, 1)
	utilityTest := v.Add(UtilityNode, 1)

	fOil := v.CreateFactor([]Variable{oil}, []float64{
		0.5, 0.3, 0.2,
	})

	fResult := v.CreateFactor([]Variable{oil, test, testResult}, []float64{
		// closed, open, diffuse
		0.1, 0.3, 0.6, // dry, test+
		0.333, 0.333, 0.333, // dry, test-
		0.3, 0.4, 0.3, // wet, test+
		0.333, 0.333, 0.333, // wet, test-
		0.5, 0.4, 0.1, // soaking, test+
		0.333, 0.333, 0.333, // soaking, test-
	})

	fUtilityDrill := v.CreateFactor([]Variable{oil, drill, utilityDrill}, []float64{
		-70, // dry, drill+
		0,   // dry, drill-
		50,  // wet, test+
		0,   // wet, drill-
		200, // soaking, test+
		0,   // soaking, drill-
	})

	fUtilityTest := v.CreateFactor([]Variable{test, utilityTest}, []float64{
		-10, // drill+
		0,   // drill-
	})

	query := []Variable{}
	evidence := []Evidence{}
	ve := New(v,
		[]Factor{fOil, fResult, fUtilityTest, fUtilityDrill},
		[]Dependencies{{Decision: drill, Parents: []Variable{test, testResult}}},
		evidence, query)

	ve.eliminateEvidence()
	fmt.Println("Eliminate evidence")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.sumUtilities()
	fmt.Println("Sum utilities")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.eliminateDecisions()

	fmt.Println("Eliminate hidden")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	result1 := ve.summarize()
	result := v.Rearrange(result1, []Variable{test, testResult, drill})

	fmt.Println("Summarize")
	fmt.Println(result)

	fmt.Println("Marginalize")
	for _, q := range query {
		marg := v.Marginal(&result, q)
		if q.nodeType == ChanceNode {
			marg.Normalize()
		}
		fmt.Println(marg)
	}

	exp := []float64{
		// drill + -
		18.6, -2.4, // test+, closed
		8.0, -3.5, // test+, open
		-16.6, -4.1, // test+, diffuse
		6.66, 0.0, // test+, closed
		6.66, 0.0, // test+, open
		6.66, 0.0, // test+, diffuse
	}
	assert.Equal(t, variables{test, testResult, drill}, result.variables)
	assert.Equal(t, len(exp), len(result.data))

	for i := range exp {
		assert.Less(t, math.Abs(exp[i]-result.data[i]), 0.00001)
	}
	expPolicy := []float64{
		// drill + -
		1, 0, // test+, closed
		1, 0, // test+, open
		0, 1, // test+, diffuse
		1, 0, // test+, closed
		1, 0, // test+, open
		1, 0, // test+, diffuse
	}
	policy := v.Policy(&result, drill)
	assert.Equal(t, variables{test, testResult, drill}, policy.variables)
	assert.Equal(t, expPolicy, policy.data)
}

func TestDecisionRobot(t *testing.T) {
	v := NewVariables()

	accidentProb := 0.1

	short := v.Add(DecisionNode, 2)
	pads := v.Add(DecisionNode, 2)

	accident := v.Add(ChanceNode, 2)

	utility := v.Add(UtilityNode, 1)

	fAccident := v.CreateFactor([]Variable{short, accident}, []float64{
		accidentProb, 1 - accidentProb, // short
		0, 1, // long
	})

	fUtility := v.CreateFactor([]Variable{pads, short, accident, utility}, []float64{
		2,  // pads+ short accident+
		8,  // pads+ short accident-
		0,  // pads+ long accident+
		4,  // pads+ long accident-
		0,  // pads- short accident+
		10, // pads- short accident-
		0,  // pads- long accident+
		6,  // pads- long accident-
	})

	query := []Variable{}
	evidence := []Evidence{}
	ve := New(v,
		[]Factor{fAccident, fUtility},
		[]Dependencies{},
		evidence, query)

	ve.eliminateEvidence()
	fmt.Println("Eliminate evidence")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.sumUtilities()
	fmt.Println("Sum utilities")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	ve.eliminateDecisions()

	fmt.Println("Eliminate hidden")
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}

	result1 := ve.summarize()
	result := v.Rearrange(result1, []Variable{short, pads})

	fmt.Println("Summarize")
	fmt.Println(result)

	fmt.Println("Marginalize")
	for _, q := range query {
		marg := v.Marginal(&result, q)
		if q.nodeType == ChanceNode {
			marg.Normalize()
		}
		fmt.Println(marg)
	}

	expected := []float64{8 - 6*accidentProb, 10 - 10*accidentProb, 4, 6}
	assert.Equal(t, expected, result.data)
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
