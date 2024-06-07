package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type App struct {
	file  string
	graph *tview.TextView
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
	_ = nodes

	a.createWidgets()

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
