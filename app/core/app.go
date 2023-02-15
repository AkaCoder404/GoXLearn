package core

import (
	"fmt"
	"os"

	"github.com/AkaCoder404/gothulearn"
	"github.com/rivo/tview"
)

// App : Global App Variable
var (
	App        *GoLearn
	AppVersion = "GoLearn v0.0.1"
)

type GoLearn struct {
	Client     *gothulearn.LearnClient
	TView      *tview.Application
	PageHolder *tview.Pages
	LogFile    *os.File
}

// Initialize the App
func (m *GoLearn) Initialize() error {
	// Set up logging
	if err := m.setUpLogging(); err != nil {
		fmt.Println("Unable to set up logging...")
		fmt.Println("Application will not continue.")
		os.Exit(1)
	}

	m.TView.SetRoot(m.PageHolder, true).SetFocus(m.PageHolder)

	return m.RestoreSession()
}

// Shutdown : Stop all services such as logging and let the application shut down gracefully.
func (m *GoLearn) Shutdown() {
	// Stop all necessary services, such as logging.

	// Sync the screen to make sure that the terminal screen is not corrupted.
	App.TView.Sync()
	App.TView.Stop()

	// Stop the logging
	if err := m.stopLogging(); err != nil {
		fmt.Println("Error while closing log file!")
	}
	fmt.Println("Application shutdown.")
}
