package gui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetFooter() *tview.TextView {
	actions := [11][2]string{
		{"Quit", "q"},
		{"Date ←", "h"},
		{"Date →", "l"},
		{"Task ↓", "j"},
		{"Task ↑", "k"},
		{"Today", "t"},
		{"Start/Stop", "enter"},
		{"Duplicate", "d"},
		{"Modify", "m"},
		{"Delete", "x"},
		{"New", "n"},
	}

	footer := tview.NewTextView()
	footer.SetBorder(true)
	footer.SetTextColor(tcell.ColorWhite)
	footer.SetDynamicColors(true)

	text := ""
	for _, action := range actions {
		text += fmt.Sprintf("[darkgray]%s[white]: [blue]%s[white] | ", action[0], action[1])
	}
	footer.SetText(text)
	return footer
}
