package tui

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	app         *tview.Application
	file        string
	nodes       []Node
	nodesByName map[string]int
	graph       *tview.TextView
	canvas      [][]rune
	network     *bbn.Network
	rng         *rand.Rand
	samples     int

	evidence      map[string]string
	marginals     map[string][]float64
	selectedNode  int
	selectedState int
}

func New(path string, evidence map[string]string, samples int, seed int64) *App {
	if evidence == nil {
		evidence = map[string]string{}
	}
	return &App{
		file:     path,
		samples:  samples,
		rng:      rand.New(rand.NewSource(seed)),
		evidence: evidence,
	}
}

func (a *App) Run() error {
	yml, err := os.ReadFile(a.file)
	if err != nil {
		return err
	}

	net, nodes, err := bbn.FromYAML(yml)
	if err != nil {
		return err
	}
	a.network = net

	a.nodes = make([]Node, len(nodes))
	a.nodesByName = make(map[string]int, len(nodes))
	for i, n := range nodes {
		a.nodes[i] = NewNode(n)
		a.nodesByName[n.Variable] = i
	}

	a.marginals, err = a.network.Sample(a.evidence, a.samples, a.rng)
	if err != nil {
		return err
	}

	a.createCanvas()
	a.createWidgets()
	a.render(false)

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

func (a *App) createWidgets() {
	a.graph = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText("")
	a.graph.SetBorder(true)
}

func (a *App) createMainPanel() tview.Primitive {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(len(a.canvas[0])+3, 0).
		SetBorders(false)

	header := tview.NewTextView().
		SetWrap(false).
		SetText(fmt.Sprintf("BBNi - %s", a.file))
	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)

	grid.AddItem(a.graph, 1, 0, 1, 2, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Scroll: ←→↕  Nodes: Tab  Outcomes: Space/Numbers  Toggle: Enter")
	grid.AddItem(help, 2, 0, 1, 2, 0, 0, false)

	return grid
}
