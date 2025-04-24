package ui

import (
	"fmt"
	"log"

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

		endFormatted := "🚗"
		if task.End != nil {
			endFormatted = task.End.Format("15:04")
		}
		state.Table.SetCell(row+1, 5, tview.NewTableCell(endFormatted).SetAlign(tview.AlignCenter))

		state.Table.SetCell(row+1, 6, tview.NewTableCell(util.HumanizeDuration(task.Duration)).SetAlign(tview.AlignRight))
		state.Table.SetCell(row+1, 7, tview.NewTableCell(task.ReportedIcon()).SetAlign(tview.AlignCenter))
	}

	state.Table.SetFixed(1, 0) // Fix header row
}

func handleTaskManipulation(
	key rune,
	state *HomeState,
	pages *tview.Pages,
	app *tview.Application,
	deps *Dependencies,
) *tcell.EventKey {
	if len(state.Tasks) == 0 {
		return nil
	}

	task := getSelectedTask(state)

	switch key {
	case 'd':
		showDuplicateTaskModal(app, pages, state, task, deps)
		return nil
	case 'm':
		showModifyTaskModal(app, pages, state, task, deps)
		return nil
	case 'x':
		showDeleteTaskModal(app, pages, state, task, deps)
		return nil
	case 'n':
		showNewTaskModal(app, pages, state, deps)
		return nil

	}

	return nil
}

func showNewTaskModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	deps *Dependencies,
) {
	var description string
	var project, externalId *string

	form := tview.NewForm().
		AddInputField("Project: ", "", 0, nil, func(text string) { project = &text }).
		AddInputField("Description", "", 0, nil, func(text string) { description = text }).
		AddInputField("External Id", "", 0, nil, func(text string) { externalId = &text })

	ShowFormModal("New Task", 80, 11, form, pages, app, func() {
		if description == "" {
			ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		// Call service to create task
		err := deps.Service.CreateTask(description, project, externalId)
		if err != nil {
			ShowAlertModal(app, pages, fmt.Sprintf("Error creating task: %s", err.Error()), nil)
			return
		}
		state.RenderAndGotoToday()
	})

}

func showDeleteTaskModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	task model.Task,
	deps *Dependencies,
) {
	ShowConfirmModal(
		app,
		pages,
		"deleteTaskModal",
		fmt.Sprintf("Are you sure you want to delete the task?\n\n%s", task.Desc),
		[]string{"Cancel", "Ok"},
		func(button string) {
			if button == "Ok" {
				log.Printf("Deleting task %d", task.ID)
				err := deps.Service.DeleteTask(task.ID)
				if err != nil {
					ShowAlertModal(app, pages, fmt.Sprintf("Error deleting task: %s", err.Error()), nil)
					return
				}
				state.Render()
			}
		},
	)
}

func showModifyTaskModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	task model.Task,
	deps *Dependencies,
) {
	defaultEnd := ""
	if task.End != nil {
		defaultEnd = task.End.Format("15:04")
	}
	form := tview.NewForm().
		AddInputField("Project: ", *task.Project, 0, nil, func(text string) { task.Project = &text }).
		AddInputField("Description", task.Desc, 0, nil, func(text string) { task.Desc = text }).
		AddInputField("External Id", *task.ExternalId, 0, nil, func(text string) { task.ExternalId = &text }).
		AddInputField("Started", task.Start.Format("15:04"), 0, nil, func(text string) {
			newTime, err := util.UpdateTime(&task.Start.Time, text)
			if err != nil {
				return
			}
			task.Start = model.LocalTimestamp{Time: newTime}
		}).
		AddInputField("Ended", defaultEnd, 0, nil, func(text string) {
			newTime, err := util.UpdateTime(&task.Start.Time, text)
			if err != nil {
				return
			}
			task.End = &model.LocalTimestamp{Time: newTime}
		})

	ShowFormModal("Modify Task", 80, 15, form, pages, app, func() {
		if task.Desc == "" {
			ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		// Call service to update task
		err := deps.Service.UpdateTask(&task)
		if err != nil {
			ShowAlertModal(app, pages, fmt.Sprintf("Error updating task: %s", err.Error()), nil)
			return
		}
		state.Render()
	})
}

func showDuplicateTaskModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	task model.Task,
	deps *Dependencies,
) {
	task.Desc = ""

	form := tview.NewForm().
		AddTextView("Project: ", *task.Project, 0, 1, false, false).
		AddTextView("External ID: ", *task.ExternalId, 0, 1, false, false).
		AddInputField("Description: ", "", 0, nil, func(text string) { task.Desc = text })

	ShowFormModal("Duplicate Task", 80, 11, form, pages, app, func() {
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
		state.RenderAndGotoToday()
	})
}
