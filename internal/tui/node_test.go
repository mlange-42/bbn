package tui_test

import (
	"strings"
	"testing"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/tui"
	"github.com/stretchr/testify/assert"
)

func TestNodeRender(t *testing.T) {
	node := &bbn.Node{
		Name:   "TestNode with a very long name",
		States: []string{"yes", "no", "maybe"},
		Coords: [2]int{0, 0},
	}

	uiNode := tui.NewNode(node)

	runes, _ := uiNode.Render([]float64{0.1, 0.2, 0.7})

	lines := make([]string, len(runes))
	for i, line := range runes {
		lines[i] = string(line)
	}
	text := strings.Join(lines, "\n")

	assert.Equal(t,
		`╔══════════════════════════╗
║ TestNode with a very lon ║
║ yes   █░░░░░░░░░  10.000 ║
║ no    ██░░░░░░░░  20.000 ║
║ maybe ███████░░░  70.000 ║
╚══════════════════════════╝`, text)
}
