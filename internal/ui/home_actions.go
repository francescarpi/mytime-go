package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func homeViewActions(app *tview.Application, pages *tview.Pages, deps *Dependencies, state *HomeState) *[]Action {
	quitAction := GetNewAction("Quit", NewRuneKey("q", 'q'),
		func() bool { return true },
		func() {
			app.Stop()
			fmt.Println("Bye!")
			log.Println("\n----------- Bye! ----------")
		},
	)

	prevDay := GetNewAction("Prev Day", NewRuneKey("h", 'h'),
		func() bool { return true },
		func() {
			state.Date = state.Date.AddDate(0, 0, -1)
			state.Render()
		},
	)

	nextDay := GetNewAction("Next Day", NewRuneKey("l", 'l'),
		func() bool { return true },
		func() {
			next := state.Date.AddDate(0, 0, 1)
			if next.After(time.Now()) {
				return
			}
			state.Date = next
			state.Render()
		},
	)

	today := GetNewAction("Today", NewRuneKey("t", 't'),
		func() bool { return true },
		func() {
			state.Date = time.Now()
			state.Render()
		},
	)

	nextTask := GetNewAction("Next Task", NewRuneKey("j", 'j'),
		func() bool { return len(state.Tasks) > 0 },
		func() {},
	)

	prevTask := GetNewAction("Prev Task", NewRuneKey("k", 'k'),
		func() bool { return len(state.Tasks) > 0 },
		func() {},
	)

	startStop := GetNewAction("Start/Stop", NewSpecialKey("Enter", tcell.KeyEnter),
		func() bool {
			_, err := getSelectedTask(state)
			return err == nil
		},
		func() {
			task, _ := getSelectedTask(state)
			err := deps.Service.StartStopTask(task.ID)
			if err != nil {
				log.Printf("Error starting/stopping task: %s", err)
			}
			state.RenderAndGotoToday()
		},
	)

	duplicate := GetNewAction("Duplicate", NewRuneKey("d", 'd'),
		func() bool {
			_, err := getSelectedTask(state)
			return err == nil
		},
		func() {
			task, _ := getSelectedTask(state)
			showDuplicateTaskModal(app, pages, state, task, deps)
		},
	)

	modify := GetNewAction("Modify", NewRuneKey("m", 'm'),
		func() bool {
			_, err := getSelectedTask(state)
			return err == nil
		},
		func() {
			task, _ := getSelectedTask(state)
			showModifyTaskModal(app, pages, state, task, deps)
		},
	)

	deleteAction := GetNewAction("Delete", NewRuneKey("x", 'x'),
		func() bool {
			_, err := getSelectedTask(state)
			return err == nil
		},
		func() {
			task, _ := getSelectedTask(state)
			showDeleteTaskModal(app, pages, state, task, deps)
		},
	)

	newAction := GetNewAction("New", NewRuneKey("n", 'n'),
		func() bool { return true },
		func() {
			showNewTaskModal(app, pages, state, deps)
		},
	)

	syncView := GetNewAction("Sync", NewRuneKey("s", 's'),
		func() bool {
			tasksToSync := len(deps.Service.GetTasksToSync())
			return tasksToSync > 0
		},
		func() {
			pages.
				RemovePage("home").
				AddPage("sync", SyncView(app, pages, deps), true, true)
		},
	)

	summaryAction := GetNewAction("Summary", NewRuneKey("y", 'y'),
		func() bool { return true },
		func() {
			showSummaryModal(app, pages, state, deps)
		},
	)

	markAsReport := GetNewAction("Report", NewRuneKey("r", 'r'),
		func() bool {
			task, err := getSelectedTask(state)
			return err == nil && !task.Reported
		},
		func() {
			task, _ := getSelectedTask(state)
			showReportConfirmModal(app, pages, state, task, deps)
		},
	)

	return &[]Action{
		quitAction,
		prevDay,
		nextDay,
		today,
		nextTask,
		prevTask,
		newAction,
		startStop,
		duplicate,
		modify,
		deleteAction,
		syncView,
		summaryAction,
		markAsReport,
	}
}
