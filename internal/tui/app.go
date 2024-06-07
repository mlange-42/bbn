package tui

import (
	"math/rand"
	"strings"

	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	app     *tview.Application
	file    string
	nodes   []Node
	graph   *tview.TextView
	canvas  [][]rune
	network *bbn.Network
	rng     *rand.Rand
	samples int

	evidence      map[string]string
	marginals     map[string][]float64
	selectedNode  int
	selectedState int
}

func New(path string, samples int, seed int64) *App {
	return &App{
		file:     path,
		samples:  samples,
		rng:      rand.New(rand.NewSource(seed)),
		evidence: map[string]string{},
	}
}

func (a *App) Run() error {
	nodes, err := bbn.NodesFromYAML(a.file)
	if err != nil {
		return err
	}
	a.nodes = make([]Node, len(nodes))
	for i, n := range nodes {
		a.nodes[i] = NewNode(n)
	}

	a.network, err = bbn.New(nodes...)
	if err != nil {
		return err
	}
	a.marginals, err = a.network.Sample(a.evidence, a.samples, a.rng)
	if err != nil {
		return err
	}

	a.createCanvas()
	a.createWidgets()
	a.draw()

	a.app = tview.NewApplication()
	grid := a.createMainPanel()

	a.app.SetInputCapture(a.input)

	if err := a.app.SetRoot(grid, true).Run(); err != nil {
		return err
	}
	return nil
}

func (a *App) createCanvas() {
	bounds := NewBounds(0, 0, 1, 1)
	for _, n := range a.nodes {
		bounds.Extend(n.Bounds())
	}
	a.canvas = make([][]rune, bounds.H)
	for i := range a.canvas {
		a.canvas[i] = make([]rune, bounds.W)
		for j := range a.canvas[i] {
			a.canvas[i][j] = Empty
		}
	}
}

func (a *App) draw() {
	for i, node := range a.nodes {
		data := a.marginals[node.Node().Name]
		_, hasEvidence := a.evidence[node.Node().Name]
		runes, _ := node.Render(data, i == a.selectedNode, a.selectedState, hasEvidence)
		b := node.Bounds()
		for i, line := range runes {
			copy(a.canvas[b.Y+i][b.X:], line)
		}
	}

	b := strings.Builder{}
	for i, line := range a.canvas {
		b.WriteString(string(line))
		if i < len(a.canvas)-1 {
			b.WriteRune('\n')
		}
	}
	a.graph.SetText(b.String())
}

func (a *App) createWidgets() {
	a.graph = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText("")
	a.graph.SetBorder(true)
}

func (a *App) createMainPanel() tview.Primitive {
	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(-1).
		SetBorders(false)

	grid.AddItem(a.graph, 0, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Scroll: ←→↕  Nodes: Tab  States: Space/Numbers  Toggle: Enter")
	grid.AddItem(help, 1, 0, 1, 1, 0, 0, false)

	return grid
}
