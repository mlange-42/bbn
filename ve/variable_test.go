package ve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariables(t *testing.T) {
	v := NewVariables()

	v1 := v.Add(2)
	v2 := v.Add(3)

	f := v.CreateFactor([]*Variable{v1, v2}, []float64{0.1, 0.2, 0.7, 0.6, 0.2, 0.2})

	assert.Equal(t, 5, f.Index([]int{1, 2}))
	assert.Equal(t, 0.2, f.Get([]int{1, 2}))
}
