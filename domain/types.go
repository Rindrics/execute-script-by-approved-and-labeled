package domain

import (
	"log/slog"
)

type Labels []string

func (l Labels) Contains(label string) bool {
	for _, v := range l {
		if v == label {
			return true
		}
	}
	return false
}

type Branches struct {
	Base string
	Head string
}

type ParsedEvent struct {
	Branches Branches
	Labels   Labels
}

type ExecutionDirectiveList struct {
	ExecutionDirectives []ExecutionDirective
	Directory           string
}

func (e *ExecutionDirectiveList) LoadExecutionDirectives(parser EventParser) error {
	executionDirectives, err := parser.ParseExecutionDirectives()
	if err != nil {
		return err
	}
	slog.Info("ExecutionDirectiveList():", "executionDirectives", executionDirectives)
	e.ExecutionDirectives = executionDirectives

	return nil
}

type ExecutionDirective string

type EventParser interface {
	ParseEvent() (ParsedEvent, error)
	ParseExecutionDirectives() ([]ExecutionDirective, error)
}
