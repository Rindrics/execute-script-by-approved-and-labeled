package application

import "github.com/Rindrics/execute-script-with-merge/domain"

type Config struct {
	RquiredLabel  string
	DefaultBranch string
}

type App struct {
	Config                 Config
	ExecutionDirectiveList string
}

func New(requiredLabel, defaultBranch string) App {
	return App{
		Config: Config{
			RquiredLabel:  requiredLabel,
			DefaultBranch: defaultBranch,
		},
	}
}

func (a App) IsValid(event domain.ParsedEvent) bool {
	if a.HasRequiredLabel(event) && a.IsDefaultBranch(event) {
		return true
	}
	return false
}

func (a App) HasRequiredLabel(event domain.ParsedEvent) bool {
	if event.Labels.Contains(a.Config.RquiredLabel) {
		return true
	}
	return false
}

func (a App) IsDefaultBranch(event domain.ParsedEvent) bool {
	if event.Branch == a.Config.DefaultBranch {
		return true
	}
	return false
}
