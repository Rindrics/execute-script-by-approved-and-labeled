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
	"github.com/google/go-github/github"
	"github.com/waigani/diffparser"
)

type EventParser struct {
	Logger *slog.Logger
}

func parsePullRequestEvent() (*github.PullRequestEvent, error) {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		return nil, fmt.Errorf("GITHUB_EVENT_PATH environment variable not set")
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
	}, nil
}

func (e EventParser) ParseExecutionDirectives(pe domain.ParsedEvent, edlPath string) ([]domain.ExecutionDirective, error) {
	e.Logger.Debug("infrastructure.ParseExecutionDirectives", "head", pe.Branches.Head, "base", pe.Branches.Base)
	diff, err := getGitDiff(pe.Branches.Base, pe.Branches.Head, edlPath, e.Logger)
	if err != nil {
		return []domain.ExecutionDirective{}, err
	}
	ed := parseExecutionDirectivesFromGitDiff(diff, e.Logger)
	e.Logger.Info("infrastructure.ParseExecutionDirectives", "ExecutionDirectives", ed)

	return ed, nil
}

func getGitDiff(base, head, targetFile string, logger *slog.Logger) (*diffparser.Diff, error) {
	// TODO:
	// - define application.Config.ExecutionDirectiveListDir as new type
	// - define Validate() to check whether it exists in git index
	cmd := exec.Command("git", "diff", "--no-color", base+"..."+head, "--", targetFile)
	logger.Debug("infrastructure.getGitDiff", "cmd", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return &diffparser.Diff{}, err
	}

	err = cmd.Start()
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "cmd.Start() failed with", err)
		return &diffparser.Diff{}, err
	}

	output, err := ioutil.ReadAll(stdout)
	logger.Debug("infrastructure.getGitDiff", "output", string(output))
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "ReadAll failed with", err)
		return &diffparser.Diff{}, err
	}

	err = cmd.Wait()
	if err != nil {
		logger.Error("infrastructure.getGitDiff", "cmd.Run() failed with", err)
		return &diffparser.Diff{}, err
	}

	return diffparser.Parse(string(output))
}

func parseExecutionDirectivesFromGitDiff(diff *diffparser.Diff, logger *slog.Logger) []domain.ExecutionDirective {
	executionDirectives := []domain.ExecutionDirective{}

	for _, file := range diff.Files {
		logger.Debug("infrastructure.parseExecutionDirectivesFromGitDiff", "file", file)
		for _, hunk := range file.Hunks {
			for _, line := range hunk.NewRange.Lines {
				logger.Debug("infrastructure.parseExecutionDirectivesFromGitDiff", "line", line)
				executionDirectives = append(executionDirectives, domain.ExecutionDirective(line.Content))
			}
		}
	}
	logger.Info("infrastructure.parseExecutionDirectivesFromGitDiff", "executionDirectives", executionDirectives)
	return executionDirectives
}
