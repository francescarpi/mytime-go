package ui

import (
	"fmt"
	"log"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/ui/components"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/rivo/tview"
)

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
