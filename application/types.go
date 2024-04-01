package application

import (
	"github.com/Rindrics/execute-script-with-merge/domain"
)

type Config struct {
	RequiredLabel       string
	DefaultBranch       string
	TargetScriptListDir string
}

type App struct {
	Config           Config
	TargetScriptList domain.TargetScriptList
	Parser           domain.EventParser
	Logger           Logger
}

func New(config Config, parser domain.EventParser, logger Logger) *App {
	logger.Debug("application.New", "config", config)
	return &App{
		Config: config,
		TargetScriptList: domain.TargetScriptList{
			Directory: config.TargetScriptListDir,
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

func (a *App) LoadTargetScripts(event domain.ParsedEvent) error {
	targetScripts, err := a.Parser.ParseTargetScripts(event, a.Config.TargetScriptListDir)
	if err != nil {
		return err
	}
	a.Logger.Info("application.LoadTargetScripts()", "targetScripts", targetScripts)
	a.TargetScriptList.TargetScripts = targetScripts

	return nil
}

func (a *App) LoadTargetScriptList(event domain.ParsedEvent) error {
	a.TargetScriptList.Directory = a.Config.TargetScriptListDir
	if err := a.LoadTargetScripts(event); err != nil {
		return err
	}

	return nil
}

type ShellInvoker interface {
	Execute(domain.TargetScriptList) error
}

func (a *App) Run(invoker ShellInvoker) error {
	return invoker.Execute(a.TargetScriptList)
}

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Error(string, ...any)
}
