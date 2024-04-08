package main

import (
	"os"
	"testing"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestMainValid(t *testing.T) {
	os.Setenv(domain.EnvVarGitHubEventPath, "./infrastructure/pull_request.json")
	os.Setenv(domain.EnvVarGitHubRepositoryUrl, "https://github.com/Rindrics/execute-scripts-github-flow.git")
	os.Setenv(domain.EnvVarToken, "token")
	os.Setenv(domain.EnvVarRequiredLabel, "test-label")
	os.Setenv(domain.EnvVarBaseBranch, "main")
	os.Setenv(domain.EnvVarTargetScriptListDir, "infrastructure/assets")

	config, err := infrastructure.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	logger := infrastructure.NewLogger()
	app := application.New(
		config,
		infrastructure.EventParser{
			os.Getenv(domain.EnvVarGitHubRepositoryUrl),
			config.Token,
			logger,
		},
		&infrastructure.TargetScriptListValidator{logger},
		&infrastructure.ParsedEventValidator{logger, *config},
		logger,
	)
	logger.Debug("main.TestMainValidEvent", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		t.Fatal(err)
	}
	logger.Debug("main.TestMainValidEvent", "event:", event)
	assert.NotNil(t, event)

	assert.True(t, app.IsValid(event))

	app.LoadTargetScriptList(event)
	logger.Debug("main", "app.TargetScriptList", app.TargetScriptList)

	// cannot assert result on test
	app.ValidateTargetScripts()

	err = app.Run(infrastructure.NewInvoker(logger))
	if err != nil {
		logger.Error("failed to run", "error", err)
		t.Fatal(err)
	}

	os.Unsetenv(domain.EnvVarGitHubEventPath)
	os.Unsetenv(domain.EnvVarRequiredLabel)
	os.Unsetenv(domain.EnvVarBaseBranch)
	os.Unsetenv(domain.EnvVarTargetScriptListDir)
}

func TestMainInvalidEvent(t *testing.T) {
	os.Setenv(domain.EnvVarGitHubEventPath, "./infrastructure/pull_request_opened.json")
	os.Setenv(domain.EnvVarGitHubRepositoryUrl, "https://github.com/Rindrics/execute-scripts-github-flow.git")
	os.Setenv(domain.EnvVarToken, "token")
	os.Setenv(domain.EnvVarRequiredLabel, "test-label")
	os.Setenv(domain.EnvVarBaseBranch, "main")
	os.Setenv(domain.EnvVarTargetScriptListDir, "infrastructure/assets")

	logger := infrastructure.NewLogger()
	config, err := infrastructure.LoadConfig()
	app := application.New(
		config,
		infrastructure.EventParser{
			os.Getenv(domain.EnvVarGitHubRepositoryUrl),
			config.Token,
			logger,
		},
		&infrastructure.TargetScriptListValidator{logger},
		&infrastructure.ParsedEventValidator{logger, *config},
		logger,
	)
	logger.Debug("main.TestMainValidEvent", "app:", app)

	event, err := app.ParseEvent()
	if err != nil {
		logger.Error("failed to parse event", "error", err)
		t.Fatal(err)
	}
	logger.Debug("main.TestMainValidEvent", "event:", event)
	assert.NotNil(t, event)

	assert.False(t, app.IsValid(event))

	os.Unsetenv(domain.EnvVarGitHubEventPath)
	os.Unsetenv(domain.EnvVarRequiredLabel)
	os.Unsetenv(domain.EnvVarBaseBranch)
	os.Unsetenv(domain.EnvVarTargetScriptListDir)
}
