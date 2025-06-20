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
	}, nil)

}

func showDeleteTaskModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	task model.Task,
	deps *Dependencies,
) {
	state.Table.SetDisableAutomaticDeselect(true)
	components.ShowConfirmModal(
		app,
		pages,
		"deleteTaskModal",
		fmt.Sprintf("Are you sure you want to delete the task?\n\n%s", task.Desc),
		[]string{"Cancel", "Ok"},
		func(button string) {
			state.Table.SetDisableAutomaticDeselect(false)
			state.Table.Deselect()
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

	externalId := ""
	if task.ExternalId != nil {
		externalId = *task.ExternalId
	}

	form := tview.NewForm().
		AddInputField("Project: ", *task.Project, 0, nil, func(text string) { task.Project = &text }).
		AddInputField("Description", task.Desc, 0, nil, func(text string) { task.Desc = text }).
		AddInputField("External Id", externalId, 0, nil, func(text string) { task.ExternalId = &text }).
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

	state.Table.SetDisableAutomaticDeselect(true)
	components.ShowFormModal("Modify Task", 80, 15, form, pages, app, func() {
		if task.Desc == "" {
			components.ShowAlertModal(app, pages, "Description cannot be empty", nil)
			return
		}

		if task.ExternalId != nil && *task.ExternalId == "" {
			task.ExternalId = nil
		}

		// Call service to update task
		err := deps.Service.UpdateTask(&task)
		if err != nil {
			components.ShowAlertModal(app, pages, fmt.Sprintf("Error updating task: %s", err.Error()), nil)
			return
		}
		state.Render()
	}, func() {
		state.Table.SetDisableAutomaticDeselect(false)
		state.Table.Deselect()
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

	externalId := ""
	if task.ExternalId != nil {
		externalId = *task.ExternalId
	}

	form := tview.NewForm().
		AddTextView("Project: ", *task.Project, 0, 1, false, false).
		AddTextView("External ID: ", externalId, 0, 1, false, false).
		AddInputField("Description: ", "", 0, nil, func(text string) { task.Desc = text })

	state.Table.SetDisableAutomaticDeselect(true)
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
	}, func() {
		state.Table.SetDisableAutomaticDeselect(false)
		state.Table.Deselect()
	})
}

func showSummaryModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	deps *Dependencies,
) {
	summary, err := deps.Service.GetSummaryDuration(state.Date)
	if err != nil {
		return
	}

	form := tview.NewForm().
		AddTextView("Reported: ", summary.Reported, 0, 1, false, false).
		AddTextView("Not Reported: ", summary.NotReported, 0, 1, false, false)
	components.ShowFormModal("Summary", 80, 11, form, pages, app, nil, nil)
}

func showReportConfirmModal(
	app *tview.Application,
	pages *tview.Pages,
	state *HomeState,
	task model.Task,
	deps *Dependencies,
) {
	state.Table.SetDisableAutomaticDeselect(true)
	components.ShowConfirmModal(
		app,
		pages,
		"reportConfirmModal",
		fmt.Sprintf("Are you sure you want to mark as reported the task?\n\n%s", task.Desc),
		[]string{"Cancel", "Ok"},
		func(button string) {
			state.Table.SetDisableAutomaticDeselect(false)
			state.Table.EnableSelection()
			if button == "Ok" {
				err := deps.Service.SetTaskAsReported(task.ID)
				if err != nil {
					components.ShowAlertModal(app, pages, fmt.Sprintf("Error reporting tasks: %s", err.Error()), nil)
					return
				}
				state.Render()
			}
		},
	)
}
