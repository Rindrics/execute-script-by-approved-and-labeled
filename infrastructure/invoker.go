package infrastructure

import (
	"io/ioutil"
	"log/slog"
	"os/exec"
	"path"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

type ShellInvoker struct {
	Logger *slog.Logger
}

func NewShellInvoker(logger *slog.Logger) ShellInvoker {
	return ShellInvoker{
		Logger: logger,
	}
}

func (s ShellInvoker) executeShellScript(dir string, ts domain.TargetScript) error {
	scriptType, _ := getScriptType(ts)

	commandArgs := domain.ScriptCommandMapping[scriptType]
	filePath := path.Join(dir, string(ts))
	commandArgs = append(commandArgs, filePath)

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	s.Logger.Debug("infrastructure.ShellInvoker.executeShellScript", "cmd", cmd.String())
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
	s.Logger.Debug("infrastructure.ShellInvoker.executeShellScript", "output", string(output))

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s ShellInvoker) Execute(tsl domain.TargetScriptList) error {
	s.Logger.Debug("infrastructure.ShellInvoker.Execute", "Directory:", tsl.Directory, "TargetScripts:", tsl.TargetScripts)

	for _, ts := range tsl.TargetScripts {
		s.Logger.Debug("infrastructure.ShellInvoker.Execute", "TargetScript:", ts)
		err := s.executeShellScript(tsl.Directory, ts)
		if err != nil {
			return err
		}
	}

	return nil
}
