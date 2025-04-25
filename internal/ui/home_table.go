package ui

import (
	"fmt"
	"log"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/ui/components"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RenderTasksTable renders the tasks into the provided table for the current date in state.
func renderTasksTable(state *HomeState) {
	renderer := state.Table.GetRowRenderer()
	for row, task := range state.Tasks {
		row := row + 1
		renderer(row, 0, fmt.Sprintf("%d", task.ID), 0, tview.AlignLeft)
		renderer(row, 1, *task.Project, 0, tview.AlignLeft)
		renderer(row, 2, task.Desc, 1, tview.AlignLeft)
		renderer(row, 3, *task.ExternalId, 0, tview.AlignLeft)
		renderer(row, 4, task.Start.Format("15:04"), 0, tview.AlignCenter)

		endFormatted := "🚗"
		if task.End != nil {
			endFormatted = task.End.Format("15:04")
		}

		renderer(row, 5, endFormatted, 0, tview.AlignCenter)
		renderer(row, 6, util.HumanizeDuration(task.Duration), 0, tview.AlignRight)
		renderer(row, 7, task.ReportedIcon(), 0, tview.AlignCenter)
	}

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

	components.ShowFormModal("New Task", 80, 11, form, pages, app, func() {
		if description == "" {
			components.ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		// Call service to create task
		err := deps.Service.CreateTask(description, project, externalId)
		if err != nil {
			components.ShowAlertModal(app, pages, fmt.Sprintf("Error creating task: %s", err.Error()), nil)
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
	components.ShowConfirmModal(
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
					components.ShowAlertModal(app, pages, fmt.Sprintf("Error deleting task: %s", err.Error()), nil)
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

	components.ShowFormModal("Modify Task", 80, 15, form, pages, app, func() {
		if task.Desc == "" {
			components.ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		// Call service to update task
		err := deps.Service.UpdateTask(&task)
		if err != nil {
			components.ShowAlertModal(app, pages, fmt.Sprintf("Error updating task: %s", err.Error()), nil)
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

	components.ShowFormModal("Duplicate Task", 80, 11, form, pages, app, func() {
		if task.Desc == "" {
			components.ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		// Call service to create duplicated task
		err := deps.Service.CreateTask(task.Desc, task.Project, task.ExternalId)
		if err != nil {
			components.ShowAlertModal(app, pages, fmt.Sprintf("Error creating task: %s", err.Error()), nil)
			return
		}
		state.RenderAndGotoToday()
	})
}
