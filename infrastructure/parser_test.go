package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePullRequestEvent(t *testing.T) {
	t.Run("pull request", func(t *testing.T) {
		os.Setenv("GITHUB_EVENT_PATH", "./pull_request.json")
		event, err := parsePullRequestEvent()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, event)
	})
	t.Run("others", func(t *testing.T) {
		t.Run("issue", func(t *testing.T) {
			os.Setenv("GITHUB_EVENT_PATH", "./issue.json")
			_, err := parsePullRequestEvent()
			assert.Error(t, err)
		})
		t.Run("invalid event", func(t *testing.T) {
			os.Setenv("GITHUB_EVENT_PATH", "./invalid.json")
			_, err := parsePullRequestEvent()
			assert.Error(t, err)
		})
	})
	os.Unsetenv("GITHUB_EVENT_PATH")
}
