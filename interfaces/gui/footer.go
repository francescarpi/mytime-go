package gui

import (
	"fmt"
	"mytime/tasks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Action struct {
	Name     string
	Key      string
	Disabled func() bool
	Format   func(tasksManager *tasks.TasksManager) string
}

func formatSync(tasksManager *tasks.TasksManager) string {
	numTasks := len(tasksManager.GetTasksToSync())
	if numTasks == 0 {
		return "Sync (0)"
	}
	return fmt.Sprintf("Sync ([red]%d[white])", numTasks)
}

func disabledIfNoTaskSelected() bool {
	return tasksTable.taskSelected == nil
}

var actions = []Action{
	{"Quit", "q", nil, nil},
	{"Date Prev.", "h", nil, nil},
	{"Date Next", "l", nil, nil},
	{"Task Next", "j", disabledIfNoTaskSelected, nil},
	{"Task Prev.", "k", disabledIfNoTaskSelected, nil},
	{"Today", "t", nil, nil},
	{"Start/Stop", "enter", disabledIfNoTaskSelected, nil},
	{"Duplicate", "d", disabledIfNoTaskSelected, nil},
	{"Modify", "m", disabledIfNoTaskSelected, nil},
	{"Delete", "x", disabledIfNoTaskSelected, nil},
	{"New", "n", nil, nil},
	{"Sync", "s", nil, formatSync},
}

type Footer struct {
	tasksManager *tasks.TasksManager
	container    *tview.TextView
}

func GetNewFooter(tasksManager *tasks.TasksManager) *Footer {

	footer := tview.NewTextView()
	footer.SetBorder(true)
	footer.SetTextColor(tcell.ColorWhite)
	footer.SetDynamicColors(true)
	return &Footer{
		tasksManager: tasksManager,
		container:    footer,
	}
}

func (f *Footer) Refresh() {
	text := ""
	for _, action := range actions {
		actionName := action.Name
		if action.Format != nil {
			actionName = action.Format(f.tasksManager)
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

	f.container.SetText(text)
}
