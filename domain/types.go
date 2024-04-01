package domain

type Labels []string

type Branches struct {
	Base string
	Head string
}

type ParsedEvent struct {
	Branches Branches
	Labels   Labels
	Merged   bool
}

type TargetScriptList struct {
	TargetScripts []TargetScript
	Directory     string
}

type TargetScript string

type EventParser interface {
	ParseEvent() (ParsedEvent, error)
	ParseTargetScripts(ParsedEvent, string) ([]TargetScript, error)
}

type TargetScriptListValidator interface {
	Validate(list TargetScriptList) bool
}

type ScriptType int

const (
	Bash ScriptType = iota
	Python
)

var ScriptExtensionMapping = map[string]ScriptType{
	".sh": Bash,
	".py": Python,
}

var ScriptCommandMapping = map[ScriptType][]string{
	Bash:   {"/bin/bash"},
	Python: {"python"},
}
