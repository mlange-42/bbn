package tui_test

import (
	"strings"
	"testing"

	"github.com/mlange-42/bbn/internal/tui"
	"github.com/mlange-42/bbn/net"
	"github.com/stretchr/testify/assert"
)

func TestNodeRender(t *testing.T) {
	node := net.Variable{
		Name:     "TestNode with a very long name",
		Outcomes: []string{"yes", "no", "maybe"},
		Position: [2]int{0, 0},
	}

	uiNode := tui.NewNode(node)

	runes, _ := uiNode.Render([]float64{0.1, 0.2, 0.7}, true, 1, false)

	lines := make([]string, len(runes))
	for i, line := range runes {
		lines[i] = string(line)
	}
	text := strings.Join(lines, "\n")

	assert.Equal(t,
		`╔══════════════════════════╗
║ TestNode with a very lon ║
║ yes   █░░░░░░░░░  10.000 ║
║[no    ██░░░░░░░░  20.000]║
║ maybe ███████░░░  70.000 ║
╚══════════════════════════╝`, text)
}
