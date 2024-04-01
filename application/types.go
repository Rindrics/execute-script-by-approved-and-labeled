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
	ScriptValidator  TargetScriptListValidator
	EventValidator   ParsedEventValidator
	Logger           Logger
}

func New(config Config, parser domain.EventParser, scriptValidator TargetScriptListValidator, eventValidator ParsedEventValidator, logger Logger) *App {
	logger.Debug("application.New", "config", config)
	return &App{
		Config: config,
		TargetScriptList: domain.TargetScriptList{
			Directory: config.TargetScriptListDir,
		},
		Parser:          parser,
		ScriptValidator: scriptValidator,
		EventValidator:  eventValidator,
		Logger:          logger,
	}
}

func (a *App) ParseEvent() (domain.ParsedEvent, error) {
	a.Logger.Debug("application.ParseEvent", "*App:", a)
	return a.Parser.ParseEvent()
}

func (a *App) IsValid(event domain.ParsedEvent) bool {
	return a.EventValidator.Validate(event)
}

func (a *App) ValidateTargetScripts() bool {
	return a.ScriptValidator.Validate(a.TargetScriptList)
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

type TargetScriptListValidator interface {
	Validate(list domain.TargetScriptList) bool
}

type ParsedEventValidator interface {
	Validate(event domain.ParsedEvent) bool
}
