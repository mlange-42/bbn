package logic_test

import (
	"testing"

	"github.com/mlange-42/bbn/internal/logic"
	"github.com/stretchr/testify/assert"
)

func TestFactorOperands(t *testing.T) {
	assert.Equal(t, 1, logic.Not.Operands())
	assert.Equal(t, 2, logic.And.Operands())
}

func TestFactorTable(t *testing.T) {
	assert.Equal(t, []float64{
		0, 1,
		1, 0,
	}, logic.Not.Table())

	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
	}, logic.And.Table())

	assert.Equal(t, []float64{
		1, 0,
		1, 0,
		1, 0,
		0, 1,
	}, logic.Or.Table())

	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		1, 0,
		1, 0,
	}, logic.Cond.Table())

	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		1, 0,
	}, logic.BiCond.Table())
}
