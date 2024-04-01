package infrastructure

import (
	"testing"

	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	logger := NewLogger()
	invoker := NewInvoker(logger)

	t.Run("run shell script", func(t *testing.T) {
		err := invoker.Execute(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{
				"hello_from_go.sh",
			},
			Directory: "assets/",
		})
		assert.Nil(t, err)
	})

	t.Run("python script", func(t *testing.T) {
		err := invoker.Execute(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{
				"foo.py",
			},
			Directory: "assets/",
		})
		assert.Nil(t, err) // unsupported extension
	})

	t.Run("unsupported script", func(t *testing.T) {
		err := invoker.Execute(domain.TargetScriptList{
			TargetScripts: []domain.TargetScript{
				"foo.R",
			},
			Directory: "assets/",
		})
		assert.NotNil(t, err) // unsupported extension
	})
}
