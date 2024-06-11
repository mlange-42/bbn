package tui

import (
	"fmt"

	"github.com/mlange-42/bbn"
	"github.com/rivo/tview"
)

type Table struct {
	tview.TableContentReadOnly
	nodes       []Node
	index       int
	nodesByName map[string]int
	header      []string
}

func NewTable(nodes []Node, index int, nodesByName map[string]int) Table {
	header := append([]string{}, nodes[index].Node().Given...)
	header = append(header, nodes[index].Node().Outcomes...)

	return Table{
		nodes:       nodes,
		index:       index,
		nodesByName: nodesByName,
		header:      header,
	}
}

func (t *Table) GetCell(row, column int) *tview.TableCell {
	if row == 0 {
		cell := tview.NewTableCell(t.header[column])
		cell.SetAlign(tview.AlignCenter)
		cell.SetExpansion(1)
		return cell
	}
	row -= 1

	node := t.nodes[t.index]
	numParents := len(node.Node().Given)

	if column < len(node.Node().Given) {
		stride := 1
		for i := len(node.Node().Given) - 1; i > column; i-- {
			parIdx := t.nodesByName[node.Node().Given[i]]
			par := t.nodes[parIdx]
			stride *= len(par.Node().Outcomes)
		}
		parent := t.nodes[t.nodesByName[node.Node().Given[column]]].Node()
		text := parent.Outcomes[(row/stride)%len(parent.Outcomes)]
		cell := tview.NewTableCell(text)
		cell.SetAlign(tview.AlignRight)
		return cell
	}

	var text string

	values := node.Node().Table[row]

	if node.Node().Type == bbn.UtilityNodeType {
		text = fmt.Sprintf("%9.3f", values[0])
	} else {
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		text = fmt.Sprintf("%7.3f%%", 100*values[column-numParents]/sum)
	}

	cell := tview.NewTableCell(text)
	cell.SetAlign(tview.AlignRight)
	return cell
}

func (t *Table) GetRowCount() int {
	return len(t.nodes[t.index].Node().Table) + 1
}

func (t *Table) GetColumnCount() int {
	return len(t.header)
}
