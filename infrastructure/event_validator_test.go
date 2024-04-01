package infrastructure

import (
	"testing"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/stretchr/testify/assert"
)

func TestParsedEventValidator_Validate(t *testing.T) {
	t.Run("event is valid", func(t *testing.T) {
		validator := ParsedEventValidator{
			Logger: NewLogger(),
			Config: application.Config{
				RequiredLabel: "required-label",
				BaseBranch:    "main",
			},
		}

		assert.True(t, validator.Validate(domain.ParsedEvent{
			Labels: []string{"required-label"},
			Branches: domain.Branches{
				Base: "main",
			},
			Merged: true,
		}))
	})

	t.Run("required label not found", func(t *testing.T) {
		validator := ParsedEventValidator{
			Logger: NewLogger(),
			Config: application.Config{
				RequiredLabel: "required-label",
				BaseBranch:    "main",
			},
		}

		assert.False(t, validator.Validate(domain.ParsedEvent{
			Labels: []string{"foo"},
			Branches: domain.Branches{
				Base: "main",
			},
			Merged: true,
		}))
	})

	t.Run("base branch did not match", func(t *testing.T) {
		validator := ParsedEventValidator{
			Logger: NewLogger(),
			Config: application.Config{
				RequiredLabel: "required-label",
				BaseBranch:    "main",
			},
		}

		assert.False(t, validator.Validate(domain.ParsedEvent{
			Labels: []string{"required-label"},
			Branches: domain.Branches{
				Base: "foo",
			},
			Merged: true,
		}))
	})

	t.Run("event is not merged", func(t *testing.T) {
		validator := ParsedEventValidator{
			Logger: NewLogger(),
			Config: application.Config{
				RequiredLabel: "required-label",
				BaseBranch:    "main",
			},
		}

		assert.False(t, validator.Validate(domain.ParsedEvent{
			Labels: []string{"required-label"},
			Branches: domain.Branches{
				Base: "main",
			},
			Merged: false,
		}))
	})
}
