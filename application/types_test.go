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

func createApp(t *testing.T, parser domain.EventParser, validator domain.TargetScriptListValidator) *application.App {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	if parser == nil {
		parser = dmock.NewMockEventParser(ctrl)
	}
	if validator == nil {
		validator = amock.NewMockTargetScriptListValidator(ctrl)
	}
	logger := infrastructure.NewLogger()

	config := application.Config{
		RequiredLabel:       "required-label",
		DefaultBranch:       "main",
		TargetScriptListDir: "../infrastructure/assets/",
	}

	app := application.New(config, parser, validator, logger)

	return app
}

func TestAppHasRequiredLabel(t *testing.T) {
	app := createApp(t, nil, nil)

	t.Run("HasRequiredLabel", func(t *testing.T) {
		event := domain.ParsedEvent{
			Labels: domain.Labels{"required-label"},
		}

		assert.True(t, app.HasRequiredLabel(event))
	})
	t.Run("NotHasRequiredLabel", func(t *testing.T) {
		event := domain.ParsedEvent{
			Labels: domain.Labels{"other-label"},
		}
		assert.False(t, app.HasRequiredLabel(event))
	})
}

func TestAppIsDefaultBranch(t *testing.T) {
	app := createApp(t, nil, nil)

	t.Run("IsDefaultBranch", func(t *testing.T) {
		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "main",
			},
		}
		assert.True(t, app.IsDefaultBranch(event))
	})
	t.Run("NotIsDefaultBranch", func(t *testing.T) {
		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "other-branhh",
			},
		}
		assert.False(t, app.IsDefaultBranch(event))
	})
}

func TestAppIsValid(t *testing.T) {
	app := createApp(t, nil, nil)

	t.Run("Valid", func(t *testing.T) {
		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "main",
			},
			Labels: domain.Labels{"required-label"},
		}
		assert.True(t, app.IsValid(event))
	})
	t.Run("NoLabel", func(t *testing.T) {
		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "main",
			},
			Labels: domain.Labels{},
		}
		assert.False(t, app.IsValid(event))
	})
	t.Run("NotDefaultBranch", func(t *testing.T) {
		event := domain.ParsedEvent{
			Branches: domain.Branches{
				Base: "other-branhh",
			},
			Labels: domain.Labels{"required-label"},
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

	app := createApp(t, mockParser, nil)

	t.Run("LoadTargetScriptList", func(t *testing.T) {
		err := app.LoadTargetScriptList(domain.ParsedEvent{})
		assert.Nil(t, err)
		assert.Equal(t, expectedScripts, app.TargetScriptList.TargetScripts,
			"The target scripts should match the expected values.")
	})
}

func TestAppValidateTargetScriptList(t *testing.T) {
	t.Run("script exist", func(t *testing.T) {
		validator := amock.NewMockTargetScriptListValidator(gomock.NewController(t))
		validator.EXPECT().Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
		}).Return(true).Times(1)

		app := createApp(t, nil, validator)

		app.TargetScriptList = domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
		}

		assert.True(t, app.ValidateTargetScripts())
	})

	t.Run("unknown script", func(t *testing.T) {
		validator := amock.NewMockTargetScriptListValidator(gomock.NewController(t))
		validator.EXPECT().Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "unknown.sh"},
		}).Return(false).Times(1)

		app := createApp(t, nil, validator)

		app.TargetScriptList = domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{"foo.sh", "unknown.sh"},
		}

		assert.False(t, app.ValidateTargetScripts())
	})
}

func TestAppRun(t *testing.T) {
	app := createApp(t, nil, nil)

	app.TargetScriptList = domain.TargetScriptList{
		TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockShellInvoker := amock.NewMockShellInvoker(ctrl)
	mockShellInvoker.EXPECT().Execute(gomock.Eq(domain.TargetScriptList{
		TargetScripts: []domain.TargetScript{"foo.sh", "bar.sh"},
	})).Return(nil).Times(1)

	t.Run("Run", func(t *testing.T) {
		err := app.Run(mockShellInvoker)
		assert.Nil(t, err)
	})
}
