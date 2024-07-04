package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountIs(t *testing.T) {
	exact := CountIs(0)
	err := exact.SetArgs(2)
	assert.Nil(t, err)

	table, err := exact.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1, // T T T
		1, 0, // T T F
		1, 0, // T F T
		0, 1, // T F F
		1, 0, // F T T
		0, 1, // F T F
		0, 1, // F F T
		0, 1, // F F F
	}, table)
}

func TestCountLess(t *testing.T) {
	exact := CountLess(0)
	err := exact.SetArgs(3)
	assert.Nil(t, err)

	table, err := exact.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1, // T T T
		1, 0, // T T F
		1, 0, // T F T
		1, 0, // T F F
		1, 0, // F T T
		1, 0, // F T F
		1, 0, // F F T
		1, 0, // F F F
	}, table)
}

func TestCountGreater(t *testing.T) {
	exact := CountGreater(0)
	err := exact.SetArgs(2)
	assert.Nil(t, err)

	table, err := exact.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		1, 0, // T T T
		0, 1, // T T F
		0, 1, // T F T
		0, 1, // T F F
		0, 1, // F T T
		0, 1, // F T F
		0, 1, // F F T
		0, 1, // F F F
	}, table)
}

func TestCountTrue(t *testing.T) {
	exact := CountTrue()

	table, err := exact.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 0, 0, 1, // T T T
		0, 0, 1, 0, // T T F
		0, 0, 1, 0, // T F T
		0, 1, 0, 0, // T F F
		0, 0, 1, 0, // F T T
		0, 1, 0, 0, // F T F
		0, 1, 0, 0, // F F T
		1, 0, 0, 0, // F F F
	}, table)
}

func TestCountFalse(t *testing.T) {
	exact := CountFalse()

	table, err := exact.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		1, 0, 0, 0, // T T T
		0, 1, 0, 0, // T T F
		0, 1, 0, 0, // T F T
		0, 0, 1, 0, // T F F
		0, 1, 0, 0, // F T T
		0, 0, 1, 0, // F T F
		0, 0, 1, 0, // F F T
		0, 0, 0, 1, // F F F
	}, table)
}
