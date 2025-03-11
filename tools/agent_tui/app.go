package agent_tui

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/openshift/agent-installer-utils/pkg/version"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/checks"
	"github.com/openshift/agent-installer-utils/tools/agent_tui/ui"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

func App(app *tview.Application, rendezvousIP string, config checks.Config, checkFuncs ...checks.CheckFunctions) {

	if err := prepareConfig(&config); err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	// initialize log
	f, err := os.OpenFile(config.LogPath, os.O_RDWR|os.O_CREATE, 0644)
	if errors.Is(err, os.ErrNotExist) {
		// handle the case where the file doesn't exist
		fmt.Printf("Error creating log file %s\n", config.LogPath)
	}
	logger.Out = f

	logger.Infof("Release Image URL: %s", config.ReleaseImageURL)
	logger.Infof("Agent TUI git version: %s", version.Commit)
	logger.Infof("Agent TUI build version: %s", version.Raw)
	logger.Infof("NODE_ZERO_IP: %s", rendezvousIP)

	var appUI *ui.UI
	if app == nil {
		app = tview.NewApplication()
	}
	appUI = ui.NewUI(app, config, logger)
	controller := ui.NewController(appUI)
	engine := checks.NewEngine(controller.GetChan(), config, logger, checkFuncs...)

	controller.Init(engine.Size(), rendezvousIP)
	engine.Init()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func prepareConfig(config *checks.Config) error {
	// Set hostname
	hostname, err := checks.ParseHostnameFromURL(config.ReleaseImageURL)
	if err != nil {
		return err
	}
	config.ReleaseImageHostname = hostname

	// Set scheme
	schemeHostnamePort, err := checks.ParseSchemeHostnamePortFromURL(config.ReleaseImageURL, "https://")
	if err != nil {
		return fmt.Errorf("error creating <scheme>://<hostname>:<port> from releaseImageURL: %s", config.ReleaseImageURL)
	}
	config.ReleaseImageSchemeHostnamePort = schemeHostnamePort

	return nil
}
