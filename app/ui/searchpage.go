package ui

import "github.com/rivo/tview"

// SearchPage : This struct contains the search bar and the table of results
// for the search. This struct reuses the MainPage struct, specifically for the guest table.
type SearchPage struct {
	MainPage
	Form *tview.Form
}

// SearchParams : Convenience struct to hold parameters for setting up a search table.
type SearchParams struct {
	term string // The term to search for
}
