package main

import (
	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()

	logger.Info("starting application")

	app := application.New("test-label", "main", infrastructure.EventParser{}, logger)
	logger.Debug("main", "app:", app)
}
