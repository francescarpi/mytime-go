package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ModalPrimitive(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

// ShowAlertModal displays a simple alert with a message and an OK button.
func ShowAlertModal(app *tview.Application, pages *tview.Pages, message string, onClose func()) {
	modal := tview.NewModal().
		SetText(message).
		SetTextColor(tcell.ColorBlack).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("alertModal")
			if onClose != nil {
				onClose()
			}
		})
		modal.SetBorder(true)

	pages.AddPage("alertModal",
		tview.NewGrid().
			SetColumns(0, 50, 0).
			SetRows(0, 7, 0).
			AddItem(modal, 1, 1, 1, 1, 0, 0, true),
		true, true)

	app.SetFocus(modal)
}

func ShowConfirmModal(
	app *tview.Application,
	pages *tview.Pages,
	modalID string,
	message string,
	buttons []string,
	onDone func(buttonLabel string),
) {
	modal := tview.NewModal().
		SetText(message).
		SetTextColor(tcell.ColorBlack).
		SetButtonBackgroundColor(tcell.ColorBlue).
		SetButtonTextColor(tcell.ColorBlack).
		AddButtons(buttons).
		SetDoneFunc(func(_ int, buttonLabel string) {
			pages.RemovePage(modalID)
			if onDone != nil {
				onDone(buttonLabel)
			}
		})
		modal.SetBorder(true)

	pages.AddPage(
		modalID,
		tview.NewGrid().
			SetColumns(0, 50, 0).
			SetRows(0, 7, 0).
			AddItem(modal, 1, 1, 1, 1, 0, 0, true),
		true,
		true,
	)

	app.SetFocus(modal)
}
