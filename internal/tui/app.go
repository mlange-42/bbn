package tui

import (
	"fmt"

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
	colors      [][]Color
	network     *bbn.Network

	evidence      map[string]string
	marginals     map[string][]float64
	selectedNode  int
	selectedState int
}

func New(path string, evidence map[string]string) *App {
	if evidence == nil {
		evidence = map[string]string{}
	}
	return &App{
		file:     path,
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
		a.nodesByName[n.Name] = i
	}

	policy, err := a.network.SolvePolicies(false)
	if err != nil {
		return err
	}

	for name, f := range policy {
		idx := a.nodesByName[name]
		node := a.nodes[idx]
		node.Node().Factor = &f
	}

	a.marginals, err = Solve(a.network, a.evidence, a.nodes)
	if err != nil {
		return err
	}

	a.createCanvas()
	a.createWidgets()
	a.render(false)

	a.app = tview.NewApplication()
	a.app.EnableMouse(true)

	a.pages = tview.NewPages()

	mainPanel := a.createMainPanel()
	a.pages.AddAndSwitchToPage("Graph", mainPanel, true)

	a.tableDialog = a.createTablePanel()
	a.pages.AddPage("Table", a.tableDialog, true, false)

	mainPanel.SetInputCapture(a.inputMainPanel)
	a.table.SetInputCapture(a.inputTable)
	a.table.SetMouseCapture(a.mouseInputTable)
	a.graph.SetMouseCapture(a.mouseInputGraph)

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
	a.colors = make([][]Color, bounds.H)
	for i := range a.canvas {
		a.canvas[i] = make([]rune, bounds.W)
		a.colors[i] = make([]Color, bounds.W)
		for j := range a.canvas[i] {
			a.canvas[i][j] = Empty
			a.colors[i][j] = White
		}
	}
}

func (a *App) createWidgets() {
	a.graph = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText("")
	a.graph.SetBorder(true)

	a.table = tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false).
		SetEvaluateAllRows(true).
		SetFixed(1, 0).
		SetSeparator(tview.Borders.Vertical)
	a.table.SetBorder(true)

}

func (a *App) createMainPanel() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(len(a.canvas[0])+3, 0).
		SetBorders(false)

	header := tview.NewTextView().
		SetWrap(false).
		SetText(fmt.Sprintf("BBNi - %s (%s)", a.network.Name(), a.file))
	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)

	grid.AddItem(a.graph, 1, 0, 1, 2, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Scroll: ←→↕  Navigate: Tab, Space/Numbers\nToggle outcome: Enter/LeftClick  Show table: T/RightClick")
	grid.AddItem(help, 2, 0, 1, 2, 0, 0, false)

	return grid
}

func (a *App) createTablePanel() *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, 72, 0).
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
