package ui

import (
	"sync/atomic"

	"github.com/openshift/agent-installer-utils/tools/agent_tui/checks"
	"github.com/rivo/tview"
)

type UI struct {
	app                 *tview.Application
	pages               *tview.Pages
	grid                *tview.Grid     // layout for the checks page
	envVars             *tview.TextView // displays release image URL
	checks              *tview.Table    // summary of all checks
	details             *tview.TextView // where errors from checks are displayed
	form                *tview.Form     // contains "Configure network" button
	timeoutModal        *tview.Modal    // popup window that times out
	nmtuiActive         atomic.Bool
	timeoutDialogActive atomic.Bool
	timeoutDialogCancel chan bool
}

func NewUI(app *tview.Application, config checks.Config) *UI {
	ui := &UI{
		app:                 app,
		timeoutDialogCancel: make(chan bool),
	}
	ui.create(config)
	return ui
}

func (u *UI) GetApp() *tview.Application {
	return u.app
}

func (u *UI) GetPages() *tview.Pages {
	return u.pages
}

func (u *UI) returnFocusToChecks() {
	u.pages.SwitchToPage(CHECK_PAGE_NAME)
	// shifting focus back to the "Configure network"
	// button requires setting focus in this sequence
	// form -> form-button
	u.app.SetFocus(u.form)
	u.app.SetFocus(u.form.GetButton(0))
}

func (u *UI) IsNMTuiActive() bool {
	return u.nmtuiActive.Load()
}

func (u *UI) setIsTimeoutDialogActive(isActive bool) {
	u.timeoutDialogActive.Store(isActive)
}

func (u *UI) isTimeoutDialogActive() bool {
	return u.timeoutDialogActive.Load()
}

func (u *UI) create(config checks.Config) {
	u.pages = tview.NewPages()
	u.createCheckPage(config)
	u.createTimeoutModal(config)
}
