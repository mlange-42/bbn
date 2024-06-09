package tui

import (
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/bbn"
)

func (a *App) input(event *tcell.EventKey) *tcell.EventKey {
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
		a.selectedState = (a.selectedState + 1) % len(a.nodes[a.selectedNode].Node().States)
		a.render(true)
		return nil
	} else if unicode.IsDigit(event.Rune()) {
		// Select states by index/number keys.
		idx, err := strconv.Atoi(string(event.Rune()))
		if err != nil {
			panic(err)
		}
		idx -= 1
		if idx >= 0 && idx < len(a.nodes[a.selectedNode].Node().States) {
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
	}
	return event
}

// inputEnter adds the currently selected node and state to the evidence.
func (a *App) inputEnter() error {
	node := a.nodes[a.selectedNode]
	value := node.Node().States[a.selectedState]

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
	a.marginals, err = a.network.Sample(a.evidence, a.samples, a.rng)
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
