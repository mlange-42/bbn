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
	assert.Equal(t, []Variable{v2}, f2.variables)
	assert.Equal(t, []float64{0.6, 0.2, 0.2}, f2.data)

	f2 = v.Restrict(&f, v2, 2)
	assert.Equal(t, []Variable{v1}, f2.variables)
	assert.Equal(t, []float64{0.7, 0.2}, f2.data)
}
