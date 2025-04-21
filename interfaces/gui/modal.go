package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ModalForm struct {
	title    string
	form     *tview.Form
	width    int
	errorMsg *tview.TextView
	done     func(buttonLabel string)
	buttons  []string
}

func GetNewModalForm(title string, width int) *ModalForm {
	form := tview.NewForm()
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	form.SetButtonTextColor(tview.Styles.PrimaryTextColor)
	form.SetFieldBackgroundColor(tcell.ColorGray)

	return &ModalForm{
		title:    title,
		form:     form,
		width:    width,
		errorMsg: tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed),
		buttons:  []string{"OK", "Cancel"},
	}
}

func (m *ModalForm) SetForm(setForm func(form *tview.Form)) {
	setForm(m.form)
}

func (m *ModalForm) SetDoneFunc(done func(buttonLabel string)) {
	m.done = done
}

func (m *ModalForm) AddButtons(buttons []string) {
	m.buttons = buttons
}

func (m *ModalForm) SetErrorMsg(msg string) {
	m.errorMsg.SetText(msg)
}

func (m *ModalForm) Draw() tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true)
	flex.SetTitle(m.title)
	flex.SetBackgroundColor(tcell.ColorBlack)
	flex.AddItem(m.form, 0, 1, true)
	flex.AddItem(m.errorMsg, 2, 1, false)

	for _, button := range m.buttons {
		m.form.AddButton(button, func() {
			if m.done != nil {
				m.done(button)
			}
		})
	}

	m.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			if m.done != nil {
				m.done("Esc")
			}
			return nil
		}
		return event
	})

	height := (m.form.GetFormItemCount() * 4) + 2

	return modal(flex, m.width, height)
}
