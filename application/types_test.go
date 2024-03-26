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

func createApp(t *testing.T, parser domain.EventParser) *application.App {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	if parser == nil {
		parser = dmock.NewMockEventParser(ctrl)
	}
	logger := infrastructure.NewLogger()

	config := application.Config{
		RequiredLabel:             "required-label",
		DefaultBranch:             "main",
		ExecutionDirectiveListDir: "../infrastructure/assets/",
	}

	app := application.New(config, parser, logger)

	return app
}

func TestAppHasRequiredLabel(t *testing.T) {
	app := createApp(t, nil)

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
	app := createApp(t, nil)

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
	app := createApp(t, nil)

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

func TestAppLoadExecutionDirectiveList(t *testing.T) {
	expectedDirectives := []domain.ExecutionDirective{"foo.sh", "bar.sh"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockParser := dmock.NewMockEventParser(ctrl)
	mockParser.EXPECT().ParseExecutionDirectives(domain.ParsedEvent{}, "../infrastructure/assets/").Return(expectedDirectives, nil).Times(1)

	app := createApp(t, mockParser)

	t.Run("LoadExecutionDirectiveList", func(t *testing.T) {
		err := app.LoadExecutionDirectiveList(domain.ParsedEvent{})
		assert.Nil(t, err)
		assert.Equal(t, expectedDirectives, app.ExecutionDirectiveList.ExecutionDirectives,
			"The execution directives should match the expected values.")
	})
}

func TestAppRun(t *testing.T) {
	app := createApp(t, nil)

	app.ExecutionDirectiveList = domain.ExecutionDirectiveList{
		ExecutionDirectives: []domain.ExecutionDirective{"foo.sh", "bar.sh"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockShellInvoker := amock.NewMockShellInvoker(ctrl)
	mockShellInvoker.EXPECT().Execute(gomock.Eq(domain.ExecutionDirectiveList{
		ExecutionDirectives: []domain.ExecutionDirective{"foo.sh", "bar.sh"},
	})).Return(nil).Times(1)

	t.Run("Run", func(t *testing.T) {
		err := app.Run(mockShellInvoker)
		assert.Nil(t, err)
	})
}
