package gui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"mytime/tasks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TaskSelected struct {
	ID         uint
	Project    *string
	Desc       string
	ExternalId *string
}

type TasksTable struct {
	container    *tview.Flex
	table        *tview.Table
	tasksManager *tasks.TasksManager
	taskSelected *TaskSelected
}

const (
	LABEL_ID    = "ID"
	GO_NEXT     = 1
	GO_PREVIOUS = -1
)

// GetNewTasksTable creates a new TasksTable with a given tasks manager.
func GetNewTasksTable(tasksManager *tasks.TasksManager) *TasksTable {
	table := tview.NewTable().SetBorders(true)

	// Set table headers
	table.SetCell(0, 0, tview.NewTableCell(" "))
	table.SetCell(0, 1, tview.NewTableCell(LABEL_ID).SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 2, tview.NewTableCell("Project").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 3, tview.NewTableCell("Description").SetTextColor(tcell.ColorYellow).SetExpansion(1))
	table.SetCell(0, 4, tview.NewTableCell("Ext.ID").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 5, tview.NewTableCell("Started").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 6, tview.NewTableCell("Ended").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 7, tview.NewTableCell("Duration").SetTextColor(tcell.ColorYellow))
	table.SetCell(0, 8, tview.NewTableCell("Reported").SetTextColor(tcell.ColorYellow))

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, false)

	return &TasksTable{
		table:        table,
		container:    flex,
		tasksManager: tasksManager,
		taskSelected: nil,
	}
}

// LoadTasks is a placeholder for loading tasks into the table.
func (t *TasksTable) LoadTasks() {}

// Refresh reloads the tasks from the database and updates the table.
func (t *TasksTable) Refresh(selectTask bool) {
	t.clear()

	log.Println("Refreshing tasks:", date.Format(time.DateOnly))
	list, err := t.tasksManager.GetTasksByDate(date.Format(time.DateOnly))
	if err != nil {
		panic("Error getting tasks from database")
	}

	log.Println("Tasks found: ", len(list))

	if selectTask {
		if len(list) == 0 {
			t.taskSelected = nil
		} else {
			t.taskSelected = &TaskSelected{
				ID:         list[0].ID,
				Project:    list[0].Project,
				Desc:       list[0].Desc,
				ExternalId: list[0].ExternalId,
			}
		}
	}

	for i, task := range list {
		log.Println("Adding task", task.ID, "to table")
		row := i + 1 // Plus 1 because row 0 is the header

		if t.taskSelected != nil && task.ID == t.taskSelected.ID {
			t.table.SetCell(row, 0, tview.NewTableCell("*").SetTextColor(tcell.ColorYellow))
		} else {
			t.table.SetCell(row, 0, tview.NewTableCell(" "))
		}

		t.table.SetCell(row, 1, tview.NewTableCell(strconv.FormatUint(uint64(task.ID), 10)))
		t.table.SetCell(row, 2, tview.NewTableCell(task.GetProjectOrBlank()))
		t.table.SetCell(row, 3, tview.NewTableCell(task.Desc).SetExpansion(1))
		t.table.SetCell(row, 4, tview.NewTableCell(task.GetExternalIdOrBlank()))
		t.table.SetCell(row, 5, tview.NewTableCell(task.Start.Format(time.Kitchen)).SetAlign(tview.AlignRight))
		t.table.SetCell(row, 6, tview.NewTableCell(task.GetEndFormatedOr("")).SetAlign(tview.AlignRight))
		t.table.SetCell(row, 7, tview.NewTableCell(task.GetDurationHumanized()).SetAlign(tview.AlignRight))
		t.table.SetCell(row, 8, tview.NewTableCell(task.GetReportedIcon()).SetAlign(tview.AlignCenter))
	}
}

// clear removes all rows from the table except the header.
func (t *TasksTable) clear() {
	rowsToRemove := t.table.GetRowCount() - 1
	for i := rowsToRemove; i > 0; i-- {
		log.Println("Removing row ", i)
		t.table.RemoveRow(i)
	}
}

// GetTaskSelectedIndex returns the index of the currently selected task.
func (t *TasksTable) GetTaskSelectedIndex() int {
	if t.taskSelected == nil {
		log.Println("No task selected")
		return -1
	}

	for i := 1; i < t.table.GetRowCount(); i++ {
		cell := t.table.GetCell(i, 0)
		if cell.Text == "*" {
			return i
		}
	}
	return -1
}

// SelectNextPrevious selects the next or previous task based on the operation parameter.
func (t *TasksTable) SelectNextPrevious(operation int) bool {
	taskSelectedIndex := t.GetTaskSelectedIndex()
	if taskSelectedIndex == -1 {
		return false
	}

	log.Println("Current index of task selected is", taskSelectedIndex)

	index := taskSelectedIndex + (1 * operation)
	if index < 1 || index >= t.table.GetRowCount() {
		log.Println("No next task")
		return false
	}

	taskId := t.table.GetCell(index, 1).Text
	project := t.table.GetCell(index, 2).Text
	desc := t.table.GetCell(index, 3).Text
	externalId := t.table.GetCell(index, 4).Text

	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		log.Println("Error converting task ID to int:", err)
		return false
	}

	t.taskSelected = &TaskSelected{
		ID:         uint(taskIdInt),
		Project:    &project,
		Desc:       desc,
		ExternalId: &externalId,
	}

	return true
}

