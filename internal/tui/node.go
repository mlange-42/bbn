package tui

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/mlange-42/bbn"
)

const maxStateLabelWidth = 6
const maxBars = 10

type Node interface {
	Node() *bbn.Node
	Bounds() Bounds
	Render(probs []float64, selected bool, state int, evidence bool) ([][]rune, [][]Color)
}

type node struct {
	node   *bbn.Node
	bounds Bounds
	runes  [][]rune
	colors [][]Color
	barsX  int
}

func NewNode(n *bbn.Node) Node {
	maxStateLen := 0
	for _, state := range n.States {
		cnt := utf8.RuneCountInString(state)
		if cnt > maxStateLen {
			maxStateLen = cnt
		}
	}
	if maxStateLen > maxStateLabelWidth {
		maxStateLen = maxStateLabelWidth
	}

	bounds := Bounds{
		X: n.Coords[0],
		Y: n.Coords[1],
		W: maxStateLen + maxBars + 7 + 6,
		H: len(n.States) + 3,
	}
	runes := make([][]rune, bounds.H)
	colors := make([][]Color, bounds.H)
	for i := range runes {
		runes[i] = make([]rune, bounds.W)
		colors[i] = make([]Color, bounds.W)
		for j := range runes[i] {
			runes[i][j] = Empty
		}
	}

	node := node{
		node:   n,
		bounds: bounds,
		runes:  runes,
		colors: colors,
		barsX:  maxStateLen + 3,
	}

	node.drawBorder(false)
	node.drawTitle()
	node.drawStateLabels()

	node.drawBars(make([]float64, len(n.States)), false, 0, false)

	return &node
}

func (n *node) Node() *bbn.Node {
	return n.node
}

func (n *node) Bounds() Bounds {
	return n.bounds
}

func (n *node) Render(probs []float64, selected bool, state int, evidence bool) ([][]rune, [][]Color) {
	n.drawBorder(selected)
	n.drawBars(probs, selected, state, evidence)
	return n.runes, n.colors
}

func (n *node) drawBorder(selected bool) {
	style := 0
	if selected {
		style = 1
	}

	n.runes[0][0] = BorderNW[style]
	n.runes[0][n.bounds.W-1] = BorderNE[style]
	n.runes[n.bounds.H-1][0] = BorderSW[style]
	n.runes[n.bounds.H-1][n.bounds.W-1] = BorderSE[style]
	for i := 1; i < n.bounds.W-1; i++ {
		n.runes[0][i] = BorderH[style]
		n.runes[n.bounds.H-1][i] = BorderH[style]
	}
	for i := 1; i < n.bounds.H-1; i++ {
		n.runes[i][0] = BorderV[style]
		n.runes[i][n.bounds.W-1] = BorderV[style]
	}
}

func (n *node) drawTitle() {
	runes := []rune(n.node.Name)
	copy(n.runes[1][2:n.bounds.W-2], runes)
}

func (n *node) drawStateLabels() {
	for i, label := range n.node.States {
		copy(n.runes[i+2][2:n.barsX-1], []rune(label))
	}
}

func (n *node) drawBars(probs []float64, selected bool, state int, evidence bool) {
	for i, p := range probs {
		full, frac := math.Modf(p * 10)
		for j := 0; j < int(full); j++ {
			n.runes[i+2][n.barsX+j] = Full
		}
		fracIdx := int(frac * float64(len(Partial)))
		if fracIdx > 0 {
			n.runes[i+2][n.barsX+int(full)] = Partial[fracIdx]
			full++
		}

		for j := int(full); j < maxBars; j++ {
			n.runes[i+2][n.barsX+j] = Shade
		}
		text := []rune(fmt.Sprintf("%7.3f", p*100))
		copy(n.runes[i+2][n.barsX+maxBars+1:], text)

		if selected && state == i {
			n.runes[i+2][1] = SelectionStart
			n.runes[i+2][n.bounds.W-2] = SelectionEnd
		} else if evidence {
			n.runes[i+2][1] = EvidenceStart
			n.runes[i+2][n.bounds.W-2] = EvidenceEnd
		} else {
			n.runes[i+2][1] = Empty
			n.runes[i+2][n.bounds.W-2] = Empty
		}
	}
}
