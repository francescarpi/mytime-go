package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/francescarpi/mytime/internal/ui/components"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const REFRESH_RATE = 30

type HomeState struct {
	Date               time.Time
	Tasks              []model.Task
	Table              *components.Table
	Render             func()
	RenderAndGotoToday func()
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

	state.Table = components.GetNewTable([]string{"ID", "Project", "Description", "Ext.ID", "Started", "Ended", "Duration", "Reported"})
	state.Table.SetInputCapture(homeInputHandler(app, pages, deps, state))

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

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

		renderHomeFooter(deps, state, footer)
		renderTasksTable(state)
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

func getSelectedTask(state *HomeState) model.Task {
	row := state.Table.GetRowSelected()
	return state.Tasks[row-1]
}

func homeInputHandler(
	app *tview.Application,
	pages *tview.Pages,
	deps *Dependencies,
	state *HomeState,
) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if len(state.Tasks) > 0 {
				task := getSelectedTask(state)
				err := deps.Service.StartStopTask(task.ID)
				if err != nil {
					log.Printf("Error starting/stopping task: %s", err)
				}
				state.RenderAndGotoToday()
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				return handleQuit(app)
			case 'h', 'l', 't':
				return handleDateNavigation(event.Rune(), state)
			case 's':
				return handleSyncNavigation(app, pages, deps)
			case 'd', 'm', 'x', 'n':
				return handleTaskManipulation(event.Rune(), state, pages, app, deps)
			}
		}
		return event
	}
}

func handleQuit(app *tview.Application) *tcell.EventKey {
	app.Stop()
	fmt.Println("Bye!")
	return nil
}

func handleDateNavigation(key rune, state *HomeState) *tcell.EventKey {
	switch key {
	case 'h':
		state.Date = state.Date.AddDate(0, 0, -1)
	case 'l':
		next := state.Date.AddDate(0, 0, 1)
		if next.After(time.Now()) {
			return nil
		}
		state.Date = next
	case 't':
		state.Date = time.Now()
	}
	state.Render()
	return nil
}

func handleSyncNavigation(app *tview.Application, pages *tview.Pages, deps *Dependencies) *tcell.EventKey {
	tasksToSync := deps.Service.GetTasksToSync()
	if len(tasksToSync) == 0 {
		return nil
	}
	pages.RemovePage("home")
	pages.AddPage("sync", SyncView(app, pages, deps), true, true)
	return nil
}

func renderHomeFooter(deps *Dependencies, state *HomeState, footer *tview.TextView) {
	hasTasks := len(state.Tasks) > 0

	tasksToSync := deps.Service.GetTasksToSync()
	tasksToSyncCount := len(tasksToSync)

	content := ""
	content += util.Colorize("Quit", "q", true)
	content += util.Colorize("Prev Day", "h", true)
	content += util.Colorize("Next Day", "l", true)
	content += util.Colorize("Today", "t", true)
	content += util.Colorize("Next Task", "j", hasTasks)
	content += util.Colorize("Prev Task", "k", hasTasks)
	content += util.Colorize("Start/Stop", "Enter", hasTasks)
	content += util.Colorize("Duplicate", "d", hasTasks)
	content += util.Colorize("Modify", "m", hasTasks)
	content += util.Colorize("Delete", "x", hasTasks)
	content += util.Colorize("New", "n", hasTasks)
	content += util.Colorize(fmt.Sprintf("Sync (%d)", tasksToSyncCount), "s", tasksToSyncCount > 0)

	footer.SetText(content)
}
