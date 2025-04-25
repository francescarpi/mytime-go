package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	table            *tview.Table
	userInputCapture func(event *tcell.EventKey) *tcell.EventKey
}

func GetNewTable(header []string) *Table {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true)
	table.SetFixed(1, 0)
	table.SetSelectedStyle(
		tcell.Style{}.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorBlue),
	)

	for col, h := range header {
		expanded := 0
		if h == "Description" {
			expanded = 1
		}
		table.SetCell(0, col, tview.NewTableCell("[yellow]"+h).SetSelectable(false).SetExpansion(expanded))
	}

	customTable := &Table{
		table: table,
	}

	table.SetInputCapture(customTable.masterInputCapture)

	return customTable
}

func (t *Table) SetTitle(title string) {
	t.table.SetTitle(" " + title + " ")
}

func (t *Table) masterInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if t.userInputCapture != nil {
		t.userInputCapture(event)
	}

	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'j', 'k':
			rowSelectable, _ := t.table.GetSelectable()
			if !rowSelectable {
				t.table.SetSelectable(true, false)
				t.table.Select(0, 0)
			}
			return event
		}
	}
	return event
}

func (t *Table) SetInputCapture(fn func(event *tcell.EventKey) *tcell.EventKey) {
	t.userInputCapture = fn
}

func (t *Table) GetTable() *tview.Table {
	return t.table
}

func (t *Table) GetRowSelected() int {
	row, _ := t.table.GetSelection()
	return row
}

func (t *Table) SetCellText(row, col int, text string) {
	t.table.GetCell(row, col).SetText(text)
}

func (t *Table) clearRows() {
	rowsToRemove := t.table.GetRowCount() - 1
	for i := rowsToRemove; i > 0; i-- {
		t.table.RemoveRow(i)
	}
}

func (t *Table) GetRowRenderer() func(int, int, string, int, int) {
	t.clearRows()
	return func(row, col int, content string, expansion, align int) {
		cell := tview.NewTableCell(content).SetExpansion(expansion).SetAlign(align)
		t.table.SetCell(row, col, cell)
	}
}

func (t *Table) Deselect() {
	t.table.Select(0, 0)
	t.table.SetSelectable(false, false)
}
