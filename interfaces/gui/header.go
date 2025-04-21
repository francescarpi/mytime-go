package gui

import (
	"fmt"
	"log"
	"mytime/settings"
	"mytime/tasks"
	"mytime/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	tasksManager   *tasks.TasksManager
	container      *tview.Flex
	todayContainer *tview.TextView
	weekContainer  *tview.TextView
	dateContainer  *tview.TextView
}

func GetNewHeader(tasksManager *tasks.TasksManager) *Header {
	tv := tview.NewTextView()
	tv.SetTitle(" MyTime ").SetBorder(true)

	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle(" MyTime ")

	todayContainer := tview.NewTextView()
	weekContainer := tview.NewTextView()
	dateContainer := tview.NewTextView().SetTextAlign(tview.AlignRight).SetTextColor(tcell.ColorBlue)

	flex.AddItem(todayContainer, 0, 1, false)
	flex.AddItem(weekContainer, 0, 1, false)
	flex.AddItem(dateContainer, 0, 1, false)

	return &Header{
		tasksManager:   tasksManager,
		container:      flex,
		todayContainer: todayContainer,
		dateContainer:  dateContainer,
		weekContainer:  weekContainer,
	}
}

func (h *Header) Refresh() {
	log.Println("Refreshing header")

	dateStr := date.Format(time.DateOnly)
	workedToday, errToday := h.tasksManager.GetWorkedDaily(dateStr)
	workedThisWeek, errWeek := h.tasksManager.GetWorkedWeekly(dateStr)
	settings, errSettings := settings.GetSettings(h.tasksManager.Conn)

	if errToday != nil || errWeek != nil || errSettings != nil {
		panic("Error getting data from database")
	}

	goalDay := settings.GoalDayInSeconds(date)
	goalWeek := settings.GoalWeekInSeconds()

	h.todayContainer.SetText(fmt.Sprintf("Today: %s/%s",
		utils.HumanizeDuration(workedToday),
		utils.HumanizeDuration(goalDay),
	))

	if workedToday > 0 && workedToday >= goalDay {
		h.todayContainer.SetTextColor(tcell.ColorGreen)
	} else {
		h.todayContainer.SetTextColor(tcell.ColorRed)
	}

	h.weekContainer.SetText(fmt.Sprintf("Week: %s/%s",
		utils.HumanizeDuration(workedThisWeek),
		utils.HumanizeDuration(goalWeek),
	))

	if workedThisWeek >= goalWeek {
		h.weekContainer.SetTextColor(tcell.ColorGreen)
	} else {
		h.weekContainer.SetTextColor(tcell.ColorRed)
	}

	h.dateContainer.SetText(fmt.Sprintf("%s", date.Format("Monday, 2006-01-02")))
}
