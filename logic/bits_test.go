package logic_test

import (
	"testing"

	"github.com/mlange-42/bbn/logic"
	"github.com/stretchr/testify/assert"
)

func TestBits(t *testing.T) {
	bits := logic.Bits()
	table, err := bits.Table(3)
	assert.Nil(t, err)

	assert.Equal(t, []float64{
		0, 0, 0, 0, 0, 0, 0, 1, // T T T
		0, 0, 0, 0, 0, 0, 1, 0, // T T F
		0, 0, 0, 0, 0, 1, 0, 0, // T F T
		0, 0, 0, 0, 1, 0, 0, 0, // T F F
		0, 0, 0, 1, 0, 0, 0, 0, // F T T
		0, 0, 1, 0, 0, 0, 0, 0, // F T F
		0, 1, 0, 0, 0, 0, 0, 0, // F F T
		1, 0, 0, 0, 0, 0, 0, 0, // F F F
	}, table)
}
