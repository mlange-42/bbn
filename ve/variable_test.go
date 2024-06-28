package ve

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablesCreateFactor(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 2)
	v2 := v.Add(ChanceNode, 3)

	f := v.CreateFactor([]Variable{v1, v2}, []float64{0.1, 0.2, 0.7, 0.6, 0.2, 0.2})

	assert.Equal(t, 5, f.Index([]int{1, 2}))
	assert.Equal(t, 0.2, f.Get([]int{1, 2}))
}

func TestVariablesRestrict(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 2)
	v2 := v.Add(ChanceNode, 3)

	f := v.CreateFactor([]Variable{v1, v2}, []float64{0.1, 0.2, 0.7, 0.6, 0.2, 0.2})

	f2 := v.Restrict(&f, v1, 1)
	assert.Equal(t, variables{v2}, f2.variables)
	assert.Equal(t, []float64{0.6, 0.2, 0.2}, f2.data)

	f2 = v.Restrict(&f, v2, 2)
	assert.Equal(t, variables{v1}, f2.variables)
	assert.Equal(t, []float64{0.7, 0.2}, f2.data)
}

func TestVariablesSumOut(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 2)
	v2 := v.Add(ChanceNode, 3)

	f := v.CreateFactor([]Variable{v1, v2}, []float64{
		1, 2, 7,
		6, 2, 2,
	})

	f2 := v.SumOut(&f, v1)
	assert.Equal(t, variables{v2}, f2.variables)
	assert.Equal(t, []float64{7, 4, 9}, f2.data)

	f2 = v.SumOut(&f, v2)
	assert.Equal(t, variables{v1}, f2.variables)
	assert.Equal(t, []float64{10, 10}, f2.data)
}

func TestVariablesProduct(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 3)
	v3 := v.Add(ChanceNode, 2)
	v2 := v.Add(ChanceNode, 2)

	f1 := v.CreateFactor([]Variable{v1, v2}, []float64{
		0.1, 0.9,
		0.5, 0.5,
		0.8, 0.2,
	})

	f2 := v.CreateFactor([]Variable{v2, v3}, []float64{
		0.1, 0.9,
		0.8, 0.2,
	})

	f3 := v.Product(&f1, &f2)

	assert.Equal(t, variables{v1, v2, v3}, f3.variables)

	assert.Equal(t, f3.Get([]int{0, 0, 0}), f1.Get([]int{0, 0})*f2.Get([]int{0, 0}))

	assert.Equal(t, f3.Get([]int{0, 0, 1}), f1.Get([]int{0, 0})*f2.Get([]int{0, 1}))
	assert.Equal(t, f3.Get([]int{0, 1, 0}), f1.Get([]int{0, 1})*f2.Get([]int{1, 0}))
	assert.Equal(t, f3.Get([]int{1, 0, 0}), f1.Get([]int{1, 0})*f2.Get([]int{0, 0}))

	assert.Equal(t, f3.Get([]int{2, 1, 1}), f1.Get([]int{2, 1})*f2.Get([]int{1, 1}))
}

func TestVariablesProductMulti1(t *testing.T) {
	v := NewVariables()

	a := v.Add(ChanceNode, 2)
	b := v.Add(ChanceNode, 2)
	c := v.Add(ChanceNode, 2)

	fA := v.CreateFactor([]Variable{a}, []float64{
		1, 2,
	})

	fB := v.CreateFactor([]Variable{b}, []float64{
		2, 3,
	})

	fC := v.CreateFactor([]Variable{c}, []float64{
		4, 5,
	})

	prod := v.Product(&fA, &fB, &fC)

	assert.Equal(t, []float64{
		8, 10, // + +
		12, 15, // + -
		16, 20, // - +
		24, 30, // - -
	}, prod.data)
}

func TestVariablesProductMulti2(t *testing.T) {
	v := NewVariables()

	a := v.Add(ChanceNode, 2)
	b := v.Add(ChanceNode, 2)

	fA := v.CreateFactor([]Variable{a}, []float64{
		1, 2,
	})

	fB := v.CreateFactor([]Variable{b}, []float64{
		2, 3,
	})

	fC := v.CreateFactor([]Variable{b}, []float64{
		4, 5,
	})

	prod := v.Product(&fA, &fB, &fC)

	assert.Equal(t, []float64{
		8, 15,
		16, 30,
	}, prod.data)
}

func TestVariablesProductMulti3(t *testing.T) {
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

	prod := v.Product(&fWeather, &fForecast, &fUtility)

	expected := []float64{
		3.15, 0, // rain+, sunny
		5.25, 0, // rain+, cloudy
		12.6, 0, // rain+, rainy
		9.8, 49, // rain-, sunny
		2.8, 14, // rain-, cloudy
		1.4, 7, // rain-, rainy
	}
	assert.Equal(t, len(expected), len(prod.data))

	for i := range expected {
		assert.Less(t, math.Abs(expected[i]-prod.data[i]), 0.0001)
	}
}

func TestVariablesProductScalar(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 3)
	v2 := v.Add(ChanceNode, 2)

	f1 := v.CreateFactor([]Variable{v1, v2}, []float64{
		1, 9,
		5, 5,
		8, 2,
	})

	f2 := v.CreateFactor([]Variable{}, []float64{
		2,
	})

	f3 := v.Product(&f1, &f2)

	assert.Equal(t, variables{v1, v2}, f3.variables)

	assert.Equal(t, []float64{2, 18, 10, 10, 16, 4}, f3.data)
}

func TestVariablesPolicy(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(ChanceNode, 3)
	v2 := v.Add(ChanceNode, 2)

	f1 := v.CreateFactor([]Variable{v1, v2}, []float64{
		0.4, 0.6,
		0.9, 0.1,
		0.2, 0.8,
	})

	p := v.Policy(&f1, v2)
	assert.Equal(t, variables{v1, v2}, p.variables)
	assert.Equal(t, []float64{
		0, 1,
		1, 0,
		0, 1,
	}, p.data)

	p = v.Policy(&f1, v1)
	assert.Equal(t, variables{v2, v1}, p.variables)
	assert.Equal(t, []float64{
		0, 1, 0,
		0, 0, 1,
	}, p.data)
}
