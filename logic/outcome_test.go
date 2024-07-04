package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutcomeIs(t *testing.T) {
	idx := OutcomeIs(0, 0)
	err := idx.SetArgs(1, 4)
	assert.Nil(t, err)

	table, err := idx.Table(1)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1,
		1, 0,
		0, 1,
		0, 1,
	}, table)
}

func TestOutcomeIsNot(t *testing.T) {
	idx := OutcomeIsNot(0, 0)
	err := idx.SetArgs(2, 4)
	assert.Nil(t, err)

	table, err := idx.Table(1)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		1, 0,
		1, 0,
		0, 1,
		1, 0,
	}, table)
}

func TestOutcomeEither(t *testing.T) {
	idx := OutcomeEither(nil, 0)
	err := idx.SetArgs(1, 2, 4)
	assert.Nil(t, err)

	table, err := idx.Table(1)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1,
		1, 0,
		1, 0,
		0, 1,
	}, table)
}

func TestOutcomeLess(t *testing.T) {
	idx := OutcomeLess(0, 0)
	err := idx.SetArgs(2, 4)
	assert.Nil(t, err)

	table, err := idx.Table(1)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		1, 0,
		1, 0,
		0, 1,
		0, 1,
	}, table)
}

func TestOutcomeGreater(t *testing.T) {
	idx := OutcomeGreater(0, 0)
	err := idx.SetArgs(2, 4)
	assert.Nil(t, err)

	table, err := idx.Table(1)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1,
		0, 1,
		1, 0,
		1, 0,
	}, table)
}
