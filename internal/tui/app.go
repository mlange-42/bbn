package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	file   string
	nodes  []Node
	graph  *tview.TextView
	canvas [][]rune
}

func New(path string) *App {
	return &App{
		file: path,
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

	a.createCanvas()
	a.createWidgets()
	a.draw()

	app := tview.NewApplication()
	grid := a.createMainPanel()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}
		return event
	})
	if err := app.SetRoot(grid, true).Run(); err != nil {
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
			a.canvas[i][j] = BorderNone
		}
	}
}

func (a *App) draw() {
	for _, node := range a.nodes {
		runes, _ := node.Render(make([]float64, len(node.Node().States)))
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
		SetText("Exit: ESC  Scroll: ←→↕")
	grid.AddItem(help, 1, 0, 1, 1, 0, 0, false)

	return grid
}
