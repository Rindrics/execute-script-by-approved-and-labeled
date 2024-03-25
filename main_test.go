package main

import (
	"os"
	"testing"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestMainValidEvent(t *testing.T) {
	os.Setenv("GITHUB_EVENT_PATH", "./infrastructure/pull_request.json")

	logger := infrastructure.NewLogger("debug")
	app := application.New("test-label", "main", infrastructure.EventParser{logger}, logger)
	logger.Debug("main.TestMainValidEvent", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		t.Fatal(err)
	}
	logger.Debug("main.TestMainValidEvent", "event:", event)
	assert.NotNil(t, event)

	assert.True(t, app.IsValid(event))
	os.Unsetenv("GITHUB_EVENT_PATH")
}
