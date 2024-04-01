package main

import (
	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()

	logger.Info("starting application")

	// TODO: load config from environment variables
	config := application.Config{
		RequiredLabel:       "test-label",
		DefaultBranch:       "main",
		TargetScriptListDir: "infrastructure/assets/",
	}

	// TODO: remove EventParser from argument
	app := application.New(config, infrastructure.EventParser{}, &infrastructure.TargetScriptListValidator{logger}, logger)
	logger.Debug("main", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		return
	}

	// TODO: add label existence to validation
	if app.IsValid(event) {
		app.LoadTargetScriptList(event)
		logger.Debug("main", "app.TargetScriptList", app.TargetScriptList)
		if app.ValidateTargetScripts() {
			logger.Info("executing TargetScriptList")
			app.Run(infrastructure.NewInvoker(logger))
		} else {
			logger.Info("exit because TargetScriptList did not meet requirements")
		}
	} else {
		logger.Info("exit because event did not meet requirements")
	}
	logger.Info("finished application")
}
