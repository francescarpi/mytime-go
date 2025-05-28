package components

import (
	"context"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DESELECT_TIMEOUT = 2 * time.Second

type Table struct {
	app                      *tview.Application
	table                    *tview.Table
	userInputCapture         func(event *tcell.EventKey) *tcell.EventKey
	cancelFunc               context.CancelFunc
	callback                 func()
	disableAutomaticDeselect bool
}

func GetNewTable(app *tview.Application, header []string, callback func()) *Table {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true)
	table.SetFixed(1, 0)

	for col, h := range header {
		expanded := 0
		if h == "Description" {
			expanded = 1
		}
		table.SetCell(0, col, tview.NewTableCell("[yellow]"+h).SetSelectable(false).SetExpansion(expanded))
	}

	customTable := &Table{
		table:    table,
		app:      app,
		callback: callback,
	}

	table.SetInputCapture(customTable.masterInputCapture)
	customTable.Deselect()

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
			t.EnableSelection()
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
	if t == nil {
		return -1
	}

	totalRows := t.table.GetRowCount()
	row, _ := t.table.GetSelection()

	if row == 0 || row >= totalRows {
		return -1
	}

	return row - 1
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

	go func(t *Table) {
		t.callback()
	}(t)
}

func (t *Table) EnableSelection() {
	if t.cancelFunc != nil {
		t.cancelFunc()
	}

	var ctx context.Context
	ctx, t.cancelFunc = context.WithCancel(context.Background())

	rowSelectable, _ := t.table.GetSelectable()
	if !rowSelectable {
		t.table.SetSelectable(true, false)
		t.table.Select(0, 0)

		go func(t *Table) {
			t.callback()
		}(t)
	}

	go func(ctx context.Context, app *tview.Application) {
		log.Println("Start goroutine to deselect")

		select {
		case <-time.After(DESELECT_TIMEOUT):
			if !t.disableAutomaticDeselect {
				app.QueueUpdateDraw(t.Deselect)
			}
			log.Println("Stop goroutine to deselect")
		case <-ctx.Done():
			log.Println("Goroutine to deselect cancelled")
		}

	}(ctx, t.app)
}

func (t *Table) SetDisableAutomaticDeselect(isThereAnOpenedModal bool) {
	t.disableAutomaticDeselect = isThereAnOpenedModal
}
