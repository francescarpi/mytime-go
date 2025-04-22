package gui

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Action struct {
	Name     string
	Key      string
	Disabled func() bool
	Format   func() string
}

func GotoPreviousDate() {
	log.Println("Goto previous date. Current is ", date.Format(time.DateOnly))
	date = date.AddDate(0, 0, -1)
}

func GotoNextDate() bool {
	log.Println("Goto next date. Current is ", date.Format(time.DateOnly))
	nextDate := date.AddDate(0, 0, 1)
	if nextDate.After(time.Now()) {
		log.Println("Cannot go to future date")
		return false
	}
	date = nextDate
	return true
}

func SetLayoutInputCapture(layout *tview.Flex) {
	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			tasksTable.StartStopTask()
			GotoToday()
			fullRefresh(true)
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				app.Stop()
				log.Println("---------- Bye! ----------")
			case 'h':
				GotoPreviousDate()
				fullRefresh(true)
			case 'l':
				if success := GotoNextDate(); success {
					fullRefresh(true)
				}
			case 'j':
				if success := tasksTable.SelectNextPrevious(GO_NEXT); success {
					tasksTable.Refresh(false)
				}
			case 'k':
				if success := tasksTable.SelectNextPrevious(GO_PREVIOUS); success {
					tasksTable.Refresh(false)
				}
			case 't':
				GotoToday()
				fullRefresh(true)
			case 'd':
				tasksTable.DuplicateWithDescription()
			case 'm':
				tasksTable.Modify()
			case 'x':
				tasksTable.Delete()
			case 'n':
				tasksTable.New()
			case 's':
				tasksTable.Sync()
			}
		}
		return event
	})
}

func RenderActions(actions *[]Action) string {
	text := ""
	for _, action := range *actions {
		actionName := action.Name
		if action.Format != nil {
			actionName = action.Format()
		}

		disabled := false
		if action.Disabled != nil {
			disabled = action.Disabled()
		}

		if disabled {
			text += fmt.Sprintf("[gray]%s: %s[white] | ", actionName, action.Key)
		} else {
			text += fmt.Sprintf("[darkgray]%s[white]: [blue]%s[white] | ", actionName, action.Key)
		}
	}
	return text
}
