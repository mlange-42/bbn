package logic_test

import (
	"testing"

	"github.com/mlange-42/bbn/logic"
	"github.com/stretchr/testify/assert"
)

func TestFactorTable(t *testing.T) {
	result, err := logic.Not.Table(1)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		0, 1,
		1, 0,
	}, result)

	_, err = logic.Not.Table(2)
	assert.NotNil(t, err)

	result, err = logic.And.Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
	}, result)

	result, err = logic.Or.Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		1, 0,
		1, 0,
		0, 1,
	}, result)

	result, err = logic.Cond.Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		1, 0,
		1, 0,
	}, result)

	result, err = logic.BiCond.Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		1, 0,
	}, result)

	result, err = logic.IfThen.Table(1)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0.5, 0.5,
	}, result)
}
