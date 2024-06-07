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
	Render(probs []float64) ([][]rune, [][]Color)
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
			runes[i][j] = BorderNone
		}
	}

	node := node{
		node:   n,
		bounds: bounds,
		runes:  runes,
		colors: colors,
		barsX:  maxStateLen + 3,
	}

	node.drawBorder()
	node.drawTitle()
	node.drawStateLabels()

	node.drawBars(make([]float64, len(n.States)))

	return &node
}

func (n *node) Node() *bbn.Node {
	return n.node
}

func (n *node) Bounds() Bounds {
	return n.bounds
}

func (n *node) Render(probs []float64) ([][]rune, [][]Color) {
	n.drawBars(probs)
	return n.runes, n.colors
}

func (n *node) drawBorder() {
	n.runes[0][0] = BorderNW
	n.runes[0][n.bounds.W-1] = BorderNE
	n.runes[n.bounds.H-1][0] = BorderSW
	n.runes[n.bounds.H-1][n.bounds.W-1] = BorderSE
	for i := 1; i < n.bounds.W-1; i++ {
		n.runes[0][i] = BorderH
		n.runes[n.bounds.H-1][i] = BorderH
	}
	for i := 1; i < n.bounds.H-1; i++ {
		n.runes[i][0] = BorderV
		n.runes[i][n.bounds.W-1] = BorderV
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

func (n *node) drawBars(probs []float64) {
	for i, p := range probs {
		full, _ := math.Modf(p * 10)
		for j := 0; j < int(full); j++ {
			n.runes[i+2][n.barsX+j] = Full
		}
		for j := int(full); j < maxBars; j++ {
			n.runes[i+2][n.barsX+j] = Shade
		}
		text := []rune(fmt.Sprintf("%7.3f", p*100))
		copy(n.runes[i+2][n.barsX+maxBars+1:], text)
	}
}
