package components

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

func ShowFormModal(
	title string,
	width, height int,
	form *tview.Form,
	pages *tview.Pages,
	app *tview.Application,
	onDone func(),
	onClose func(),
) {
	form.
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetCancelFunc(func() {
			pages.RemovePage("formModal")
			if onClose != nil {
				onClose()
			}
		})

	if onDone == nil {
		form = form.AddButton("OK", func() {
			pages.RemovePage("formModal")
			if onClose != nil {
				onClose()
			}
		})
	} else {
		form = form.
			AddButton("OK", func() {
				onDone()
				pages.RemovePage("formModal")
				if onClose != nil {
					onClose()
				}
			}).
			AddButton("Cancel", func() {
				pages.RemovePage("formModal")
				if onClose != nil {
					onClose()
				}
			})
	}

	layout := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	layout.SetTitle(title).SetBorder(true)

	pages.AddPage("formModal", ModalPrimitive(layout, width, height), true, true)

	// Show the modal
	app.SetFocus(form)
}
