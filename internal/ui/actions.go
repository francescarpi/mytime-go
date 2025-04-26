package ui

import (
	"log"

	"github.com/francescarpi/mytime/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ActionKey struct {
	Name string
	Rune rune
	Key  tcell.Key
}

type Action struct {
	label        string
	key          ActionKey
	enabledFn    func() bool
	inputHandler func()
}

type ActionsManager struct {
	actions   *[]Action
	container *tview.TextView
}

func NewRuneKey(name string, r rune) ActionKey {
	return ActionKey{Name: name, Rune: r}
}

func NewSpecialKey(name string, k tcell.Key) ActionKey {
	return ActionKey{Name: name, Key: k}
}

func GetNewAction(label string, key ActionKey, enabledFn func() bool, inputHandler func()) Action {
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
		result += util.Colorize(action.label, action.key.Name, action.enabledFn())
	}
	a.container.SetText(result)
}

func (a *ActionsManager) GetInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		for _, action := range *a.actions {
			if action.enabledFn() {
				if (action.key.Rune != 0 && event.Key() == tcell.KeyRune && event.Rune() == action.key.Rune) ||
					(action.key.Key != 0 && event.Key() == action.key.Key) {
					log.Printf("Action triggered: %s (%s)\n", action.label, action.key.Name)
					action.inputHandler()
					break
				}
			}
		}
		a.Refresh()
		return event
	}
}
