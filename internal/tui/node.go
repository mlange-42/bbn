package tui

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/ve"
)

const maxNodeLabelWidth = 48
const maxStateLabelWidth = 16
const maxBars = 10
const extraUtilityWidth = 5

type Node interface {
	Node() *bbn.Variable
	Bounds() *Bounds
	Render(probs []float64, selected bool, state int, evidence bool) ([][]rune, [][]Color)
	SelectedOutcome(x, y int) (int, bool)
}

type node struct {
	node   bbn.Variable
	bounds Bounds
	runes  [][]rune
	colors [][]Color
	barsX  int
}

func NewNode(n bbn.Variable) Node {
	maxStateLen := 0
	for _, state := range n.Outcomes {
		cnt := utf8.RuneCountInString(state)
		if cnt > maxStateLen {
			maxStateLen = cnt
		}
	}
	if maxStateLen > maxStateLabelWidth {
		maxStateLen = maxStateLabelWidth
	}
	titleLength := min(utf8.RuneCountInString(n.Name), maxNodeLabelWidth)
	bars := maxBars
	if n.NodeType == ve.UtilityNode {
		bars = extraUtilityWidth
	}

	bounds := Bounds{
		X: n.Position[0],
		Y: n.Position[1],
		W: max(maxStateLen+bars+9, titleLength) + 4,
		H: len(n.Outcomes) + 3,
	}

	if n.NodeType == ve.UtilityNode {
		bounds.H++
	}

	color, err := nodeColor(&n)
	if err != nil {
		panic(err)
	}

	runes := make([][]rune, bounds.H)
	colors := make([][]Color, bounds.H)
	for i := range runes {
		runes[i] = make([]rune, bounds.W)
		colors[i] = make([]Color, bounds.W)
		for j := range runes[i] {
			runes[i][j] = Empty
		}
		for j := range colors[i] {
			colors[i][j] = color
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

	return &node
}

func nodeColor(n *bbn.Variable) (Color, error) {
	color := White
	if n.Color == "" {
		switch n.NodeType {
		case ve.ChanceNode:
			color = White
		case ve.UtilityNode:
			color = Green
		case ve.DecisionNode:
			color = Blue
		}
	} else {
		var ok bool
		color, ok = NamedColors[n.Color]
		if !ok {
			return Black, fmt.Errorf("unknown node color %s; available colors are:\n %s", n.Color, strings.Join(ColorNames, ", "))
		}
	}
	return color, nil
}

func (n *node) Node() *bbn.Variable {
	return &n.node
}

func (n *node) Bounds() *Bounds {
	return &n.bounds
}

func (n *node) Render(probs []float64, selected bool, state int, evidence bool) ([][]rune, [][]Color) {
	n.drawBorder(selected)
	if n.node.NodeType == ve.UtilityNode {
		n.drawUtility(probs)
	} else {
		n.drawBars(probs, selected, state, evidence)
	}
	return n.runes, n.colors
}

func (n *node) SelectedOutcome(x, y int) (int, bool) {
	if !n.bounds.Contains(x, y) {
		return -1, false
	}
	for i := range n.node.Outcomes {
		if y == n.bounds.Y+i+2 {
			return i, true
		}
	}
	return -1, false
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

	if n.node.NodeType == ve.DecisionNode {
		n.runes[1][n.bounds.W-3] = '!'
	} else if n.node.NodeType == ve.UtilityNode {
		n.runes[1][n.bounds.W-3] = '$'
	}
}

func (n *node) drawStateLabels() {
	for i, label := range n.node.Outcomes {
		copy(n.runes[i+2][2:n.barsX-1], []rune(label))
	}
	if n.node.NodeType == ve.UtilityNode {
		copy(n.runes[len(n.node.Outcomes)+2][2:n.barsX-1], []rune("total"))
	}
}

func (n *node) drawBars(probs []float64, selected bool, state int, evidence bool) {
	for i, p := range probs {
		var full, frac float64
		if !math.IsNaN(p) {
			full, frac = math.Modf(p * 10)
		}
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
		var text []rune
		if math.IsNaN(p) {
			text = []rune("    ???")
		} else {
			text = []rune(fmt.Sprintf("%7.3f", p*100))
		}

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

func (n *node) drawUtility(probs []float64) {
	for i, p := range probs {
		var text []rune
		if math.IsNaN(p) {
			text = []rune("         ???")
		} else {
			text = []rune(fmt.Sprintf("%12.3f", p))
		}
		copy(n.runes[i+2][n.barsX+1:], text)
	}
}
