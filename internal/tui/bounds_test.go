package tui_test

import (
	"testing"

	"github.com/mlange-42/bbn/internal/tui"
	"github.com/stretchr/testify/assert"
)

func TestBounds(t *testing.T) {
	a := tui.NewBounds(3, 4, 5, 6)
	b := tui.NewBounds(1, 2, 12, 13)
	a.Extend(b)
	assert.Equal(t, b, a)

	a = tui.NewBounds(3, 4, 5, 6)
	b = tui.NewBounds(5, 6, 12, 13)
	a.Extend(b)
	assert.Equal(t, tui.NewBounds(3, 4, 14, 15), a)
}
