package infrastructure

import (
	"testing"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

func TestExecute(t *testing.T) {
	logger := NewLogger()
	invoker := NewShellInvoker(logger)

	err := invoker.Execute(domain.TargetScriptList{
		TargetScripts: []domain.TargetScript{
			"hello_from_go.sh",
		},
		Directory: "assets/",
	})
	if err != nil {
		t.Fatal(err)
	}
}
