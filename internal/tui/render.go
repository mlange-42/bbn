package tui

import (
	"strings"
)

func (a *App) render() {
	a.renderNodes()
	a.renderEdges()
	a.updateCanvas()
}

func (a *App) renderNodes() {
	for i, node := range a.nodes {
		data := a.marginals[node.Node().Name]
		_, hasEvidence := a.evidence[node.Node().Name]
		runes, _ := node.Render(data, i == a.selectedNode, a.selectedState, hasEvidence)
		b := node.Bounds()
		for i, line := range runes {
			copy(a.canvas[b.Y+i][b.X:], line)
		}
	}
}

func (a *App) renderEdges() {
	for i, node := range a.nodes {
		for _, p := range node.Node().Parents {
			pid := a.nodesByName[p]
			a.renderEdge(pid, i)
		}
	}
}

func (a *App) renderEdge(from, to int) {
	n1, n2 := a.nodes[from], a.nodes[to]
	b1, b2 := n1.Bounds(), n2.Bounds()

	hOverlap := b1.X+b1.W > b2.X && b1.X < b2.X+b2.W
	vOverlap := b1.Y+b1.H > b2.Y && b1.Y < b2.Y+b2.H
	if hOverlap {
		if !vOverlap {
			a.renderEdgeVertical(b1, b2)
		}
	} else if vOverlap {
		if !hOverlap {
			a.renderEdgeHorizontal(b1, b2)
		}
	} else {
		a.renderEdgeCorner(b1, b2)
	}
}

func (a *App) renderEdgeVertical(b1, b2 Bounds) {
	xMid := (max(b1.X, b2.X) + min(b1.X+b1.W, b2.X+b2.W)) / 2
	if b1.Y < b2.Y {
		for y := b1.Y + b1.H; y < b2.Y; y++ {
			a.canvas[y][xMid] = BorderV[0]
		}
		a.canvas[b2.Y-1][xMid] = ArrowDown
	} else {
		for y := b2.Y + b2.H; y < b1.Y; y++ {
			a.canvas[y][xMid] = BorderV[0]
		}
		a.canvas[b2.Y+b2.H][xMid] = ArrowUp
	}
}

func (a *App) renderEdgeHorizontal(b1, b2 Bounds) {
	yMid := (max(b1.Y, b2.Y) + min(b1.Y+b1.H, b2.Y+b2.H)) / 2
	if b1.X < b2.X {
		for x := b1.X + b1.W; x < b2.X; x++ {
			a.canvas[yMid][x] = BorderH[0]
		}
		a.canvas[yMid][b2.X-1] = ArrowRight
	} else {
		for x := b2.X + b2.W; x < b1.X; x++ {
			a.canvas[yMid][x] = BorderH[0]
		}
		a.canvas[yMid][b1.X-1] = ArrowLeft
	}
}

func (a *App) renderEdgeCorner(b1, b2 Bounds) {
	downwards := b2.Y > b1.Y
	rightwards := b2.X > b1.X

	yStart := b1.Y
	if downwards {
		yStart = b1.Y + b1.H - 1
	}
	var xStart int
	if rightwards {
		for x := b1.X + b1.W + 1; x < b2.X; x++ {
			a.canvas[yStart][x] = BorderH[0]
		}
		xStart = b2.X
		if downwards {
			a.canvas[yStart][xStart] = BorderNE[0]
		} else {
			a.canvas[yStart][xStart] = BorderSE[0]
		}
	} else {
		for x := b1.X - 1; x >= b2.X+b2.W; x-- {
			a.canvas[yStart][x] = BorderH[0]
		}
		xStart = b2.X + b2.W - 1
		if downwards {
			a.canvas[yStart][xStart] = BorderNW[0]
		} else {
			a.canvas[yStart][xStart] = BorderSW[0]
		}
	}

	if downwards {
		for y := yStart + 1; y < b2.Y; y++ {
			a.canvas[y][xStart] = BorderV[0]
		}
		a.canvas[b2.Y-1][xStart] = ArrowDown
	} else {
		for y := yStart - 1; y > b2.Y+b2.H; y-- {
			a.canvas[y][xStart] = BorderV[0]
		}
		a.canvas[b2.Y+b2.H][xStart] = ArrowUp
	}
}

func (a *App) updateCanvas() {
	b := strings.Builder{}
	for i, line := range a.canvas {
		b.WriteString(string(line))
		if i < len(a.canvas)-1 {
			b.WriteRune('\n')
		}
	}
	a.graph.SetText(b.String())
}
