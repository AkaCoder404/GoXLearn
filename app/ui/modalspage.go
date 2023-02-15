package ui

import (
	"log"
	"main/app/core"
	"main/app/ui/utils"

	"github.com/rivo/tview"
)

// ShowModal : Make the app show a modal
func ShowModal(id string, modal *tview.Modal) {
	core.App.TView.SetFocus(modal)
	core.App.PageHolder.AddPage(id, modal, true, true)
}

// okModal : Creates a new modal with an "OK" acknowledge button
func okModal(id, text string) *tview.Modal {
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.App.PageHolder.RemovePage((id))
		})

	return modal
}

// confirmModal: Creates a confirmation modal
func confirmModal(id, text, confirmButton string, f func()) *tview.Modal {
	// Create new modal
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		SetBackgroundColor(utils.ModalColor).
		AddButtons([]string{confirmButton, "Cancel"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, _ string) {
			if buttonIndex == 0 {
				f()
			}
			log.Printf("Removing %s modal\n", id)
			core.App.PageHolder.RemovePage(id)
		})
	return modal
}
