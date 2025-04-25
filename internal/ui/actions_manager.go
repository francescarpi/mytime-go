package ui

import (
	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Action struct {
	label        string
	key          string
	enabledFn    func() bool
	inputHandler func(event *tcell.EventKey) *tcell.EventKey
}

type ActionsManager struct {
	actions   *[]Action
	container *tview.TextView
}

func GetNewAction(label, key string, enabledFn func() bool, inputHandler func(event *tcell.EventKey) *tcell.EventKey) Action {
	return Action{
		label:        label,
		key:          key,
		enabledFn:    enabledFn,
		inputHandler: inputHandler,
	}
}

func GetNewActionsManager(container *tview.TextView, actions *[]Action) *ActionsManager {
	am := ActionsManager{
		actions,
		container,
	}
	am.Refresh()
	return &am
}

func (a *ActionsManager) Refresh() {
	result := ""
	for _, action := range *a.actions {
		result += util.Colorize(action.label, action.key, action.enabledFn())
	}
	a.container.SetText(result)
}

func (a *ActionsManager) GetInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		for _, action := range *a.actions {
			if action.enabledFn() {
				action.inputHandler(event)
			}
		}
		a.Refresh()
		return event
	}
}
