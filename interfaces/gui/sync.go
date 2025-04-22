package gui

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"mytime/integrations"
	"mytime/tasks"

	"github.com/rivo/tview"
)

type ActivityResult struct {
	Id              string
	DefaultActivity int
}

func loadActivity(
	wg *sync.WaitGroup,
	resultsChain chan<- ActivityResult,
	app *tview.Application,
	table *tview.Table,
	row int,
	id, externalId string,
	redmine *integrations.Redmine,
) {

	defer wg.Done()

	resp, err := redmine.LoadActivities(externalId)
	if err != nil {
		log.Println("Error loading activities:", err)
		return
	}

	defActivityId, err := strconv.Atoi(redmine.Config.DefaultActivity)
	if err != nil {
		log.Println("Error converting default activity ID:", err)
		return
	}

	// Find the activity in the response
	var defActivity integrations.RedmineProjectActivity
	for _, activity := range *resp {
		if activity.Id == defActivityId {
			defActivity = activity
			break
		}
	}

	if defActivity.Id != 0 {
		app.QueueUpdateDraw(func() {
			table.GetCell(row, 5).SetText(fmt.Sprintf("[green]%s", defActivity.Name))
		})
		resultsChain <- ActivityResult{
			Id:              id,
			DefaultActivity: defActivity.Id,
		}
	} else {
		resultsChain <- ActivityResult{
			Id:              id,
			DefaultActivity: 0,
		}
	}

}

func syncTask(
	wg *sync.WaitGroup,
	task *tasks.TasksToSync,
	activityId int,
	row int,
	app *tview.Application,
	table *tview.Table,
	redmine *integrations.Redmine,
	tasksManager *tasks.TasksManager,
) {
	defer wg.Done()
	log.Println("Syncking task", task.Id)

	app.QueueUpdateDraw(func() {
		table.GetCell(row, 6).SetText(fmt.Sprintf("[blue]*"))
	})

	err := redmine.SendTask(task.ExternalId, task.Desc, task.Date, task.Duration, activityId)
	if err != nil {
		log.Println("Error syncing task:", err)
		app.QueueUpdateDraw(func() {
			table.GetCell(row, 6).SetText(fmt.Sprintf("[red]E"))
		})
		return
	}

	app.QueueUpdateDraw(func() {
		table.GetCell(row, 6).SetText(fmt.Sprintf("[green]✓"))
	})

	for _, id := range task.Ids.Ids {
		uintId, _ := strconv.ParseUint(id, 10, 32)
		tasksManager.MarkTaskAsReported(uint(uintId))
	}
}

func syncTasks(
	tasksToSync *[]tasks.TasksToSync,
	activitiesByTask *map[string]int,
	app *tview.Application,
	table *tview.Table,
	redmine *integrations.Redmine,
	callback func(),
	tasksManager *tasks.TasksManager,
) {
	log.Println("Sync tasks", tasksToSync)

	var wg sync.WaitGroup

	for i, task := range *tasksToSync {
		row := i + 1
		activity := (*activitiesByTask)[task.Id]
		wg.Add(1)
		go syncTask(&wg, &task, activity, row, app, table, redmine, tasksManager)
	}

	go func() {
		log.Println("Waiting for all goroutines to finish...")
		wg.Wait()
		log.Println("All goroutines finished")
		callback()
	}()
}
