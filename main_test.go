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

	config := application.Config{
		RequiredLabel:             "test-label",
		DefaultBranch:             "main",
		ExecutionDirectiveListDir: "infrastructure/assets/",
	}
	logger := infrastructure.NewLogger()
	app := application.New(config, infrastructure.EventParser{logger}, logger)
	logger.Debug("main.TestMainValidEvent", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		t.Fatal(err)
	}
	logger.Debug("main.TestMainValidEvent", "event:", event)
	assert.NotNil(t, event)

	assert.True(t, app.IsValid(event))

	app.LoadExecutionDirectiveList(event)
	logger.Debug("main", "app.ExecutionDirectiveList", app.ExecutionDirectiveList)
	err = app.Run(infrastructure.NewShellInvoker(logger))
	if err != nil {
		logger.Error("failed to run", "error", err)
		t.Fatal(err)
	}

	os.Unsetenv("GITHUB_EVENT_PATH")
}

func TestMainInvalidEvent(t *testing.T) {
	os.Setenv("GITHUB_EVENT_PATH", "./infrastructure/pull_request_opened.json")

	logger := infrastructure.NewLogger()
	config := application.Config{
		RequiredLabel:             "test-label",
		DefaultBranch:             "main",
		ExecutionDirectiveListDir: "infrastructure/assets/",
	}
	app := application.New(config, infrastructure.EventParser{logger}, logger)
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
