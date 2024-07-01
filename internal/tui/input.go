package tui

import (
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/ve"
	"github.com/rivo/tview"
)

func (a *App) inputMainPanel(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		// Exit program.
		a.app.Stop()
		return nil
	} else if event.Key() == tcell.KeyTAB {
		// Tab through nodes.
		a.selectedNode = (a.selectedNode + 1) % len(a.nodes)
		a.selectedState = 0
		a.render(true)
		return nil
	} else if event.Key() == tcell.KeyBacktab {
		// Tab through nodes, backwards.
		a.selectedNode--
		if a.selectedNode < 0 {
			a.selectedNode = len(a.nodes) - 1
		}
		a.selectedState = 0
		a.render(true)
		return nil
	} else if event.Rune() == ' ' {
		// Cycle through states.
		a.selectedState = (a.selectedState + 1) % len(a.nodes[a.selectedNode].Node().Outcomes)
		a.render(true)
		return nil
	} else if unicode.IsDigit(event.Rune()) {
		// Select states by index/number keys.
		idx, err := strconv.Atoi(string(event.Rune()))
		if err != nil {
			panic(err)
		}
		idx -= 1
		if idx >= 0 && idx < len(a.nodes[a.selectedNode].Node().Outcomes) {
			a.selectedState = idx
			a.render(true)
		}
		return nil
	} else if event.Key() == tcell.KeyEnter {
		// Set selected state as evidence.
		if err := a.inputEnter(); err != nil {
			panic(err)
		}
		a.render(true)
		return nil
	} else if event.Rune() == 't' {
		a.showTable()
	}
	return event
}

func (a *App) inputTable(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		a.pages.HidePage("Table")
		return nil
	}
	return event
}

// inputEnter adds the currently selected node and state to the evidence.
func (a *App) inputEnter() error {
	node := a.nodes[a.selectedNode]
	if node.Node().Type == ve.UtilityNode {
		return nil
	}

	value := node.Node().Outcomes[a.selectedState]

	// Store old evidence in case of fail/error.
	oldEvidence := make(map[string]string, len(a.evidence))
	for k, v := range a.evidence {
		oldEvidence[k] = v
	}

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

	// Store old marginals in case of fail/error.
	oldMarginals := a.marginals

	// Perform sampling
	var err error
	a.marginals, err = Solve(a.network, a.evidence, a.nodes)
	if err != nil {
		if _, ok := err.(*bbn.ConflictingEvidenceError); ok {
			// Rollback on error
			// TODO: show alert!
			a.evidence = oldEvidence
			a.marginals = oldMarginals
		} else {
			return err
		}
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

func (a *App) mouseInputGraph(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	if action == tview.MouseLeftClick || action == tview.MouseRightClick {
		front, _ := a.pages.GetFrontPage()
		if front == "Table" {
			a.pages.HidePage("Table")
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

func (a *App) mousePosInGraph(x, y int) (int, int) {
	boxX, boxY, _, _ := a.graph.GetInnerRect()
	scrollY, scrollX := a.graph.GetScrollOffset()
	return x - boxX + scrollX, y - boxY + scrollY
}
