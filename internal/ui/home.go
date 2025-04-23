package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const REFRESH_RATE = 30

type HomeState struct {
	Date          time.Time
	Tasks         []model.Task
	SelectedIndex int
	Table         *tview.Table
	Render        func()
}

func HomeView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &HomeState{
		Date:          time.Now(),
		SelectedIndex: 0,
	}

	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.SetTitle(" MyTime ").SetBorder(true)

	state.Table = tview.NewTable().SetSelectable(true, false)
	state.Table.SetBorder(true)
	// state.Table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow))
	state.Table.SetInputCapture(HomeInputHandler(app, pages, deps, state))

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(state.Table, 0, 1, true).
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

		if state.SelectedIndex >= len(state.Tasks) {
			state.SelectedIndex = 0
		}

		header.Clear().
			AddItem(formatHeaderSection("Today", w.DailyFormatted, w.DailyGoalFormatted, w.DailyOvertime), 0, 1, false).
			AddItem(formatHeaderSection("Week", w.WeeklyFormatted, w.WeeklyGoalFormatted, w.WeeklyOvertimeFormatted), 0, 1, false).
			AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).SetText(state.Date.Format("Monday, 2006-01-02")), 0, 1, false)

		renderFooter(state, footer)
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

func HomeInputHandler(
	app *tview.Application,
	pages *tview.Pages,
	deps *Dependencies,
	state *HomeState,
) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if len(state.Tasks) > 0 {
				task := state.Tasks[state.SelectedIndex]
				err := deps.Service.StartStopTask(task.ID)
				if err != nil {
					log.Printf("Error starting/stopping task: %s", err)
				}
				state.Render()
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				return handleQuit(app)
			case 'h', 'l', 't':
				return handleDateNavigation(event.Rune(), state)
			case 's':
				return handleSyncNavigation(pages)
			case 'j', 'k':
				return handleTaskSelection(event.Rune(), state)
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

func handleSyncNavigation(pages *tview.Pages) *tcell.EventKey {
	pages.SwitchToPage("sync")
	return nil
}

func renderFooter(state *HomeState, footer *tview.TextView) {
	hasTasks := len(state.Tasks) > 0

	// Helper to apply conditional color
	colorize := func(label, key string, enabled bool) string {
		keyColor := "blue"
		if !enabled {
			keyColor = "gray"
		}
		return fmt.Sprintf("[white]%s:[%s] %s [white]| ", label, keyColor, key)
	}

	content := ""
	content += colorize("Quit", "q", true)
	content += colorize("Prev Day", "h", true)
	content += colorize("Next Day", "l", true)
	content += colorize("Today", "t", true)
	content += colorize("Next Task", "j", hasTasks)
	content += colorize("Prev Task", "k", hasTasks)
	content += colorize("Start/Stop", "Enter", hasTasks)
	content += colorize("Duplicate", "d", hasTasks)
	content += colorize("Modify", "m", hasTasks)
	content += colorize("Delete", "x", hasTasks)
	content += colorize("New", "n", hasTasks)
	content += colorize("Sync", "s", true)

	footer.SetText(content)
}
