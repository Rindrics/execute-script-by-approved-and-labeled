package application

import "github.com/Rindrics/execute-script-with-merge/domain"

type Config struct {
	RquiredLabel              string
	DefaultBranch             string
	ExecutionDirectiveListDir string
}

type App struct {
	Config                 Config
	ExecutionDirectiveList domain.ExecutionDirectiveList
	Parser                 domain.EventParser
}

func New(requiredLabel, defaultBranch string, parser domain.EventParser) *App {
	return &App{
		Config: Config{
			RquiredLabel:  requiredLabel,
			DefaultBranch: defaultBranch,
		},
		Parser: parser,
	}
}

func (a *App) IsValid(event domain.ParsedEvent) bool {
	if a.HasRequiredLabel(event) && a.IsDefaultBranch(event) {
		return true
	}
	return false
}

func (a *App) HasRequiredLabel(event domain.ParsedEvent) bool {
	if event.Labels.Contains(a.Config.RquiredLabel) {
		return true
	}
	return false
}

func (a *App) IsDefaultBranch(event domain.ParsedEvent) bool {
	if event.Branch == a.Config.DefaultBranch {
		return true
	}
	return false
}

func (a *App) LoadExecutionDirectiveList() error {
	a.ExecutionDirectiveList.Directory = a.Config.ExecutionDirectiveListDir
	if err := a.ExecutionDirectiveList.LoadExecutionDirectives(a.Parser); err != nil {
		return err
	}

	return nil
}
