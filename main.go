package main

import (
	"os"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
)

func main() {
	os.Setenv(domain.EnvVarLogLevel, "debug")
	logger := infrastructure.NewLogger()

	logger.Info("starting application")

	logger.Debug("main", "event path", os.Getenv(domain.EnvVarGitHubEventPath))

	config, err := infrastructure.LoadConfig()
	logger.Debug("main", "config", *config)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}
	// TODO: remove EventParser from argument

	app := application.New(config, infrastructure.EventParser{os.Getenv(domain.EnvVarGitHubRepositoryUrl), config.Token, logger}, &infrastructure.TargetScriptListValidator{logger}, &infrastructure.ParsedEventValidator{logger, *config}, logger)
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
