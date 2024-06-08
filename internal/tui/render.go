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

	hOverlap := b1.X+b1.W >= b2.X && b1.X <= b2.X+b2.W
	vOverlap := b1.Y+b1.H >= b2.Y && b1.Y <= b2.Y+b2.H
	if hOverlap {
		if !vOverlap {
			a.renderVertical(b1, b2)
		}
	} else if vOverlap {
		if !hOverlap {
			a.renderHorizontal(b1, b2)
		}
	}
}

func (a *App) renderVertical(b1, b2 Bounds) {
	xMid := (max(b1.X, b2.X) + min(b1.X+b1.W, b2.X+b2.W)) / 2
	if b1.Y < b2.Y {
		for y := b1.Y + b1.H; y < b2.Y; y++ {
			a.canvas[y][xMid] = BorderH[0]
		}
		a.canvas[b2.Y-1][xMid] = ArrowRight
	} else {
		for y := b2.Y + b2.H; y < b1.Y; y++ {
			a.canvas[y][xMid] = BorderH[0]
		}
		a.canvas[b1.Y-1][xMid] = ArrowLeft
	}
}

func (a *App) renderHorizontal(b1, b2 Bounds) {
	yMid := (max(b1.Y, b2.Y) + min(b1.Y+b1.H, b2.Y+b2.H)) / 2
	if b1.X < b2.X {
		for x := b1.X + b1.W; x < b2.X; x++ {
			a.canvas[yMid][x] = BorderV[0]
		}
		a.canvas[yMid][b2.Y-1] = ArrowDown
	} else {
		for x := b2.X + b2.W; x < b1.X; x++ {
			a.canvas[yMid][x] = BorderV[0]
		}
		a.canvas[yMid][b1.X-1] = ArrowUp
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
