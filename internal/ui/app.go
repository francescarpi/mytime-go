package ui

import "github.com/rivo/tview"

func StartApp() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	deps := InitDeps()

	pages.AddPage("home", HomeView(app, pages, deps), true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
