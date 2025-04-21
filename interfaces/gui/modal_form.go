package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ModalForm struct {
	title    string
	form     *tview.Form
	width    int
	height   int
	done     func(m *ModalForm)
	cancel   func(m *ModalForm)
	errorMsg *tview.TextView
}

func GetNewModalForm(title string) *ModalForm {
	form := tview.NewForm()
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	form.SetButtonTextColor(tview.Styles.PrimaryTextColor)
	form.SetFieldBackgroundColor(tcell.ColorGray)

	return &ModalForm{
		title:    title,
		form:     form,
		width:    50,
		height:   10,
		errorMsg: tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed),
	}
}

func (m *ModalForm) SetDoneFunc(done func(m *ModalForm)) {
	m.done = done
}

func (m *ModalForm) SetCancelFunc(cancel func(m *ModalForm)) {
	m.cancel = cancel
}

func (m *ModalForm) SetForm(setForm func(form *tview.Form)) {
	setForm(m.form)
}

func (m *ModalForm) addButtons() {
	m.form.AddButton("OK", func() {
		if m.done != nil {
			m.done(m)
		}
	})
	m.form.AddButton("Cancel", func() {
		if m.cancel != nil {
			m.cancel(m)
		}
	})
}

func (m *ModalForm) SetDimensions(width, height int) {
	m.width = width
	m.height = height
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

	m.addButtons()

	return modal(flex, m.width, m.height)
}
