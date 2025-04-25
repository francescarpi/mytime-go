package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/francescarpi/mytime/internal/service/redmine"
	"github.com/francescarpi/mytime/internal/types"
	"github.com/francescarpi/mytime/internal/ui/components"
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TaskToSyncActivities struct {
	Activities *[]redmine.RedmineProjectActivity
	Default    *redmine.RedmineProjectActivity
	Index      int
}

type SyncState struct {
	Tasks                []types.TasksToSync
	Table                *components.Table
	AllTasksHaveActivity bool
	ActionsLock          bool
	TasksActivities      []TaskToSyncActivities
	ActionsManager       *ActionsManager
}

func SyncView(app *tview.Application, pages *tview.Pages, deps *Dependencies) tview.Primitive {
	state := &SyncState{
		Tasks:                deps.Service.GetTasksToSync(),
		AllTasksHaveActivity: false,
		ActionsLock:          true,
	}

	state.TasksActivities = make([]TaskToSyncActivities, len(state.Tasks))

	state.Table = components.GetNewTable([]string{"Description", "Date", "Duration", "Ext.ID", "Tasks Ids", "Activity", "Status"})
	state.Table.SetTitle("Tasks Synchronization")

	footer := tview.NewTextView()
	footer.SetDynamicColors(true).SetBorder(true)

	state.ActionsManager = GetNewActionsManager(footer, syncViewActions(app, pages, deps, state))
	state.Table.SetInputCapture(state.ActionsManager.GetInputHandler())

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(state.Table.GetTable(), 0, 1, true).
		AddItem(footer, 3, 0, false)

	renderSyncTable(state)
	loadTasksActivity(app, deps, state)

	return layout
}

func renderSyncTable(state *SyncState) {
	renderer := state.Table.GetRowRenderer()

	for row, task := range state.Tasks {
		row := row + 1
		renderer(row, 0, task.Desc, 1, tview.AlignLeft)
		renderer(row, 1, task.Date, 0, tview.AlignLeft)
		renderer(row, 2, util.HumanizeDuration(task.Duration), 0, tview.AlignRight)
		renderer(row, 3, task.ExternalId, 0, tview.AlignLeft)
		renderer(row, 4, strings.Join(task.Ids.IDs, ","), 0, tview.AlignLeft)
		renderer(row, 5, "[red]Loading...", 0, tview.AlignLeft)
		renderer(row, 6, "[red]✗", 0, tview.AlignCenter)
	}

	state.Table.Deselect()

}

func syncViewActions(app *tview.Application, pages *tview.Pages, deps *Dependencies, state *SyncState) *[]Action {
	closeAction := GetNewAction("Close", NewSpecialKey("Esc", tcell.KeyEsc),
		func() bool {
			return !state.ActionsLock
		},
		func() {
			pages.RemovePage("sync")
			pages.AddPage("home", HomeView(app, pages, deps), true, true)
		},
	)

	syncAction := GetNewAction("Sync", NewRuneKey("s", 's'),
		func() bool {
			return !state.ActionsLock && state.AllTasksHaveActivity
		},
		func() {
			handleSyncTasks(app, state, deps)
		},
	)

	selectAction := GetNewAction("Select Activity", NewRuneKey("a", 'a'),
		func() bool {
			return !state.ActionsLock
		},
		func() {
			handleSelectActivity(app, pages, state)
		},
	)

	return &[]Action{closeAction, syncAction, selectAction}
}

func loadTasksActivity(app *tview.Application, deps *Dependencies, state *SyncState) {
	var wg sync.WaitGroup
	resultsChan := make(chan TaskToSyncActivities, len(state.Tasks))

	for row, task := range state.Tasks {
		wg.Add(1)
		go loadTaskActivity(&wg, resultsChan, app, deps, state, &task, row+1)
	}

	go func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines finished")
		close(resultsChan)

		allTasksWithDefaultActivity := true
		for result := range resultsChan {
			state.TasksActivities[result.Index] = result
			if result.Default.Name == "" {
				allTasksWithDefaultActivity = false
			}
		}

		state.AllTasksHaveActivity = allTasksWithDefaultActivity
		state.ActionsLock = false
		app.QueueUpdateDraw(state.ActionsManager.Refresh)
	}()
}

