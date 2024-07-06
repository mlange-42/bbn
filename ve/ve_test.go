package ve

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEliminate(t *testing.T) {
	vars := NewVariables()

	rain := vars.AddVariable(0, ChanceNode, 2)
	sprinkler := vars.AddVariable(1, ChanceNode, 2)
	grass := vars.AddVariable(2, ChanceNode, 2)

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
		nil, nil)
	result := ve.SolveQuery(evidence, query)

	for _, q := range query {
		fmt.Println(vars.Marginal(result, q))
	}

	pRain := vars.Marginal(result, rain)
	pRain = vars.Normalize(&pRain)
	assert.Equal(t, []float64{1, 0}, pRain.Data())
}

func TestDecisionUmbrella(t *testing.T) {
	v := NewVariables()

	weather := v.AddVariable(0, ChanceNode, 2)
	forecast := v.AddVariable(1, ChanceNode, 3)
	umbrella := v.AddVariable(2, DecisionNode, 2)
	utility := v.AddVariable(3, UtilityNode, 1)
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
	ve := New(v,
		[]Factor{fWeather, fForecast, fUtility},
		map[Variable][]Variable{umbrella: {forecast}},
		nil)

	result1 := ve.SolveUtility(evidence, nil, nil)

	fmt.Println("Summarize")
	fmt.Println(result1)

	result := ve.variables.Rearrange(result1, []Variable{forecast, umbrella})
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

func TestDecisionUmbrella2(t *testing.T) {
	v := NewVariables()

	weather := v.AddVariable(0, ChanceNode, 2)
	forecast := v.AddVariable(1, ChanceNode, 3)
	umbrella := v.AddVariable(2, DecisionNode, 2)
	utility := v.AddVariable(3, UtilityNode, 1)
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

	ve := New(v,
		[]Factor{fWeather, fForecast, fUtility},
		map[Variable][]Variable{umbrella: {weather, forecast}},
		nil)

	result := ve.SolvePolicies(false)

	fmt.Println("Summarize")
	for k, v := range result {
		fmt.Println(k, v[0], v[1])
	}

	policy := result[umbrella][1]
	assert.Equal(t, []Variable{weather, umbrella}, policy.Variables())
}

func TestDecisionEvacuate(t *testing.T) {
	v := NewVariables()

	earthquake := v.AddVariable(0, ChanceNode, 3)
	sensor := v.AddVariable(1, ChanceNode, 3)
	maintenance := v.AddVariable(2, ChanceNode, 2)
	evacuate := v.AddVariable(3, DecisionNode, 2)

	materialDamage := v.AddVariable(4, UtilityNode, 1)
	humanDamage := v.AddVariable(5, UtilityNode, 1)
	evacCost := v.AddVariable(6, UtilityNode, 1)

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

	evidence := []Evidence{}
	ve := New(v,
		[]Factor{fEarthquake, fSensor, fMaintenance, fMaterialDamage, fHumanDamage, fEvacCost},
		map[Variable][]Variable{evacuate: {sensor}},
		nil)

	result := ve.SolveUtility(evidence, nil, nil)

	fmt.Println("Summarize")
	fmt.Println(result)

	expectedUtility := v.Marginal(result, evacuate)
	assert.Equal(t, []float64{-124.5, -85.0}, expectedUtility.Data())
}

func TestDecisionOil(t *testing.T) {
	v := NewVariables()

	oil := v.AddVariable(0, ChanceNode, 3)
	test := v.AddVariable(1, DecisionNode, 2)
	testResult := v.AddVariable(2, ChanceNode, 3)
	drill := v.AddVariable(3, DecisionNode, 2)

	utilityDrill := v.AddVariable(4, UtilityNode, 1)
	utilityTest := v.AddVariable(5, UtilityNode, 1)

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
		-10, // test+
		0,   // test-
	})

	//query := []Variable{}
	//evidence := []Evidence{}
	ve := New(v,
		[]Factor{fOil, fResult, fUtilityTest, fUtilityDrill},
		map[Variable][]Variable{drill: {test, testResult}},
		nil)

	policies := ve.SolvePolicies(false)

	testPolicy := policies[test][1]
	drillPolicy := v.Rearrange(policies[drill][1], []Variable{test, testResult, drill})

	fmt.Println("Test policy:", testPolicy)
	fmt.Println("Drill policy:", drillPolicy)

	expTest := []float64{
		1, 0,
	}
	assert.Equal(t, expTest, testPolicy.Data())

	expDrill := []float64{
		// drill + -
		1, 0, // test+, closed
		1, 0, // test+, open
		0, 1, // test+, diffuse
		1, 0, // test+, closed
		1, 0, // test+, open
		1, 0, // test+, diffuse
	}
	assert.Equal(t, expDrill, drillPolicy.Data())

}

func TestDecisionRobot(t *testing.T) {
	v := NewVariables()

	accidentProb := 0.1

	short := v.AddVariable(0, DecisionNode, 2)
	pads := v.AddVariable(1, DecisionNode, 2)

	accident := v.AddVariable(2, ChanceNode, 2)

	utility := v.AddVariable(3, UtilityNode, 1)

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

	evidence := []Evidence{}
	ve := New(v,
		[]Factor{fAccident, fUtility},
		map[Variable][]Variable{},
		nil)

	result1 := ve.SolveUtility(evidence, nil, nil)
	result := v.Rearrange(result1, []Variable{short, pads})

	fmt.Println("Summarize")
	fmt.Println(result)

	expected := []float64{8 - 6*accidentProb, 10 - 10*accidentProb, 4, 6}
	assert.Equal(t, expected, result.Data())
}

func TestSortDecisions(t *testing.T) {
	v := NewVariables()

	d3 := v.AddVariable(0, DecisionNode, 2)
	d2 := v.AddVariable(1, DecisionNode, 2)
	d1 := v.AddVariable(2, DecisionNode, 2)

	deps := map[Variable][]Variable{
		d2: {d1},
		d3: {d2},
	}

	ve := New(v,
		[]Factor{},
		deps, nil)

	assert.Equal(t, []Variable{d1, d2, d3}, ve.getDecisions())
}
