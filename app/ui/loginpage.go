package ui

import (
	"log"
	"main/app/core"
	"main/app/ui/utils"

	"github.com/rivo/tview"
)

// LoginPage : This struct contains the grid and form for the login page.
type LoginPage struct {
	Grid *tview.Grid
	Form *tview.Form
}

// ShowLoginPage : Make the app show the login page.
func ShowLoginPage() {
	// Create the new login page
	loginPage := newLoginPage()

	core.App.TView.SetFocus(loginPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.LoginPageID, loginPage.Grid, true)
}

// func
func newLoginPage() *LoginPage {
	// Create the LoginPage
	loginPage := &LoginPage{}

	// Create the form
	form := tview.NewForm()

	// Set form attributes.
	form.SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(utils.LoginFormLabelColor).
		SetTitle("Login to GoXLearn").
		SetTitleColor(utils.LoginPageTitleColor).
		SetBorder(true).
		SetBorderColor(utils.LoginFormBorderColor)

	// Add form fields.
	form.AddInputField("账号", "", 0, nil, nil).
		AddPasswordField("密码", "", 0, '*', nil).
		AddCheckbox("记得我", false, nil).
		AddButton("Login", func() {
			loginPage.attemptLogin()
		}).
		AddButton("Quit", func() {
			core.App.PageHolder.RemovePage(utils.LoginPageID)
			core.App.Shutdown()
		})

	// Setup Login Page
	dimension := []int{0, 0, 0}
	grid := utils.NewGrid(dimension, dimension)

	grid.AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	loginPage.Grid = grid
	loginPage.Form = form
	return loginPage
}

// attemptLogin : Attempts to log in with given form fields. If success, bring user to main page.
func (p *LoginPage) attemptLogin() {
	log.Println("Attempting to log in...")
	// Get username and password input.
	form := p.Form
	user := form.GetFormItemByLabel("账号").(*tview.InputField).GetText()
	pwd := form.GetFormItemByLabel("密码").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("记得我").(*tview.Checkbox).IsChecked()

	// Attempt to Log into Learn
	if err := core.App.Client.Auth.Login(user, pwd); err != nil {
		log.Printf("Error trying to log	in: %s\n", err.Error())
		modal := okModal(utils.GenericAPIErrorModalID, "Authentication failed\n Try again!")
		ShowModal(utils.GenericAPIErrorModalID, modal)
		return
	}

	// TODO: Remember the user's login credentials
	if remember {
		if err := core.App.StoreCredentials(user, pwd); err != nil {
			log.Printf("Error storing credentials: %s\n", err.Error())
			modal := okModal(utils.StoreCredentialErrorModalID,
				"Failed to store login token.\nCheck logs for details.")
			ShowModal(utils.StoreCredentialErrorModalID, modal)
		}
	}

	log.Println("Log in successfully")
	core.App.PageHolder.RemovePage(utils.LoginPageID) // Remove the login page as we no longer need it.
	ShowMainPage()
}
