package application_test

import (
	"testing"

	"github.com/Rindrics/execute-script-with-merge/application"
	amock "github.com/Rindrics/execute-script-with-merge/application/mock"
	"github.com/Rindrics/execute-script-with-merge/domain"
	dmock "github.com/Rindrics/execute-script-with-merge/domain/mock"
	"github.com/Rindrics/execute-script-with-merge/infrastructure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createApp(t *testing.T, ctrl *gomock.Controller, parser domain.EventParser, scriptValidator application.TargetScriptListValidator, eventValidator application.ParsedEventValidator) *application.App {
	if parser == nil {
		parser = dmock.NewMockEventParser(ctrl)
	}
	if scriptValidator == nil {
		scriptValidator = amock.NewMockTargetScriptListValidator(ctrl)
	}
	if eventValidator == nil {
		eventValidator = amock.NewMockParsedEventValidator(ctrl)
	}
	logger := infrastructure.NewLogger()

	config := application.Config{
		RequiredLabel:       "required-label",
		DefaultBranch:       "main",
		TargetScriptListDir: "../infrastructure/assets/",
	}

	app := application.New(config, parser, scriptValidator, eventValidator, logger)

	return app
}

func TestAppIsValid(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eventValidator := amock.NewMockParsedEventValidator(ctrl)
		eventValidator.EXPECT().Validate(gomock.Any()).Return(true).AnyTimes()

		app := createApp(t, ctrl, nil, nil, eventValidator)

		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "main",
			},
			Labels: domain.Labels{"required-label"},
			Merged: true,
		}
		assert.True(t, app.IsValid(event))
	})
	t.Run("InValid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eventValidator := amock.NewMockParsedEventValidator(ctrl)
		eventValidator.EXPECT().Validate(gomock.Any()).Return(false).AnyTimes()

		app := createApp(t, ctrl, nil, nil, eventValidator)

		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "main",
			},
			Labels: domain.Labels{},
			Merged: true,
		}
		assert.False(t, app.IsValid(event))
	})
}

func TestAppLoadTargetScriptList(t *testing.T) {
	expectedScripts := []domain.TargetScript{"foo.sh", "bar.sh"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := dmock.NewMockEventParser(ctrl)
	mockParser.EXPECT().ParseTargetScripts(domain.ParsedEvent{}, "../infrastructure/assets/").Return(expectedScripts, nil).Times(1)

	app := createApp(t, ctrl, mockParser, nil, nil)

	t.Run("LoadTargetScriptList", func(t *testing.T) {
		err := app.LoadTargetScriptList(domain.ParsedEvent{})
		assert.Nil(t, err)
		assert.Equal(t, expectedScripts, app.TargetScriptList.TargetScripts,
			"The target scripts should match the expected values.")
	})
}

func TestAppValidateTargetScriptList(t *testing.T) {
	t.Run("script exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		scriptValidator := amock.NewMockTargetScriptListValidator(gomock.NewController(t))
		scriptValidator.EXPECT().Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
		}).Return(true).Times(1)

		app := createApp(t, ctrl, nil, scriptValidator, nil)

		app.TargetScriptList = domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
		}

		assert.True(t, app.ValidateTargetScripts())
	})

	t.Run("unknown script", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		scriptValidator := amock.NewMockTargetScriptListValidator(gomock.NewController(t))
		scriptValidator.EXPECT().Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "unknown.sh"},
		}).Return(false).Times(1)

		app := createApp(t, ctrl, nil, scriptValidator, nil)

		app.TargetScriptList = domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "unknown.sh"},
		}

		assert.False(t, app.ValidateTargetScripts())
	})
}

func TestAppRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := createApp(t, ctrl, nil, nil, nil)

	app.TargetScriptList = domain.TargetScriptList{
		TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
	}

	mockShellInvoker := amock.NewMockShellInvoker(ctrl)
	mockShellInvoker.EXPECT().Execute(gomock.Eq(domain.TargetScriptList{
		TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
	})).Return(nil).Times(1)

	t.Run("Run", func(t *testing.T) {
		err := app.Run(mockShellInvoker)
		assert.Nil(t, err)
	})
}
