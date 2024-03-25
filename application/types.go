package application

import (
	"github.com/Rindrics/execute-script-with-merge/domain"
)

type Config struct {
	RquiredLabel              string
	DefaultBranch             string
	ExecutionDirectiveListDir string
}

type App struct {
	Config                 Config
	ExecutionDirectiveList domain.ExecutionDirectiveList
	Parser                 domain.EventParser
	Logger                 Logger
}

func New(requiredLabel, defaultBranch string, parser domain.EventParser, logger Logger) *App {
	logger.Debug("application.New", "requiredLabel:", requiredLabel, "defaultBranch", defaultBranch)
	return &App{
		Config: Config{
			RquiredLabel:  requiredLabel,
			DefaultBranch: defaultBranch,
		},
		Parser: parser,
		Logger: logger,
	}
}

func (a *App) ParseEvent() (domain.ParsedEvent, error) {
	a.Logger.Debug("application.ParseEvent", "*App:", a)
	return a.Parser.ParseEvent()
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
	if event.Branches.Base == a.Config.DefaultBranch {
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

type ShellInvoker interface {
	Execute(Config, domain.ExecutionDirectiveList) error
}

func (a *App) Run(invoker ShellInvoker) error {
	return invoker.Execute(a.Config, a.ExecutionDirectiveList)
}

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Error(string, ...any)
}
