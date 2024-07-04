package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexIs(t *testing.T) {
	idx := IndexIs(0, 0)
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

func TestIndexIsNot(t *testing.T) {
	idx := IndexIsNot(0, 0)
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

func TestIndexEither(t *testing.T) {
	idx := IndexEither(nil, 0)
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
