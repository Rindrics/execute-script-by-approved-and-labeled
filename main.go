package main

import (
	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()

	logger.Info("starting application")

	config := application.Config{
		RequiredLabel:             "test-label",
		DefaultBranch:             "main",
		ExecutionDirectiveListDir: "infrastructure/assets/",
	}

	app := application.New(config, infrastructure.EventParser{}, logger)
	logger.Debug("main", "app:", app)
}
