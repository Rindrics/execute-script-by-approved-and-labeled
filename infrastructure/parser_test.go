package infrastructure_test

import (
	"os"
	"testing"

	"github.com/Rindrics/execute-script-with-merge/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestParsePullRequestEvent(t *testing.T) {
	t.Run("pull request", func(t *testing.T) {
		os.Setenv("GITHUB_EVENT_PATH", "./pull_request.json")
		event, err := infrastructure.ParsePullRequestEvent()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, event)
	})
	t.Run("others", func(t *testing.T) {
		os.Setenv("GITHUB_EVENT_PATH", "./invalid.json")
		_, err := infrastructure.ParsePullRequestEvent()
		assert.Error(t, err)
	})
	os.Unsetenv("GITHUB_EVENT_PATH")
}
