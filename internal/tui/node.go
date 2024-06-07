package tui

import "github.com/mlange-42/bbn"

type Bounds struct {
	X int
	Y int
	W int
	H int
}

type Node interface {
	Bounds() Bounds
	Render(probs []float64) ([][]rune, [][]Color)
}

type node struct {
	node   *bbn.Node
	bounds Bounds
}

func NewNode(n *bbn.Node) Node {
	return &node{
		node: n,
		bounds: Bounds{
			X: n.Coords[0],
			Y: n.Coords[1],
			W: 20,
			H: len(n.States) + 3,
		},
	}
}

func (n *node) Bounds() Bounds {
	return n.bounds
}

func (n *node) Render(probs []float64) ([][]rune, [][]Color) {
	return nil, nil
}
