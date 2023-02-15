package ui

import (
	"context"
	"log"
	"main/app/core"
	"main/app/ui/utils"

	"github.com/AkaCoder404/gothulearn"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SetUniversalHandlers : Set universal inputs for the app.
func SetUniversalHandlers() {
	// Enable mouse inputs.
	core.App.TView.EnableMouse(true)

	// Set universal keybindings
	core.App.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput()
		case tcell.KeyCtrlK: // Help page.
			// ctrlKInput()
		case tcell.KeyCtrlS: // Search page.
			// ctrlSInput()
		case tcell.KeyCtrlC: // Ctrl-C interrupt.
			// ctrlCInput()
		}
		return event // Forward the event to the actual current primitive.
	})
}

// ctrlInput : Enables user to toggle logout modal
func ctrlLInput() {
	log.Println("Received toggle logout modal event")
	// prevent logout modal from being shown when on login page
	if page, _ := core.App.PageHolder.GetFrontPage(); page == utils.LoginPageID {
		return
	}

	var modal *tview.Modal
	// Create modal for logout confirmation
	text := "Logout?\n Stored credentials will be deleted!"
	modal = confirmModal(utils.LoginLogoutCfmModalID, text, "Logout", func() {
		// attempt to logout
		// if err := core.App.Client.Logout(); err != nil {
		// 	okM := okModal(utils.GenericAPIErrorModalID, "Error logging out", err.Error())
		// 	ShowModal(utils.GenericAPIErrorModalID, okM)
		// 	return
		// }
		// If loggout successfully, then delete stored credentials and direct user to login page
		core.App.DeleteCredentials()
		ShowLoginPage()
	})

	ShowModal(utils.LoginLogoutCfmModalID, modal)
}

// setHandlers : Set handlers for the main page
func (p *MainPage) setHandlers(cancel context.CancelFunc, searchParams *SearchParams) {
	// TODO: Set table input captures

	// Set table entry selected functions
	p.CurrentCoursesTable.SetSelectedFunc(func(row, _ int) {
		log.Printf("Selected row %d on main page\n", row)
		classRef := p.CurrentCoursesTable.GetCell(row, 0).GetReference()
		if classRef == nil {
			return
		} else if class, ok := classRef.(*gothulearn.Class); ok {
			ShowClassPage(class)
		}
	})

	p.AllCoursesTable.SetSelectedFunc(func(row, _ int) {
		log.Printf("Selected row %d on main page\n", row)
		classRef := p.AllCoursesTable.GetCell(row, 0).GetReference()
		if classRef == nil {
			return
		} else if class, ok := classRef.(*gothulearn.ClassAll); ok {
			log.Println(class)
			// ShowClassPage(class)
			// TODO : cast ClassAll to Class object
			var classTemp gothulearn.Class
			classTemp.Wlkcid = class.Wlkcid
			ShowClassPage(&classTemp)
		}
	})

}

// setHandlers : Set handlers for the class page
func (p *ClassPage) setHandlers(cancel context.CancelFunc) {
	// Set grid input captures
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // go back to main page
			cancel()
			core.App.PageHolder.RemovePage(utils.ClassPageID)
		}
		return event
	})
}
