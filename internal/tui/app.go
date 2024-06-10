package tui

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	app         *tview.Application
	file        string
	nodes       []Node
	nodesByName map[string]int
	pages       *tview.Pages
	graph       *tview.TextView
	tableDialog *tview.Grid
	table       *tview.Table
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
	net, nodes, err := bbn.FromFile(a.file)
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

	a.pages = tview.NewPages()

	mainPanel := a.createMainPanel()
	a.pages.AddAndSwitchToPage("Graph", mainPanel, true)

	a.tableDialog = a.createTablePanel()
	a.pages.AddPage("Table", a.tableDialog, true, false)

	mainPanel.SetInputCapture(a.inputMainPanel)
	a.table.SetInputCapture(a.inputTable)

	a.app.SetFocus(a.graph)

	if err := a.app.SetRoot(a.pages, true).Run(); err != nil {
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

	data := NewTable(a.nodes, 2, a.nodesByName)
	a.table = tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		SetContent(&data).
		SetEvaluateAllRows(true).
		SetFixed(1, 0).
		SetSeparator(tview.Borders.Vertical)
	a.table.SetBorder(true)

}

func (a *App) createMainPanel() *tview.Grid {
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
		SetText("Exit: ESC  Scroll: ←→↕  Nodes: Tab  Outcomes: Space/Numbers  Toggle: Enter  Table: T")
	grid.AddItem(help, 2, 0, 1, 2, 0, 0, false)

	return grid
}

func (a *App) createTablePanel() *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 16, 0)

	subGrid := tview.NewGrid().
		SetColumns(0).
		SetRows(0, 1)
	subGrid.SetBorder(true)

	subGrid.AddItem(a.table, 0, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Close table: ESC  Scroll: ←→↕")
	subGrid.AddItem(help, 1, 0, 1, 1, 0, 0, false)

	grid.AddItem(subGrid, 1, 1, 1, 1, 0, 0, false)

	return grid
}
