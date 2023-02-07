package agent_tui

import (
	"time"

	"github.com/openshift/agent-installer-utils/tools/agent_tui/dialogs"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/forms"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/newt"
	"github.com/rivo/tview"
)

func App(app *tview.Application) {
	//New UI
	if app == nil {
		app = tview.NewApplication()
	}
	pages := tview.NewPages()

	background := tview.NewBox().
		SetBorder(false).
		SetBackgroundColor(newt.ColorBlue)

	pages.AddPage("background", background, true, true).
		AddPage("Node0", forms.IsNode0Modal(app, pages), true, true)

	controller := NewController(app)
	engine := NewChecksEngine(controller.GetChan())

	engine.Init()
	controller.Init()

	//UI init
	if err := app.SetRoot(pages, true).Run(); err != nil {
		dialogs.PanicDialog(app, err)
	}
}

// ChecksEngine is the model part, and is composed by a number
// of different checks.
// Each Check has a type, frequency and evaluation loop.
// Different checks could have the same type

type CheckResult struct {
	Type    string
	Success bool
	Details string // In case of failure
}

type Check struct {
	Type string
	Freq time.Duration //Note: a ticker could be useful
	Run  func(c chan CheckResult)
}

type ChecksEngine struct {
	checks []*Check
	c      chan CheckResult
}

func NewChecksEngine(c chan CheckResult) *ChecksEngine {
	checks := []*Check{}

	checkRegistry := &Check{
		Type: "Registry",
		Freq: 5 * time.Second,
		Run: func(c chan CheckResult) {
			for {
				//Do the check
				//...

				//Send the result
				res := CheckResult{}
				c <- res

				//Wait
			}
		},
	}
	checks = append(checks, checkRegistry)
	//.. add other checks

	return &ChecksEngine{
		checks: checks,
		c:      c,
	}
}

func (ce *ChecksEngine) Init() {
	for _, chk := range ce.checks {
		go chk.Run(ce.c)
	}
}

// Controller
type Controller struct {
	app    *tview.Application
	c      chan CheckResult
	status bool
}

func NewController(app *tview.Application) *Controller {
	return &Controller{
		c:   make(chan CheckResult),
		app: app,
	}
}

func (c *Controller) GetChan() chan CheckResult {
	return c.c
}

func (c *Controller) Init() {
	go func() {
		for {
			select {
			case r := <-c.c:

				//Update the internal state. If it changes, update the ui
				c.status = r.Success
				//...

				//Update the widgets
				switch r.Type {
				case "Registry":
					c.app.QueueUpdate(func() {
						//... Update the registry label
					})
				}
			}
		}
	}()
}
