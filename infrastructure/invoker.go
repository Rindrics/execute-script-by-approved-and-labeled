package infrastructure

import (
	"io/ioutil"
	"log/slog"
	"os/exec"
	"path"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

type Invoker struct {
	Logger *slog.Logger
}

func NewInvoker(logger *slog.Logger) Invoker {
	return Invoker{
		Logger: logger,
	}
}

func (s Invoker) executeScript(dir string, ts domain.TargetScript) error {
	scriptType, _ := getScriptType(ts)

	commandArgs := domain.ScriptCommandMapping[scriptType]
	filePath := path.Join(dir, string(ts))
	commandArgs = append(commandArgs, filePath)

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	s.Logger.Debug("infrastructure.Invoker.executeScript", "cmd", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	s.Logger.Debug("infrastructure.Invoker.executeScript", "output", string(output))

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s Invoker) Execute(tsl domain.TargetScriptList) error {
	s.Logger.Debug("infrastructure.Invoker.Execute", "Directory:", tsl.Directory, "TargetScripts:", tsl.TargetScripts)

	for _, ts := range tsl.TargetScripts {
		s.Logger.Debug("infrastructure.Invoker.Execute", "TargetScript:", ts)
		err := s.executeScript(tsl.Directory, ts)
		if err != nil {
			return err
		}
	}

	return nil
}
