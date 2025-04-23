package ui

import (
	"fmt"
	"time"

	"github.com/francescarpi/mytime/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const REFRESH_RATE = 30

type HomeState struct {
	Date   time.Time
	Tasks  []model.Task
	Render func()
}

func HomeView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &HomeState{
		Date: time.Now(),
	}

	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.SetTitle(" MyTime ").SetBorder(true)
	body := tview.NewTable().SetBorder(true)
	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	state.Render = func() {
		w, err := deps.TaskService.GetWorkedDuration(state.Date)
		if err != nil {
			panic(err)
		}

		header.Clear().
			AddItem(formatHeaderSection("Today", w.DailyFormatted, w.DailyGoalFormatted, w.DailyOvertime), 0, 1, false).
			AddItem(formatHeaderSection("Week", w.WeeklyFormatted, w.WeeklyGoalFormatted, w.WeeklyOvertimeFormatted), 0, 1, false).
			AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).SetText(state.Date.Format("Monday, 2006-01-02")), 0, 1, false)

		renderFooter(state, footer)
	}

	body.SetInputCapture(HomeInputHandler(app, pages, deps, state))

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(body, 0, 1, true).
		AddItem(footer, 4, 0, false)

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
		text = fmt.Sprintf("[red]%s: %s/%s[-]", title, formatted, goal)
	} else {
		text = fmt.Sprintf("[green]%s: %s/%s (+%s)[-]", title, formatted, goal, overtime)
	}

	return tview.NewTextView().
		SetDynamicColors(true).
		SetText(text)
}

func HomeInputHandler(app *tview.Application, pages *tview.Pages, deps *Dependencies, state *HomeState) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
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
			case 'n':
				return handleTaskCreation()
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

func handleTaskSelection(key rune, state *HomeState) *tcell.EventKey {
	if len(state.Tasks) == 0 {
		return nil
	}

	switch key {
	case 'j':
		// TODO: next selection
	case 'k':
		// TODO: previous selection
	}

	state.Render()
	return nil
}

func handleTaskCreation() *tcell.EventKey {
	// TODO: logic to create a new task
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
