package ui

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/francescarpi/mytime/internal/types"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DefaultActivity struct {
	Id              string
	Ids             []string
	DefaultActivity int
}

type SyncState struct {
	Tasks                []types.TasksToSync
	Table                *tview.Table
	SelectedIndex        int
	AllTasksHaveActivity bool
	UpdateFooter         func()
	DefaultActivities    []DefaultActivity
}

func SyncView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &SyncState{
		Tasks:                deps.Service.GetTasksToSync(),
		SelectedIndex:        0,
		AllTasksHaveActivity: false,
	}

	state.Table = tview.NewTable().SetSelectable(true, false)
	state.Table.SetTitle(" Tasks Synchronization ").SetBorder(true)
	state.Table.SetInputCapture(syncInputHandler(pages))

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(state.Table, 0, 1, true).
		AddItem(footer, 3, 0, false)

	state.UpdateFooter = func() {
		renderSyncFooter(state, footer)
	}

	renderSyncFooter(state, footer)
	renderSyncTable(state)
	loadTasksActivity(app, deps, state)

	return layout
}

func renderSyncFooter(state *SyncState, footer *tview.TextView) {
	content := ""
	content += util.Colorize("Close", "Esc", true)
	content += util.Colorize("Sync", "s", state.AllTasksHaveActivity)
	footer.SetText(content)
}

func renderSyncTable(state *SyncState) {
	state.Table.Clear()

	// Define header
	headers := []string{"Description", "Date", "Duration", "Ext.ID", "Tasks Ids", "Activity", "Status"}
	for col, h := range headers {
		expanded := 0
		if h == "Description" {
			expanded = 1
		}
		state.Table.SetCell(0, col, tview.NewTableCell(fmt.Sprintf("[::b]%s", h)).SetExpansion(expanded).SetSelectable(false))
	}

	// Fill rows with tasks
	for row, task := range state.Tasks {
		state.Table.SetCell(row+1, 0, tview.NewTableCell(task.Desc).SetExpansion(1))
		state.Table.SetCell(row+1, 1, tview.NewTableCell(task.Date))
		state.Table.SetCell(row+1, 2, tview.NewTableCell(util.HumanizeDuration(task.Duration)).SetAlign(tview.AlignRight))
		state.Table.SetCell(row+1, 3, tview.NewTableCell(task.ExternalId))
		state.Table.SetCell(row+1, 4, tview.NewTableCell(strings.Join(task.Ids.IDs, ",")))
		state.Table.SetCell(row+1, 5, tview.NewTableCell("[red]Loading..."))
		state.Table.SetCell(row+1, 6, tview.NewTableCell("[red]✗").SetAlign(tview.AlignCenter))
	}

	state.Table.Select(state.SelectedIndex+1, 0) // Select first row (first task)
	state.Table.SetFixed(1, 0)                   // Fix header row
}

func syncInputHandler(pages *tview.Pages) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.SwitchToPage("home")
			return nil
		}
		return event
	}
}

func loadTasksActivity(app *tview.Application, deps *Dependencies, state *SyncState) {
	var wg sync.WaitGroup
	resultsChan := make(chan DefaultActivity, len(state.Tasks))

	for row, task := range state.Tasks {
		wg.Add(1)
		go loadTaskActivity(&wg, resultsChan, app, deps, state, &task, row+1)
	}

	go func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines finished")
		close(resultsChan)

		for result := range resultsChan {
			log.Println("Result received:", result)
			state.DefaultActivities = append(state.DefaultActivities, result)
		}

		state.AllTasksHaveActivity = true
		app.QueueUpdateDraw(state.UpdateFooter)
	}()
}

func loadTaskActivity(
	wg *sync.WaitGroup,
	resultsChain chan<- DefaultActivity,
	app *tview.Application,
	deps *Dependencies,
	state *SyncState,
	task *types.TasksToSync,
	row int,
) {
	defer wg.Done()

	log.Println("Loading task activity for externalId:", task.ExternalId)
	_, defaultActivity, err := deps.Redmine.LoadActivities(task.ExternalId)
	if err != nil {
		log.Println("Error loading task activity:", err)
		return
	}

	app.QueueUpdateDraw(func() {
		state.Table.SetCell(row, 5, tview.NewTableCell(fmt.Sprintf("[green]%s", defaultActivity.Name)))
	})

	resultsChain <- DefaultActivity{
		Id:              task.Id,
		Ids:             task.Ids.IDs,
		DefaultActivity: defaultActivity.Id,
	}
}
