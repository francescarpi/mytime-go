package ui

import (
	"fmt"

	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RenderTasksTable renders the tasks into the provided table for the current date in state.
func renderTasksTable(state *HomeState) {
	state.Table.Clear()

	// Define header
	headers := []string{"ID", "Project", "Description", "Ext.ID", "Started", "Ended", "Duration", "Reported"}
	for col, h := range headers {
		expanded := 0
		if h == "Description" {
			expanded = 1
		}
		state.Table.SetCell(0, col, tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetExpansion(expanded).SetSelectable(false))
	}

	// Fill rows with tasks
	for row, task := range state.Tasks {
		state.Table.SetCell(row+1, 0, tview.NewTableCell(fmt.Sprintf("%d", task.ID)))
		state.Table.SetCell(row+1, 1, tview.NewTableCell(*task.Project))
		state.Table.SetCell(row+1, 2, tview.NewTableCell(task.Desc).SetExpansion(1))
		state.Table.SetCell(row+1, 3, tview.NewTableCell(*task.ExternalId))
		state.Table.SetCell(row+1, 4, tview.NewTableCell(task.Start.Format("15:04")))
		endFormatted := ""
		if task.End != nil {
			endFormatted = task.End.Format("15:04")
		}
		state.Table.SetCell(row+1, 5, tview.NewTableCell(endFormatted))
		state.Table.SetCell(row+1, 6, tview.NewTableCell(util.HumanizeDuration(task.Duration)).SetAlign(tview.AlignRight))
		state.Table.SetCell(row+1, 7, tview.NewTableCell(task.ReportedIcon()).SetAlign(tview.AlignCenter))
	}

	state.Table.Select(state.SelectedIndex+1, 0) // Select first row (first task)
	state.Table.SetFixed(1, 0)                   // Fix header row
}

func handleTaskSelection(key rune, state *HomeState) *tcell.EventKey {
	if len(state.Tasks) == 0 {
		return nil
	}

	switch key {
	case 'j':
		if len(state.Tasks) > 0 && state.SelectedIndex < len(state.Tasks)-1 {
			state.SelectedIndex++
			state.Table.Select(state.SelectedIndex+1, 0)
		}
		return nil
	case 'k':
		if len(state.Tasks) > 0 && state.SelectedIndex > 0 {
			state.SelectedIndex--
			state.Table.Select(state.SelectedIndex+1, 0)
		}
	}

	return nil
}

func handleTaskCreation() *tcell.EventKey {
	// TODO: logic to create a new task
	return nil
}
