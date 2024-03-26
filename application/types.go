package application

import (
	"github.com/Rindrics/execute-script-with-merge/domain"
)

type Config struct {
	RequiredLabel             string
	DefaultBranch             string
	ExecutionDirectiveListDir string
}

type App struct {
	Config                 Config
	ExecutionDirectiveList domain.ExecutionDirectiveList
	Parser                 domain.EventParser
	Logger                 Logger
}

func New(config Config, parser domain.EventParser, logger Logger) *App {
	logger.Debug("application.New", "config", config)
	return &App{
		Config: config,
		ExecutionDirectiveList: domain.ExecutionDirectiveList{
			Directory: config.ExecutionDirectiveListDir,
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
	if event.Labels.Contains(a.Config.RequiredLabel) {
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

func (a *App) LoadExecutionDirectives(event domain.ParsedEvent) error {
	executionDirectives, err := a.Parser.ParseExecutionDirectives(event, a.Config.ExecutionDirectiveListDir)
	if err != nil {
		return err
	}
	a.Logger.Info("application.LoadExecutionDirectives()", "executionDirectives", executionDirectives)
	a.ExecutionDirectiveList.ExecutionDirectives = executionDirectives

	return nil
}

func (a *App) LoadExecutionDirectiveList(event domain.ParsedEvent) error {
	a.ExecutionDirectiveList.Directory = a.Config.ExecutionDirectiveListDir
	if err := a.LoadExecutionDirectives(event); err != nil {
		return err
	}

	return nil
}

type ShellInvoker interface {
	Execute(domain.ExecutionDirectiveList) error
}

func (a *App) Run(invoker ShellInvoker) error {
	return invoker.Execute(a.ExecutionDirectiveList)
}

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Error(string, ...any)
}