// StartStopTask toggles the start/stop state of the selected task.
func (t *TasksTable) StartStopTask() {
	if t.taskSelected == nil {
		log.Println("No task selected")
		return
	}
	log.Println("Starting/stopping task", t.taskSelected.ID)
	t.tasksManager.StartStopTask(t.taskSelected.ID)
}

// DuplicateWithDescription duplicates the selected task with a new description.
func (t *TasksTable) DuplicateWithDescription() {
	if t.taskSelected == nil {
		log.Println("No task selected")
		return
	}
	log.Println("Duplicating task", t.taskSelected.ID)

	var description string
	m := GetNewModalForm(" Duplicate Task ", 50)

	m.SetForm(func(form *tview.Form) {
		form.AddTextView("Project", *t.taskSelected.Project, 0, 1, false, false)
		form.AddTextView("External ID", *t.taskSelected.ExternalId, 0, 1, false, false)
		form.AddInputField("Description", "", 0, nil, func(text string) { description = text })
	})

	m.SetDoneFunc(func(buttonLabel string) {
		log.Println("Button pressed:", buttonLabel)
		if buttonLabel == "OK" {
			if description == "" {
				m.SetErrorMsg("Description cannot be empty")
				return
			}
			t.tasksManager.DuplicateTaskWithDescription(t.taskSelected.ID, description)
			GotoToday()
			fullRefresh(true)
		}
		pages.RemovePage("modal")
	})

	pages.AddPage("modal", m.Draw(), true, true)
}

// Modify opens a modal to modify the selected task's details.
func (t *TasksTable) Modify() {
	if t.taskSelected == nil {
		log.Println("No task selected")
		return
	}
	log.Println("Modifying task", t.taskSelected.ID)

	task, err := t.tasksManager.GetTaskById(t.taskSelected.ID)
	if (err != nil) || (task.ID == 0) {
		log.Println("Error getting task from database")
		return
	}

	project := ""
	if task.Project != nil {
		project = *task.Project
	}

	description := task.Desc

	externalId := ""
	if task.ExternalId != nil {
		externalId = *task.ExternalId
	}

	started := task.Start.Format("15:04")
	ended := ""

	if task.End != nil {
		ended = task.End.Format("15:04")
	}

	m := GetNewModalForm(" Modify Task ", 50)

	m.SetForm(func(form *tview.Form) {
		form.AddInputField("Project", project, 0, nil, func(text string) { project = text })
		form.AddInputField("Description", description, 0, nil, func(text string) { description = text })
		form.AddInputField("Ext.ID", externalId, 0, nil, func(text string) { externalId = text })
		form.AddInputField("Started", started, 0, nil, func(text string) { started = text })
		form.AddInputField("Ended", ended, 0, nil, func(text string) { ended = text })
	})

	m.SetDoneFunc(func(buttonLabel string) {
		log.Println("Button pressed:", buttonLabel)
		pages.RemovePage("modal")
		if buttonLabel == "OK" {
			t.tasksManager.Update(t.taskSelected.ID, project, description, externalId, started, ended)
			GotoToday()
			fullRefresh(true)
		}
	})

	pages.AddPage("modal", m.Draw(), true, true)
}

// Delete removes the selected task after confirmation.
func (t *TasksTable) Delete() {
	if t.taskSelected == nil {
		log.Println("No task selected")
		return
	}
	log.Println("Delete task", t.taskSelected.ID)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete task %d?", t.taskSelected.ID)).
		AddButtons([]string{"Yes", "No"}).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				t.tasksManager.Delete(t.taskSelected.ID)
				GotoToday()
				fullRefresh(true)
			}
			pages.RemovePage("modal")
		})
	pages.AddPage("modal", modal, true, true)
}

// New creates a new task with the current date and time.
func (t *TasksTable) New() {
	log.Println("New task")

	var project, description, externalId string

	m := GetNewModalForm(" New Task ", 50)

	m.SetForm(func(form *tview.Form) {
		form.AddInputField("Project", project, 0, nil, func(text string) { project = text })
		form.AddInputField("Description", description, 0, nil, func(text string) { description = text })
		form.AddInputField("Ext.ID", externalId, 0, nil, func(text string) { externalId = text })
	})

	m.SetDoneFunc(func(buttonLabel string) {
		log.Println("Button pressed:", buttonLabel)
		if buttonLabel == "OK" {
			if description == "" {
				m.SetErrorMsg("Description cannot be empty")
				return
			}
			t.tasksManager.New(project, description, externalId)
			GotoToday()
			fullRefresh(true)
		}
		pages.RemovePage("modal")
	})

	pages.AddPage("modal", m.Draw(), true, true)
}

// Sync synchronizes the tasks with (at the moment) Redmine...
func (t *TasksTable) Sync() {
	log.Println("Syncing tasks...")

	tasks := t.tasksManager.GetTasksToSync()
	log.Println("Tasks to sync:", len(tasks))
}
