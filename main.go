package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
)

func main() {
	os.Setenv(domain.EnvVarLogLevel, "debug")
	logger := infrastructure.NewLogger()

	logger.Info("starting application")

	logger.Debug("main", "path", os.Getenv("PATH"))

	echo, err := exec.Command("echo", "$PATH").Output()
	fmt.Printf("echo:\n%s :Error:\n%v\n", echo, err)

	git, err := exec.Command("git", "--version").Output()
	fmt.Printf("git:\n%s :Error:\n%v\n", git, err)

	logger.Debug("main", "event path", os.Getenv(domain.EnvVarGitHubEventPath))

	config, err := infrastructure.LoadConfig()
	logger.Debug("main", "config", *config)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}
	// TODO: remove EventParser from argument
	app := application.New(config, infrastructure.EventParser{logger}, &infrastructure.TargetScriptListValidator{logger}, &infrastructure.ParsedEventValidator{logger, *config}, logger)
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
