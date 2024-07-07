package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorTable(t *testing.T) {
	result, err := Not().Table(1)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		0, 1,
		1, 0,
	}, result)

	_, err = Not().Table(2)
	assert.NotNil(t, err)

	result, err = And().Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		0, 1,
	}, result)

	result, err = Or().Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		1, 0,
		1, 0,
		0, 1,
	}, result)

	result, err = Cond().Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		1, 0,
		1, 0,
	}, result)

	result, err = BiCond().Table(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0, 1,
		0, 1,
		1, 0,
	}, result)

	result, err = IfThen().Table(1)
	assert.Nil(t, err)
	assert.Equal(t, []float64{
		1, 0,
		0.5, 0.5,
	}, result)

	var f boolFactor = []bool{
		true, // T T T
		true, // T T F
		true, // T F T
		true, // T F F
		true, // F T T
		true, // F T F
		true, // F F T
		true, // F F F
	}

	_, err = f.Table(3)
	assert.Nil(t, err)
}
