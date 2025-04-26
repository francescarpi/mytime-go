package ui

import (
	"fmt"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/ui/components"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/rivo/tview"
)

const REFRESH_RATE = 10

type HomeState struct {
	Date               time.Time
	Tasks              []model.Task
	Table              *components.Table
	Render             func()
	RenderAndGotoToday func()
	ActionsManager     *ActionsManager
}

func HomeView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &HomeState{
		Date: time.Now(),
	}

	state.RenderAndGotoToday = func() {
		state.Date = time.Now()
		state.Render()
	}

	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.SetTitle(" MyTime ").SetBorder(true)

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	state.Table = components.GetNewTable([]string{"ID", "Project", "Description", "Ext.ID", "Started", "Ended", "Duration", "Reported"})

	state.ActionsManager = GetNewActionsManager(footer, homeViewActions(app, pages, deps, state))
	state.Table.SetInputCapture(state.ActionsManager.GetInputHandler())

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(state.Table.GetTable(), 0, 1, true).
		AddItem(footer, 4, 0, false)

	state.Render = func() {
		w, err := deps.Service.GetWorkedDuration(state.Date)
		if err != nil {
			panic(err)
		}

		tasks, err := deps.Service.GetTasksByDate(state.Date)
		if err != nil {
			panic(err)
		}
		state.Tasks = tasks

		header.Clear().
			AddItem(formatHeaderSection("Today", w.DailyFormatted, w.DailyGoalFormatted, w.DailyOvertime), 0, 1, false).
			AddItem(formatHeaderSection("Week", w.WeeklyFormatted, w.WeeklyGoalFormatted, w.WeeklyOvertimeFormatted), 0, 1, false).
			AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).SetText(state.Date.Format("Monday, 2006-01-02")), 0, 1, false)

		renderTasksTable(state)
		state.ActionsManager.Refresh()
	}

	state.Render()

	go func() {
		for {
			time.Sleep(REFRESH_RATE * time.Second)
			app.QueueUpdateDraw(state.Render)
		}
	}()

	return layout
}

func formatHeaderSection(title, formatted, goal, overtime string) *tview.TextView {
	text := ""
	if overtime == "" {
		text = fmt.Sprintf("[red]%s: %s of %s[-]", title, formatted, goal)
	} else {
		text = fmt.Sprintf("[green]%s: %s of %s (+%s)[-]", title, formatted, goal, overtime)
	}

	return tview.NewTextView().
		SetDynamicColors(true).
		SetText(text)
}

func getSelectedTask(state *HomeState) (model.Task, error) {
	if len(state.Tasks) == 0 {
		return model.Task{}, fmt.Errorf("no tasks available")
	}

	row := state.Table.GetRowSelected()
	if row == 0 {
		return model.Task{}, fmt.Errorf("no task selected")
	}
	return state.Tasks[row-1], nil
}

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

	state.Table.Deselect()
}
