package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/francescarpi/mytime/internal/types"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DefaultActivity struct {
	Id         string
	ActivityId int
	Row        int
}

type SyncState struct {
	Tasks                []types.TasksToSync
	Table                *tview.Table
	SelectedIndex        int
	AllTasksHaveActivity bool
	UpdateFooter         func()
	DefaultActivities    []DefaultActivity
	ActionsLock          bool
}

func SyncView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &SyncState{
		Tasks:                deps.Service.GetTasksToSync(),
		SelectedIndex:        0,
		AllTasksHaveActivity: false,
		ActionsLock:          false,
	}

	state.DefaultActivities = make([]DefaultActivity, len(state.Tasks))

	state.Table = tview.NewTable().SetSelectable(true, false)
	state.Table.SetTitle(" Tasks Synchronization ").SetBorder(true)
	state.Table.SetInputCapture(syncInputHandler(app, pages, state, deps))

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
	content += util.Colorize("Close", "Esc", !state.ActionsLock)
	content += util.Colorize("Sync", "s", state.AllTasksHaveActivity && !state.ActionsLock)
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

func syncInputHandler(
	app *tview.Application,
	pages *tview.Pages,
	state *SyncState,
	deps *Dependencies,
) func(event *tcell.EventKey) *tcell.EventKey {
	if state.ActionsLock {
		return nil
	}

	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage("sync")
			pages.AddPage("home", HomeView(app, pages, deps), true, true)
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 's':
				if state.AllTasksHaveActivity {
					handleSyncTasks(app, state, deps)
				}
				return nil
			}
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
			state.DefaultActivities[result.Row-1] = result
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
		Id:         task.Id,
		ActivityId: defaultActivity.Id,
		Row:        row,
	}
}

func handleSyncTasks(app *tview.Application, state *SyncState, deps *Dependencies) {
	log.Println("Syncing tasks...")

	state.ActionsLock = true
	app.QueueUpdateDraw(state.UpdateFooter)

	var wg sync.WaitGroup

	for i, task := range state.Tasks {
		wg.Add(1)
		go syncTask(&wg, app, &task, state.DefaultActivities[i].ActivityId, i+1, state, deps)
	}

	go func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines finished")

		state.ActionsLock = false
		app.QueueUpdateDraw(state.UpdateFooter)
	}()
}

func syncTask(
	wg *sync.WaitGroup,
	app *tview.Application,
	task *types.TasksToSync,
	activityId int,
	row int,
	state *SyncState,
	deps *Dependencies,
) {
	defer wg.Done()

	log.Println("Syncing task:", task.Id, "with activityId:", activityId)
	app.QueueUpdateDraw(func() {
		state.Table.SetCell(row, 6, tview.NewTableCell("[yellow]Syncing..."))
	})

	err := deps.Redmine.SendTask(task.ExternalId, task.Desc, task.Date, task.Duration, activityId)
	if err != nil {
		log.Println("Error syncing task:", err)
		app.QueueUpdateDraw(func() {
			state.Table.SetCell(row, 6, tview.NewTableCell("[red]Error!").SetAlign(tview.AlignCenter))
		})
		return
	}

	app.QueueUpdateDraw(func() {
		state.Table.SetCell(row, 6, tview.NewTableCell("[green]Success!").SetAlign(tview.AlignCenter))
	})

	for _, idStr := range task.Ids.IDs {
		id, _ := strconv.Atoi(idStr)
		deps.Service.SetTaskAsReported(uint(id))
	}
}
