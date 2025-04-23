package ui

import (
	"fmt"

	"github.com/francescarpi/mytime/internal/model"
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

func handleTaskManipulation(key rune, state *HomeState, pages *tview.Pages, app *tview.Application, deps *Dependencies) *tcell.EventKey {
	if len(state.Tasks) == 0 {
		return nil
	}

	switch key {
	case 'd':
		if len(state.Tasks) > 0 {
			taskToDuplicate := state.Tasks[state.SelectedIndex]
			showDuplicateTaskModal(app, pages, state, taskToDuplicate, deps)
		}
		return nil
	}

	return nil
}

func showDuplicateTaskModal(app *tview.Application, pages *tview.Pages, state *HomeState, task model.Task, deps *Dependencies) {
	task.Desc = ""

	form := tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetFieldBackgroundColor(tcell.ColorGray).
		AddTextView("Project: ", *task.Project, 0, 1, false, false).
		AddTextView("External ID: ", *task.ExternalId, 0, 1, false, false).
		AddInputField("Description: ", "", 0, nil, func(text string) { task.Desc = text }).
		AddButton("Ok", func() {
			if task.Desc == "" {
				ShowAlertModal(app, pages, "Description cannot be empty", nil)
				return
			}

			// Call service to create duplicated task
			err := deps.Service.CreateTask(task.Desc, task.Project, task.ExternalId)
			if err != nil {
				ShowAlertModal(app, pages, fmt.Sprintf("Error creating task: %s", err.Error()), nil)
				return
			}
			state.Render()
			pages.RemovePage("duplicateTaskModal")
		}).
		AddButton("Cancel", func() {
			pages.RemovePage("duplicateTaskModal")
		})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.RemovePage("duplicateTaskModal")
			return nil
		}
		return event
	})

	layout := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	layout.SetTitle("Duplicate Task").SetBorder(true)

	pages.AddPage("duplicateTaskModal", ModalPrimitive(layout, 80, 11), true, true)

	// Show the modal
	app.SetFocus(form)
}
