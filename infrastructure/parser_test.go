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

func TestParseExecutionDirectives(t *testing.T) {
	pe := domain.ParsedEvent{
		Branches: domain.Branches{
			Base: "main",
			Head: "branch-for-test",
		},
	}

	parser := EventParser{NewLogger()}
	ed, err := parser.ParseExecutionDirectives(pe, "assets/execution_directive_list.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ed)
	assert.Equal(t, domain.ExecutionDirective("for_test.sh"), ed[3])
	assert.Equal(t, domain.ExecutionDirective("for_test2.sh"), ed[4])
}

func TestGetGitDiff(t *testing.T) {
	diff, err := getGitDiff("main", "branch-for-test", "assets/execution_directive_list.txt", NewLogger())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, &diffparser.Diff{}, diff) // get non-empty diff
}

func TestParseExecutionDirectivesFromGitDiff(t *testing.T) {
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
	ed := parseExecutionDirectivesFromGitDiff(diff, logger)
	assert.NotNil(t, ed)
	assert.Equal(t, domain.ExecutionDirective("foo.sh"), ed[0])
}
