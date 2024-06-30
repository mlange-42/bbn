package ve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorIndex(t *testing.T) {
	f := Factor{
		Variables: []Variable{
			{Id: 0, outcomes: 3},
			{Id: 1, outcomes: 2},
		},
		Data: make([]float64, 6),
	}

	assert.Equal(t, 0, f.Index([]int{0, 0}))
	assert.Equal(t, 1, f.Index([]int{0, 1}))
	assert.Equal(t, 2, f.Index([]int{1, 0}))
	assert.Equal(t, 5, f.Index([]int{2, 1}))
}

func TestFactorRowIndex(t *testing.T) {
	f := Factor{
		Variables: []Variable{
			{Id: 0, outcomes: 3},
			{Id: 1, outcomes: 2},
			{Id: 2, outcomes: 4},
		},
		Data: make([]float64, 24),
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

func TestFactorOutcomes(t *testing.T) {
	f := Factor{
		Variables: []Variable{
			{Id: 0, outcomes: 3},
			{Id: 1, outcomes: 2},
		},
		Data: make([]float64, 6),
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
