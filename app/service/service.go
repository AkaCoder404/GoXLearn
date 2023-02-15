package service

import (
	"main/app/core"
	"main/app/ui"

	"github.com/AkaCoder404/gothulearn"
	"github.com/rivo/tview"
)

// Start : Set up the application.
func Start() {
	// create application
	core.App = &core.GoLearn{
		Client:     gothulearn.NewLearnClient(),
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}

	// TODO: show appropriate screen based on restore session
	if err := core.App.Initialize(); err != nil {
		ui.ShowLoginPage()
	} else {
		ui.ShowMainPage()
	}
	ui.SetUniversalHandlers()

	core.App.TView.Run()
}

// Shutdown : Shutdown the application.
func Shutdown() {
	core.App.Shutdown()
}
