package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

const REFRESH_RATE = 30

func HomeView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	date := time.Now()

	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.SetTitle(" MyTime ").SetBorder(true)

	body := tview.NewTable().SetBorder(true)
	footer := tview.NewTextView().SetDynamicColors(true).SetBorder(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(body, 0, 1, false).
		AddItem(footer, 4, 0, false)

	render := func() {
		w, err := deps.TaskService.GetWorkedDuration(date)
		if err != nil {
			panic(err)
		}

		header.Clear().
			AddItem(formatHeaderSection("Today", w.DailyFormatted, w.DailyGoalFormatted, w.DailyOvertime), 0, 1, false).
			AddItem(formatHeaderSection("Week", w.WeeklyFormatted, w.WeeklyGoalFormatted, w.WeeklyOvertimeFormatted), 0, 1, false).
			AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).SetText(date.Format("Monday, 2006-01-02")), 0, 1, false)
	}

	render()

	go func() {
		for {
			time.Sleep(REFRESH_RATE * time.Second)
			app.QueueUpdateDraw(render)
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
