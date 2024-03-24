package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Rindrics/execute-script-with-merge/domain"
	"github.com/google/go-github/github"
)

type EventParser struct{}

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
	event, err := parsePullRequestEvent()
	if err != nil {
		return domain.ParsedEvent{}, err
	}
	var labels domain.Labels
	for _, label := range event.PullRequest.Labels {
		labels = append(labels, *label.Name)
	}

	return domain.ParsedEvent{
		Branch: *event.PullRequest.Head.Ref,
		Labels: labels,
	}, nil

}

func (e EventParser) ParseExecutionDirectives() ([]domain.ExecutionDirective, error) {
	return []domain.ExecutionDirective{}, nil
}
