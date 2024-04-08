package infrastructure

import (
	"os"
	"testing"

	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/stretchr/testify/assert"
	"github.com/waigani/diffparser"
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

func TestParseTargetScripts(t *testing.T) {
	pe := domain.ParsedEvent{
		Branches: domain.Branches{
			Base: "origin/main",
			Head: "origin/branch-for-test",
		},
	}

	parser := EventParser{
		"https://github.com/Rindrics/execute-scripts-github-flow.git",
		"token",
		NewLogger()}
	ed, err := parser.ParseTargetScripts(pe, "assets/target_script_list.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ed)
	assert.Equal(t, domain.TargetScript("for_test.sh"), ed[3])
	assert.Equal(t, domain.TargetScript("for_test2.sh"), ed[4])
}

func TestGetGitDiff(t *testing.T) {
	diff, err := getGitDiff(
		"https://github.com/Rindrics/execute-scripts-github-flow.git",
		"token",
		"origin/main",
		"origin/branch-for-test",
		"assets/target_script_list.txt",
		NewLogger(),
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, &diffparser.Diff{}, diff) // get non-empty diff
}

func TestParseTargetScripsFromGitDiff(t *testing.T) {
	logger := NewLogger()
	diff := &diffparser.Diff{
		Files: []*diffparser.DiffFile{
			{
				Hunks: []*diffparser.DiffHunk{
					{
						NewRange: diffparser.DiffRange{
							Lines: []*diffparser.DiffLine{
								{Content: "foo.sh"},
							},
						},
					},
				},
			},
		},
	}
	ts := parseTargetScriptsFromGitDiff(diff, logger)
	assert.NotNil(t, ts)
	assert.Equal(t, domain.TargetScript("foo.sh"), ts[0])
}
