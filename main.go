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
		RequiredLabel:             "test-label",
		DefaultBranch:             "main",
		ExecutionDirectiveListDir: "infrastructure/assets/",
	}

	// TODO: remove EventParser from argument
	app := application.New(config, infrastructure.EventParser{}, logger)
	logger.Debug("main", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		return
	}

	// TODO: add label existence to validation
	if app.IsValid(event) {
		app.LoadExecutionDirectiveList(event)
		logger.Debug("main", "app.ExecutionDirectiveList", app.ExecutionDirectiveList)
		logger.Info("executing ExecutingDirectiveList")
		app.Run(infrastructure.NewShellInvoker(logger))
	} else {
		logger.Info("exit because event did not meet requirements")
	}
	logger.Info("finished application")
}
