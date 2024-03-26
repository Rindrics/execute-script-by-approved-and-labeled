package domain

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

type ExecutionDirective string

type EventParser interface {
	ParseEvent() (ParsedEvent, error)
	ParseExecutionDirectives(ParsedEvent, string) ([]ExecutionDirective, error)
}