func loadTaskActivity(
	wg *sync.WaitGroup,
	resultsChan chan<- TaskToSyncActivities,
	app *tview.Application,
	deps *Dependencies,
	state *SyncState,
	task *types.TasksToSync,
	row int,
) {
	defer wg.Done()

	log.Println("Loading task activity for externalId:", task.ExternalId)
	activities, defaultActivity, err := deps.Redmine.LoadActivities(task.ExternalId)
	if err != nil {
		log.Println("Error loading task activity:", err)
		return
	}

	app.QueueUpdateDraw(func() {
		if defaultActivity.Name == "" {
			state.Table.SetCellText(row, 5, "[red]Select activity!")
		} else {
			state.Table.SetCellText(row, 5, "[green]"+defaultActivity.Name)
		}
	})

	resultsChan <- TaskToSyncActivities{
		Activities: activities,
		Default:    defaultActivity,
		Index:      row - 1,
	}
}

func handleSyncTasks(app *tview.Application, state *SyncState, deps *Dependencies) {
	log.Println("Syncing tasks...")

	state.ActionsLock = true
	state.ActionsManager.Refresh()

	var wg sync.WaitGroup

	for i, task := range state.Tasks {
		wg.Add(1)
		go syncTask(&wg, app, &task, state.TasksActivities[i].Default.Id, i+1, state, deps)
	}

	go func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines finished")

		state.ActionsLock = false
		app.QueueUpdateDraw(state.ActionsManager.Refresh)
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
		state.Table.SetCellText(row, 6, "[yellow]Syncing...!")
	})

	err := deps.Redmine.SendTask(task.ExternalId, task.Desc, task.Date, task.Duration, activityId)
	if err != nil {
		log.Println("Error syncing task:", err)
		app.QueueUpdateDraw(func() {
			state.Table.SetCellText(row, 6, "[red]Error!")
		})
		return
	}

	app.QueueUpdateDraw(func() {
		state.Table.SetCellText(row, 6, "[green]Success!")
	})

	for _, idStr := range task.Ids.IDs {
		id, _ := strconv.Atoi(idStr)
		deps.Service.SetTaskAsReported(uint(id))
	}
}

func getSelectedTaskToSync(state *SyncState) (types.TasksToSync, int, error) {
	row := state.Table.GetRowSelected()
	if row == 0 {
		return types.TasksToSync{}, 0, fmt.Errorf("there is nit a selected task")
	}
	return state.Tasks[row-1], row, nil
}

func handleSelectActivity(app *tview.Application, pages *tview.Pages, state *SyncState) {
	task, taskRow, err := getSelectedTaskToSync(state)
	if err != nil {
		return
	}

	taskActivities := state.TasksActivities[taskRow-1]

	log.Println("Select activity for task", task.Id, "with row", taskRow)

	options := []string{}
	currentOption := -1
	for i, activity := range *taskActivities.Activities {
		options = append(options, activity.Name)
		if activity.Id == taskActivities.Default.Id {
			currentOption = i
		}
	}

	log.Println("Current option is", currentOption)

	dropdown := tview.NewDropDown().SetLabel("Activity: ")
	dropdown.SetOptions(options, func(text string, index int) {}).SetCurrentOption(currentOption)

	form := tview.NewForm().
		AddTextView("Task: ", task.Desc, 0, 1, false, false).
		AddFormItem(dropdown)

	components.ShowFormModal("Select Activity", 80, 9, form, pages, app, func() {
		idx, _ := dropdown.GetCurrentOption()
		newActivity := (*taskActivities.Activities)[idx]
		log.Println("Option selected", newActivity)
		(*taskActivities.Default) = newActivity

		state.Table.SetCellText(taskRow, 5, "[green]"+newActivity.Name)
		state.checkAllTasksHaveDefaultActivity(state)
		state.Table.Deselect()
	})
}

func (s *SyncState) checkAllTasksHaveDefaultActivity(state *SyncState) {
	haveDefault := true
	for _, task := range s.TasksActivities {
		if task.Default.Name == "" {
			haveDefault = false
			break
		}
	}

	s.AllTasksHaveActivity = haveDefault
	state.ActionsManager.Refresh()
}
