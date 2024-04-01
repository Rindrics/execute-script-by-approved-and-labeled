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

func (s ShellInvoker) executeShellScript(dir string, ed domain.ExecutionDirective) error {
	cmd := exec.Command("/bin/bash", path.Join(dir, string(ed)))
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

func (s ShellInvoker) Execute(edl domain.TargetScriptList) error {
	s.Logger.Debug("infrastructure.ShellInvoker.Execute", "Directory:", edl.Directory, "ExecutionDirectives:", edl.ExecutionDirectives)

	for _, ed := range edl.ExecutionDirectives {
		s.Logger.Debug("infrastructure.ShellInvoker.Execute", "ExecutionDirective:", ed)
		err := s.executeShellScript(edl.Directory, ed)
		if err != nil {
			return err
		}
	}

	return nil
}
