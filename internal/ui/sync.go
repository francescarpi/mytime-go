package ui

import (
	"fmt"
	"strings"

	"github.com/francescarpi/mytime/internal/types"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SyncState struct {
	Tasks         []types.TasksToSync
	Table         *tview.Table
	SelectedIndex int
}

func SyncView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &SyncState{
		Tasks:         deps.Service.GetTasksToSync(),
		SelectedIndex: 0,
	}

	state.Table = tview.NewTable().SetSelectable(true, false)
	state.Table.SetTitle(" Tasks Synchronization ").SetBorder(true)
	state.Table.SetInputCapture(syncInputHandler(pages))

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(state.Table, 0, 1, true).
		AddItem(footer, 3, 0, false)

	renderSyncFooter(state, footer)
	renderSyncTable(state)

	return layout
}

func renderSyncFooter(state *SyncState, footer *tview.TextView) {
	// hasTasks := len(state.Tasks) > 0

	content := ""
	content += util.Colorize("Close", "Esc", true)
	// content += util.Colorize("Sync", "s", true)

	footer.SetText(content)
}

func renderSyncTable(state *SyncState) {
	state.Table.Clear()

	// Define header
	headers := []string{"Description", "Date", "Duration", "Ext.ID", "Tasks Ids", "Activity", "Status"}
	for col, h := range headers {
		expanded := 0
		if h == "Description" {
			expanded = 1
		}
		state.Table.SetCell(0, col, tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetExpansion(expanded).SetSelectable(false))
	}

	// Fill rows with tasks
	for row, task := range state.Tasks {
		state.Table.SetCell(row+1, 0, tview.NewTableCell(task.Desc).SetExpansion(1))
		state.Table.SetCell(row+1, 1, tview.NewTableCell(task.Date))
		state.Table.SetCell(row+1, 2, tview.NewTableCell(util.HumanizeDuration(task.Duration)).SetAlign(tview.AlignRight))
		state.Table.SetCell(row+1, 3, tview.NewTableCell(task.ExternalId))
		state.Table.SetCell(row+1, 4, tview.NewTableCell(strings.Join(task.Ids.IDs, ",")))
		state.Table.SetCell(row+1, 5, tview.NewTableCell("[red]Loading..."))
		state.Table.SetCell(row+1, 6, tview.NewTableCell("[red]✗").SetAlign(tview.AlignCenter))
	}

	state.Table.Select(state.SelectedIndex+1, 0) // Select first row (first task)
	state.Table.SetFixed(1, 0)                   // Fix header row
}

func syncInputHandler(pages *tview.Pages) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.SwitchToPage("home")
			return nil
		}
		return event
	}
}
