package gui

import (
	"log"
	"time"

	"mytime/tasks"

	"github.com/rivo/tview"
	"gorm.io/gorm"
)

var (
	app        *tview.Application
	pages      *tview.Pages
	date       time.Time
	tasksTable *TasksTable
	header     *Header
	footer     *Footer
)

const REFRESH_RATE = 30

func fullRefresh(selectTask bool) {
	log.Println("Full refresh")
	tasksTable.Refresh(selectTask)
	header.Refresh()
	footer.Refresh()
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
	tasksTable = GetNewTasksTable(&tasksManager)
	header = GetNewHeader(&tasksManager)
	footer = GetNewFooter(&tasksManager)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header.container, 3, 0, false).
		AddItem(tasksTable.container, 0, 1, false).
		AddItem(footer.container, 4, 0, false)
	pages.AddPage("main", layout, true, true)

	SetLayoutInputCapture(layout)
	fullRefresh(true)

	go updateContent()

	if err := app.SetRoot(pages, true).SetFocus(layout).Run(); err != nil {
		panic(err)
	}

}
