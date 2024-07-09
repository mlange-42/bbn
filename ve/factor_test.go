package ve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorIndex(t *testing.T) {
	f := Factor{
		variables: []Variable{
			{id: 0, outcomes: 3},
			{id: 1, outcomes: 2},
		},
		data: make([]float64, 6),
	}

	assert.Equal(t, 0, f.Index([]int{0, 0}))
	assert.Equal(t, 1, f.Index([]int{0, 1}))
	assert.Equal(t, 2, f.Index([]int{1, 0}))
	assert.Equal(t, 5, f.Index([]int{2, 1}))
}

func TestFactorRowIndex(t *testing.T) {
	f := Factor{
		variables: []Variable{
			{id: 0, outcomes: 3},
			{id: 1, outcomes: 2},
			{id: 2, outcomes: 4},
		},
		data: make([]float64, 24),
	}

	idx, ln := f.RowIndex([]int{0, 0})
	assert.Equal(t, 0, idx)
	assert.Equal(t, 4, ln)

	idx, ln = f.RowIndex([]int{0, 1})
	assert.Equal(t, 4, idx)
	assert.Equal(t, 4, ln)

	idx, ln = f.RowIndex([]int{1, 0})
	assert.Equal(t, 8, idx)
	assert.Equal(t, 4, ln)

	idx, ln = f.RowIndex([]int{1, 1})
	assert.Equal(t, 12, idx)
	assert.Equal(t, 4, ln)
}

func TestFactorIndexWithNoData(t *testing.T) {
	f := Factor{
		variables: []Variable{
			{id: 0, outcomes: 3},
			{id: 1, outcomes: 2},
		},
		data: make([]float64, 6),
	}

	idx, ok := f.IndexWithNoData([]int{0, 0})
	assert.True(t, ok)
	assert.Equal(t, 0, idx)

	idx, ok = f.IndexWithNoData([]int{2, 1})
	assert.True(t, ok)
	assert.Equal(t, 5, idx)

	_, ok = f.IndexWithNoData([]int{2, -1})
	assert.False(t, ok)
}

func TestFactorOutcomes(t *testing.T) {
	f := Factor{
		variables: []Variable{
			{id: 0, outcomes: 3},
			{id: 1, outcomes: 2},
		},
		data: make([]float64, 6),
	}

	result := []int{0, 0}

	f.Outcomes(0, result)
	assert.Equal(t, []int{0, 0}, result)

	f.Outcomes(1, result)
	assert.Equal(t, []int{0, 1}, result)

	f.Outcomes(2, result)
	assert.Equal(t, []int{1, 0}, result)

	f.Outcomes(5, result)
	assert.Equal(t, []int{2, 1}, result)
}

func BenchmarkFactorIndex(b *testing.B) {
	b.StopTimer()

	vars := NewVariables()

	v1 := vars.AddVariable(0, ChanceNode, 2)
	v2 := vars.AddVariable(1, ChanceNode, 2)
	v3 := vars.AddVariable(2, ChanceNode, 2)

	f1 := vars.CreateFactor([]Variable{v1, v2, v3}, make([]float64, 8))

	var v int
	indices := []int{1, 0, 1}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v = f1.Index(indices)
	}
	b.StopTimer()

	vv := v + 1
	_ = vv
}

func BenchmarkFactorIndexWithNoData(b *testing.B) {
	b.StopTimer()

	vars := NewVariables()

	v1 := vars.AddVariable(0, ChanceNode, 2)
	v2 := vars.AddVariable(1, ChanceNode, 2)
	v3 := vars.AddVariable(2, ChanceNode, 2)

	f1 := vars.CreateFactor([]Variable{v1, v2, v3}, make([]float64, 8))

	var v int
	indices := []int{1, 0, 1}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v, _ = f1.IndexWithNoData(indices)
	}
	b.StopTimer()

	vv := v + 1
	_ = vv
}

func BenchmarkFactorOutcome(b *testing.B) {
	b.StopTimer()

	vars := NewVariables()

	v1 := vars.AddVariable(0, ChanceNode, 2)
	v2 := vars.AddVariable(1, ChanceNode, 2)
	v3 := vars.AddVariable(2, ChanceNode, 2)

	f1 := vars.CreateFactor([]Variable{v1, v2, v3}, make([]float64, 8))

	indices := []int{0, 0, 0}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f1.Outcomes(7, indices)
	}
	b.StopTimer()

	indices[0] = 1
	_ = indices
}
