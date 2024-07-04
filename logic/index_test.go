package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	idx := Index(0)
	err := idx.SetArgs(1)
	assert.Nil(t, err)

	table, err := idx.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		1, 0, // T T T
		1, 0, // T T F
		0, 1, // T F T
		0, 1, // T F F
		1, 0, // F T T
		1, 0, // F T F
		0, 1, // F F T
		0, 1, // F F F
	}, table)
}

func TestIndexNot(t *testing.T) {
	idx := IndexNot(0)
	err := idx.SetArgs(1)
	assert.Nil(t, err)

	table, err := idx.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1, // T T T
		0, 1, // T T F
		1, 0, // T F T
		1, 0, // T F F
		0, 1, // F T T
		0, 1, // F T F
		1, 0, // F F T
		1, 0, // F F F
	}, table)
}

func TestIndexExcl(t *testing.T) {
	idx := IndexExcl(0)
	err := idx.SetArgs(1)
	assert.Nil(t, err)

	table, err := idx.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 1, // T T T
		0, 1, // T T F
		0, 1, // T F T
		0, 1, // T F F
		0, 1, // F T T
		1, 0, // F T F
		0, 1, // F F T
		0, 1, // F F F
	}, table)
}
