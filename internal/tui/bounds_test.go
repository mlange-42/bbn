package tui_test

import (
	"testing"

	"github.com/mlange-42/bbn/internal/tui"
	"github.com/stretchr/testify/assert"
)

func TestNewBounds(t *testing.T) {
	b := tui.NewBounds(1, 2, 3, 4)

	assert.Equal(t, tui.Bounds{X: 1, Y: 2, W: 3, H: 4}, b)
}

func TestBoundsContains(t *testing.T) {
	b := tui.NewBounds(10, 20, 20, 30)

	assert.True(t, b.Contains(11, 22))
	assert.False(t, b.Contains(9, 22))
	assert.False(t, b.Contains(35, 22))
	assert.False(t, b.Contains(11, 15))
	assert.False(t, b.Contains(11, 51))
}

func TestBoundsExtend(t *testing.T) {
	a := tui.NewBounds(3, 4, 5, 6)
	b := tui.NewBounds(1, 2, 12, 13)
	a.Extend(b)
	assert.Equal(t, b, a)

	a = tui.NewBounds(3, 4, 5, 6)
	b = tui.NewBounds(5, 6, 12, 13)
	a.Extend(b)
	assert.Equal(t, tui.NewBounds(3, 4, 14, 15), a)
}
