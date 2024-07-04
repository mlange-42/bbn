package tui

import (
	"fmt"

	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	app          *tview.Application
	file         string
	trainingFile string
	csvDelimiter rune
	noData       string

	nodes       []Node
	nodesByName map[string]int
	pages       *tview.Pages
	graph       *tview.TextView
	tableDialog *tview.Grid
	table       *tview.Table
	helpDialog  *tview.Grid
	help        *tview.TextView
	infoDialog  *tview.Grid
	info        *tview.TextView
	canvas      [][]rune
	colors      [][]Color
	network     *bbn.Network

	evidence      map[string]string
	marginals     map[string][]float64
	selectedNode  int
	selectedState int

	ignorePolicies bool
}

func New(path string, evidence map[string]string, trainingFile, noData string, csvDelimiter rune) *App {
	if evidence == nil {
		evidence = map[string]string{}
	}
	return &App{
		file:         path,
		trainingFile: trainingFile,
		csvDelimiter: csvDelimiter,
		evidence:     evidence,
	}
}

func (a *App) Run() error {
	net, nodes, err := bbn.FromFile(a.file)
	if err != nil {
		return err
	}

	a.network = net

	if a.trainingFile != "" {
		a.network, err = TrainNetwork(net, nodes, a.trainingFile, a.noData, a.csvDelimiter)
		if err != nil {
			return err
		}
	}

	a.nodes = make([]Node, len(nodes))
	a.nodesByName = make(map[string]int, len(nodes))
	for i, n := range nodes {
		a.nodes[i] = NewNode(n)
		a.nodesByName[n.Name] = i
	}

	policy, err := a.network.SolvePolicies(true)
	if err != nil {
		return err
	}

	for name, f := range policy {
		idx := a.nodesByName[name]
		node := a.nodes[idx]
		node.Node().Factor = &f
	}

	err = a.updateMarginals()
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

	a.helpDialog = a.createHelpPanel()
	a.pages.AddPage("Help", a.helpDialog, true, false)

	a.infoDialog = a.createInfoPanel()
	a.pages.AddPage("Info", a.infoDialog, true, false)

	mainPanel.SetInputCapture(a.inputMainPanel)
	a.graph.SetMouseCapture(a.mouseInputGraph)

	a.table.SetInputCapture(a.inputTable)
	a.table.SetMouseCapture(a.mouseInputTable)

	a.help.SetInputCapture(a.inputHelp)
	a.help.SetMouseCapture(a.mouseInputHelp)

	a.info.SetInputCapture(a.inputInfo)
	a.info.SetMouseCapture(a.mouseInputInfo)

	rooted := a.app.SetRoot(a.pages, true)
	a.showInfo()

	if err := rooted.Run(); err != nil {
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

	a.help = tview.NewTextView().
		SetWrap(true).
		SetText(` Set/unset evidence by clicking on the probability bars of nodes.

                    Keyboard              Mouse

 Exit               ESC
 Scroll             ←→↕
 Navigate nodes     Tab/BackTab
 Navigate bars      Space/Numbers
 Toggle evidence    Enter                 left click
 Show table         T                     right click
 Ignore policies    P
 Move node          W/A/S/D
 Save network       Ctrl+S
`)
	a.info = tview.NewTextView().
		SetWrap(true).
		SetText(a.network.Info())
}

func (a *App) createMainPanel() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(len(a.canvas[0])+3, 0).
		SetBorders(false)

	header := tview.NewTextView().
		SetWrap(false).
		SetText(fmt.Sprintf("BBNi - %s (%s)", a.network.Name(), a.file))
	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)

	grid.AddItem(a.graph, 1, 0, 1, 2, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Press H for help, I for network info!")
	grid.AddItem(help, 2, 0, 1, 2, 0, 0, false)

	return grid
}

func (a *App) createTablePanel() *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, 72, 0).
		SetRows(0, 17, 0)

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

func (a *App) createHelpPanel() *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, 72, 0).
		SetRows(0, 17, 0)

	subGrid := tview.NewGrid().
		SetColumns(0).
		SetRows(0, 1)
	subGrid.SetBorder(true)
	subGrid.SetTitle(" Help ")

	subGrid.AddItem(a.help, 0, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText(" Close help: ESC  Scroll: ←→↕")
	subGrid.AddItem(help, 1, 0, 1, 1, 0, 0, false)

	grid.AddItem(subGrid, 1, 1, 1, 1, 0, 0, false)

	return grid
}

func (a *App) createInfoPanel() *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, 72, 0).
		SetRows(0, 17, 0)

	subGrid := tview.NewGrid().
		SetColumns(0).
		SetRows(0, 1)
	subGrid.SetBorder(true)
	subGrid.SetTitle(" Info ")

	subGrid.AddItem(a.info, 0, 0, 1, 1, 0, 0, true)

	info := tview.NewTextView().
		SetWrap(false).
		SetText(" Close info: ESC  Scroll: ←→↕")
	subGrid.AddItem(info, 1, 0, 1, 1, 0, 0, false)

	grid.AddItem(subGrid, 1, 1, 1, 1, 0, 0, false)

	return grid
}
