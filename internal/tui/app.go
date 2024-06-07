package tui

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
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

func (a *App) input(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		a.app.Stop()
		return nil
	} else if event.Key() == tcell.KeyTAB {
		a.selectedNode = (a.selectedNode + 1) % len(a.nodes)
		a.selectedState = 0
		a.draw()
		return nil
	} else if event.Rune() == ' ' {
		a.selectedState = (a.selectedState + 1) % len(a.nodes[a.selectedNode].Node().States)
		a.draw()
		return nil
	} else if event.Key() == tcell.KeyEnter {
		node := a.nodes[a.selectedNode]
		value := node.Node().States[a.selectedState]

		oldEvidence := make(map[string]string, len(a.evidence))
		for k, v := range a.evidence {
			oldEvidence[k] = v
		}

		if oldValue, ok := a.evidence[node.Node().Name]; ok {
			if oldValue == value {
				delete(a.evidence, node.Node().Name)
			} else {
				a.evidence[node.Node().Name] = value
			}
		} else {
			a.evidence[node.Node().Name] = value
		}

		oldMarginals := a.marginals
		var err error
		a.marginals, err = a.network.Sample(a.evidence, a.samples, a.rng)
		if err != nil {
			if _, ok := err.(*bbn.ConflictingEvidenceError); ok {
				// TODO: show alert!
				a.evidence = oldEvidence
				a.marginals = oldMarginals
			} else {
				panic(err)
			}
		}
		a.draw()
		return nil
	} else if unicode.IsDigit(event.Rune()) {
		idx, err := strconv.Atoi(string(event.Rune()))
		if err != nil {
			panic(err)
		}
		idx -= 1
		if idx >= 0 && idx < len(a.nodes[a.selectedNode].Node().States) {
			a.selectedState = idx
			a.draw()
		}
		return nil
	}
	return event
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
		SetRows(0, 2).
		SetColumns(-1).
		SetBorders(false)

	grid.AddItem(a.graph, 0, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Scroll: ←→↕  Cycle nodes: TAB  Cycle states: SPACE/numbers\nToggle state: ENTER")
	grid.AddItem(help, 1, 0, 1, 1, 0, 0, false)

	return grid
}
