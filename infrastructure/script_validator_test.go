package infrastructure

import (
	"testing"

	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/stretchr/testify/assert"
)

func TestScriptListValidatorValidate(t *testing.T) {
	validator := ScriptListValidator{NewLogger()}

	t.Run("a script exists", func(t *testing.T) {
		assert.True(t, validator.Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{
				"hello_from_go.sh",
			},
			Directory: "assets/",
		}))
	})

	t.Run("scripts exist", func(t *testing.T) {
		assert.True(t, validator.Validate(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{
				"foo.sh",
				"bar.sh",
			},
			Directory: "assets/",
		}))
	})
}
