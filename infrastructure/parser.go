package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"github.com/Rindrics/execute-script-with-merge/domain"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"github.com/waigani/diffparser"
)

type EventParser struct {
	Url    string
	Token  string
	Logger *slog.Logger
}

func parsePullRequestEvent() (*github.PullRequestEvent, error) {
	eventPath := os.Getenv(domain.EnvVarGitHubEventPath)
	if eventPath == "" {
		return nil, fmt.Errorf("%s environment variable not set", domain.EnvVarGitHubEventPath)
	}

	data, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the data into a generic map
	var genericData map[string]interface{}
	if err := json.Unmarshal(data, &genericData); err != nil {
		return nil, err
	}

	// Return if pull_request event
	if _, ok := genericData["pull_request"]; ok {
		var event github.PullRequestEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return &event, nil
	} else {
		return nil, fmt.Errorf("unknown event type")
	}
}

func (e EventParser) ParseEvent() (domain.ParsedEvent, error) {
	e.Logger.Debug("infrastructure.ParseEvent", "EventParser", e)

	event, err := parsePullRequestEvent()
	if err != nil {
		return domain.ParsedEvent{}, err
	}
	e.Logger.Debug("infrastructure.ParseEvent", "event:", event)
	var labels domain.Labels
	for _, label := range event.PullRequest.Labels {
		labels = append(labels, *label.Name)
	}
	e.Logger.Debug("infrastructure.ParseEvent", "labels:", labels)

	return domain.ParsedEvent{
		Branches: domain.Branches{
			Head: *event.PullRequest.Head.Ref,
			Base: *event.PullRequest.Base.Ref,
		},
		Labels: labels,
		Merged: *event.PullRequest.Merged,
	}, nil
}

func (e EventParser) ParseTargetScripts(pe domain.ParsedEvent, tslPath string) ([]domain.TargetScript, error) {
	e.Logger.Debug("infrastructure.ParseTargetScripts", "head", pe.Branches.Head, "base", pe.Branches.Base)
	diff, err := getGitDiff(e.Url, e.Token, pe.Branches.Base, pe.Branches.Head, tslPath, e.Logger)
	if err != nil {
		return []domain.TargetScript{}, err
	}
	ts := parseTargetScriptsFromGitDiff(diff, e.Logger)
	e.Logger.Info("infrastructure.ParseTargetScripts", "TargetScripts", ts)

	return ts, nil
}

func getGitDiff(url, token, base, head, targetFile string, logger *slog.Logger) (*diffparser.Diff, error) {
	logger.Info("infrastructure.getGitDiff", "cloning", url)
	dir, err := os.MkdirTemp("", "clone-example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/git-fixtures/basic.git",
		Auth: &http.BasicAuth{
			Username: "execute-script-with-merge",
			Password: token,
		},
	})
	// TODO:
	// - define application.Config.ExecutionDirectiveListDir as new type
	// - define Validate() to check whether it exists in git index
	output, err := ExecuteCommandWithLogging(logger, "ls", "-al", targetFile)
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "failed with", err)
	}
	logger.Debug("infrastructure.getGitDiff", "output", output)

	branches, err := repo.Branches()
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "repo.Branches() fails with:", err)
	}
	logger.Debug("infrastructure.getGitDiff", "branches", branches)

	output, err = ExecuteCommandWithLogging(logger, "git", "show")
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "failed with", err)
	}
	logger.Debug("infrastructure.getGitDiff", "output", output)

	output, err = ExecuteCommandWithLogging(logger, "git", "diff", "--no-color", base+".."+head, "--", targetFile)
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "failed with", err)
	}
	logger.Debug("infrastructure.getGitDiff", "output", output)

	return diffparser.Parse(string(output))
}

func ExecuteCommandWithLogging(logger *slog.Logger, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	logger.Debug("infrastructure.ExecuteCommandWithLogging", "cmd", cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		logger.Debug("infrastructure.ExecuteCommandWithLogging", "cmd.Start() failed with", err)
		return "", err
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		logger.Debug("infrastructure.ExecuteCommandWithLogging", "ReadAll() failed with", err)
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		logger.Debug("infrastructure.ExecuteCommandWithLogging", "cmd.Run() failed with", err)
		return "", err
	}

	return string(output), nil
}

func parseTargetScriptsFromGitDiff(diff *diffparser.Diff, logger *slog.Logger) []domain.TargetScript {
	targetScripts := []domain.TargetScript{}

	for _, file := range diff.Files {
		logger.Debug("infrastructure.parseTargetScriptsFromGitDiff", "file", file)
		for _, hunk := range file.Hunks {
			for _, line := range hunk.NewRange.Lines {
				logger.Debug("infrastructure.parseTargetScriptsFromGitDiff", "line", line)
				targetScripts = append(targetScripts, domain.TargetScript(line.Content))
			}
		}
	}
	logger.Info("infrastructure.parseTargetScriptsFromGitDiff", "targetScripts", targetScripts)
	return targetScripts
}
