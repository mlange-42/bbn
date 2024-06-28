package ve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablesCreateFactor(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(2)
	v2 := v.Add(3)

	f := v.CreateFactor([]Variable{v1, v2}, []float64{0.1, 0.2, 0.7, 0.6, 0.2, 0.2})

	assert.Equal(t, 5, f.Index([]int{1, 2}))
	assert.Equal(t, 0.2, f.Get([]int{1, 2}))
}

func TestVariablesRestrict(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(2)
	v2 := v.Add(3)

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

	v1 := v.Add(2)
	v2 := v.Add(3)

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

	v1 := v.Add(3)
	v3 := v.Add(2)
	v2 := v.Add(2)

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

func TestVariablesProductScalar(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(3)
	v2 := v.Add(2)

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
