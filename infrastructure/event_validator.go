package infrastructure

import (
	"log/slog"

	"github.com/Rindrics/execute-script-with-merge/application"
	"github.com/Rindrics/execute-script-with-merge/domain"
)

type ParsedEventValidator struct {
	Logger *slog.Logger
	Config application.Config
}

func (pev *ParsedEventValidator) Validate(event domain.ParsedEvent) bool {
	if !doExistRequiredLabel(event, pev.Config.RequiredLabel) {
		pev.Logger.Error("required label not found", "event", event)
		return false
	}

	if !doMatchDefaultBranch(event, pev.Config.DefaultBranch) {
		pev.Logger.Error("default branch did not match", "event", event)
		return false
	}

	if !isMerged(event) {
		pev.Logger.Error("event is not merged", "event", event)
		return false
	}

	return true
}

func doExistRequiredLabel(event domain.ParsedEvent, label string) bool {
	labels := event.Labels

	for _, v := range labels {
		if v == label {
			return true
		}
	}
	return false
}

func doMatchDefaultBranch(event domain.ParsedEvent, branch string) bool {
	return event.Branches.Base == branch
}

func isMerged(event domain.ParsedEvent) bool {
	return event.Merged
}
