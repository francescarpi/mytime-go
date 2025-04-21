package gui

import (
	"log"
	"time"

	"mytime/tasks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gorm.io/gorm"
)

var (
	app        *tview.Application
	pages      *tview.Pages
	date       time.Time
	tasksTable *TasksTable
	header     *Header
)

const REFRESH_RATE = 30

func fullRefresh(selectTask bool) {
	log.Println("Full refresh")
	tasksTable.Refresh(selectTask)
	header.Refresh()
}

func updateContent() {
	for {
		time.Sleep(REFRESH_RATE * time.Second)
		app.QueueUpdateDraw(func() {
			fullRefresh(false)
		})
	}
}

func GotoToday() {
	log.Println("Goto today")
	date = time.Now()
}

func Init(conn *gorm.DB) {
	log.Println("========== GUI Initiates ==========")

	app = tview.NewApplication()
	pages = tview.NewPages()
	date = time.Now()

	tasksManager := tasks.TasksManager{Conn: conn}
	header = GetNewHeader(&tasksManager)
	footer := GetFooter()
	tasksTable = GetNewTasksTable(&tasksManager)

	layout := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header.container, 3, 0, false).
			AddItem(tasksTable.container, 0, 1, false).
			AddItem(footer, 4, 0, false), 0, 1, false)

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
				if success := tasksTable.SelectNextTask(); success {
					tasksTable.Refresh(false)
				}
			case 'k':
				if success := tasksTable.SelectPreviousTask(); success {
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
			}
		}
		return event
	})

	fullRefresh(true)
	go updateContent()

	pages.AddPage("main", layout, true, true)

	if err := app.SetRoot(pages, true).SetFocus(layout).Run(); err != nil {
		panic(err)
	}

}
