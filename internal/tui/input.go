package tui

import (
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/ve"
	"github.com/rivo/tview"
)

func (a *App) inputMainPanel(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		// Exit program.
		a.app.Stop()
		return nil
	} else if event.Key() == tcell.KeyTAB {
		// Tab through nodes.
		a.selectNextNode()
		return nil
	} else if event.Key() == tcell.KeyBacktab {
		// Tab through nodes, backwards.
		a.selectPreviousNode()
		return nil
	} else if event.Rune() == ' ' {
		// Cycle through states.
		a.selectedState = (a.selectedState + 1) % len(a.nodes[a.selectedNode].Node().Outcomes)
		a.render(true)
		return nil
	} else if unicode.IsDigit(event.Rune()) {
		// Select states by index/number keys.
		a.selectNodeOutcome(string(event.Rune()))
		return nil
	} else if event.Key() == tcell.KeyEnter {
		// Set selected state as evidence.
		if err := a.inputEnter(); err != nil {
			panic(err)
		}
		a.render(true)
		return nil
	} else if event.Rune() == 'h' {
		a.showHelp()
		return nil
	} else if event.Rune() == 'i' {
		a.showInfo()
		return nil
	} else if event.Rune() == 't' {
		a.showTable()
		return nil
	} else if event.Rune() == 'p' {
		a.toggleIgnorePolicy()
		return nil
	} else if event.Key() == tcell.KeyCtrlS {
		a.saveNetwork()
		return nil
	} else {
		return a.moveNode(event)
	}
}

func (a *App) toggleIgnorePolicy() {
	a.ignorePolicies = !a.ignorePolicies

	if a.ignorePolicies {
		a.graph.SetBorderColor(tcell.ColorBlue)
	} else {
		a.graph.SetBorderColor(tcell.ColorDefault)
	}

	a.updateMarginals()
	a.render(true)
}

func (a *App) saveNetwork() {
	yml, err := bbn.ToYAML(a.network)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(a.file)
	saveFile := strings.TrimSuffix(a.file, ext) + "-save.yml"

	err = os.WriteFile(saveFile, yml, 0644)
	if err != nil {
		panic(err)
	}
}

func (a *App) moveNode(event *tcell.EventKey) *tcell.EventKey {
	node := a.nodes[a.selectedNode]

	dx, dy := 0, 0
	createCanvas := false
	if event.Rune() == 'w' {
		if node.Bounds().Y > 0 {
			dy--
		}
	} else if event.Rune() == 'a' {
		if node.Bounds().X > 0 {
			dx--
		}
	} else if event.Rune() == 's' {
		dy++
		if node.Bounds().Y+node.Bounds().H >= len(a.canvas) {
			createCanvas = true
		}
	} else if event.Rune() == 'd' {
		dx++
		if node.Bounds().X+node.Bounds().W >= len(a.canvas[0]) {
			createCanvas = true
		}
	}
	if dx != 0 || dy != 0 {
		b := node.Bounds()
		b.X += dx
		b.Y += dy

		p := a.network.Variables()[a.selectedNode].Position
		a.network.Variables()[a.selectedNode].Position = [2]int{p[0] + dx, p[1] + dy}

		if createCanvas {
			a.createCanvas()
		}

		a.render(false)
		return nil
	}
	return event
}

// inputEnter adds the currently selected node and state to the evidence.
func (a *App) inputEnter() error {
	node := a.nodes[a.selectedNode]
	if node.Node().NodeType == ve.UtilityNode {
		return nil
	}

	value := node.Node().Outcomes[a.selectedState]

	// Add/clear selected state
	if oldValue, ok := a.evidence[node.Node().Name]; ok {
		if oldValue == value {
			delete(a.evidence, node.Node().Name)
		} else {
			a.evidence[node.Node().Name] = value
		}
	} else {
		a.evidence[node.Node().Name] = value
	}

	return a.updateMarginals()
}

