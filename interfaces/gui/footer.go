package gui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var actions = [][]any{
	{"Quit", "q", false},
	{"Date Prev.", "h", false},
	{"Date Next", "l", false},
	{"Task Next", "j", true},
	{"Task Prev.", "k", true},
	{"Today", "t", false},
	{"Start/Stop", "enter", true},
	{"Duplicate", "d", true},
	{"Modify", "m", true},
	{"Delete", "x", true},
	{"New", "n", false},
}

type Footer struct {
	container *tview.TextView
}

func GetNewFooter() *Footer {

	footer := tview.NewTextView()
	footer.SetBorder(true)
	footer.SetTextColor(tcell.ColorWhite)
	footer.SetDynamicColors(true)
	return &Footer{
		container: footer,
	}
}

func (f *Footer) Refresh() {
	text := ""
	for _, action := range actions {
		if tasksTable.taskSelected == nil && action[2].(bool) {
			text += fmt.Sprintf("[gray]%s: %s[white] | ", action[0], action[1])
		} else {
			text += fmt.Sprintf("[darkgray]%s[white]: [blue]%s[white] | ", action[0], action[1])
		}
	}

	f.container.SetText(text)
}
