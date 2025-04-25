package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	table *tview.Table
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

	return &Table{
		table,
	}
}

func (t *Table) SetTitle(title string) {
	t.table.SetTitle(" " + title + " ")
}

func (t *Table) SetInputCapture(fn func(event *tcell.EventKey) *tcell.EventKey) {
	t.table.SetInputCapture(fn)
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
