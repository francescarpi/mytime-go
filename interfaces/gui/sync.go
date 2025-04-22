package gui

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"mytime/integrations"

	"github.com/rivo/tview"
)

type ActivityResult struct {
	Id              string
	DefaultActivity int
}

func loadActivity(wg *sync.WaitGroup, resultsChain chan<- ActivityResult, app *tview.Application,
	table *tview.Table, row int, id, externalId string, redmine *integrations.Redmine) {

	log.Println("Reduce 1 to Wg")
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