func (a *App) updateMarginals() error {
	var err error
	a.marginals, err = Solve(a.network, a.evidence, a.nodes, a.ignorePolicies)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) showTable() {
	node := a.nodes[a.selectedNode].Node()

	data := NewTable(a.nodes, a.selectedNode, a.nodesByName)

	a.tableDialog.SetRows(0, data.GetRowCount()+5, 0)
	a.table.SetContent(&data)
	a.table.SetTitle(" " + node.Name + " ")
	a.pages.ShowPage("Table")

	a.app.SetFocus(a.table)
}

func (a *App) showHelp() {
	a.pages.ShowPage("Help")
	a.app.SetFocus(a.help)
}

func (a *App) showInfo() {
	a.pages.ShowPage("Info")
	a.app.SetFocus(a.info)
}

func (a *App) mouseInputGraph(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	if action == tview.MouseLeftClick || action == tview.MouseRightClick {
		front, _ := a.pages.GetFrontPage()
		if front == "Table" {
			a.pages.HidePage("Table")
			return tview.MouseConsumed, nil
		} else if front == "Help" {
			a.pages.HidePage("Help")
			return tview.MouseConsumed, nil
		} else if front == "Info" {
			a.pages.HidePage("Info")
			return tview.MouseConsumed, nil
		}
	}

	if action == tview.MouseLeftClick {
		x, y := a.mousePosInGraph(event.Position())
		for i, node := range a.nodes {
			if node.Bounds().Contains(x, y) {
				a.selectedNode = i
				if outcome, ok := node.SelectedOutcome(x, y); ok {
					a.selectedState = outcome
					if err := a.inputEnter(); err != nil {
						panic(err)
					}
				}
				a.render(true)
				break
			}
		}
		return tview.MouseConsumed, nil
	} else if action == tview.MouseRightClick {
		x, y := a.mousePosInGraph(event.Position())
		for i, node := range a.nodes {
			if node.Bounds().Contains(x, y) {
				a.selectedNode = i
				a.render(true)
				a.showTable()
				break
			}
		}
		return tview.MouseConsumed, nil
	}

	return action, event
}

func (a *App) mouseInputTable(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	if action == tview.MouseLeftClick || action == tview.MouseRightClick {
		a.pages.HidePage("Table")
		return tview.MouseConsumed, nil
	}
	return action, event
}

func (a *App) mouseInputHelp(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	if action == tview.MouseLeftClick || action == tview.MouseRightClick {
		a.pages.HidePage("Help")
		return tview.MouseConsumed, nil
	}
	return action, event
}

func (a *App) mouseInputInfo(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	if action == tview.MouseLeftClick || action == tview.MouseRightClick {
		a.pages.HidePage("Info")
		return tview.MouseConsumed, nil
	}
	return action, event
}

func (a *App) inputTable(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
		a.pages.HidePage("Table")
		return nil
	}
	return event
}

func (a *App) inputHelp(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
		a.pages.HidePage("Help")
		return nil
	}
	return event
}

func (a *App) inputInfo(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
		a.pages.HidePage("Info")
		return nil
	}
	return event
}

func (a *App) mousePosInGraph(x, y int) (int, int) {
	boxX, boxY, _, _ := a.graph.GetInnerRect()
	scrollY, scrollX := a.graph.GetScrollOffset()
	return x - boxX + scrollX, y - boxY + scrollY
}

func (a *App) selectNextNode() {
	a.selectedNode = (a.selectedNode + 1) % len(a.nodes)
	a.selectedState = 0
	a.render(true)
}

func (a *App) selectPreviousNode() {
	a.selectedNode--
	if a.selectedNode < 0 {
		a.selectedNode = len(a.nodes) - 1
	}
	a.selectedState = 0
	a.render(true)
}

func (a *App) selectNodeOutcome(index string) {
	idx, err := strconv.Atoi(index)
	if err != nil {
		panic(err)
	}
	idx -= 1
	if idx >= 0 && idx < len(a.nodes[a.selectedNode].Node().Outcomes) {
		a.selectedState = idx
		a.render(true)
	}
}
