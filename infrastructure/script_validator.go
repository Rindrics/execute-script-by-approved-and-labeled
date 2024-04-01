package infrastructure

import (
	"io/ioutil"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

type ScriptListValidator struct {
	Logger *slog.Logger
}

func (slv *ScriptListValidator) Validate(list domain.TargetScriptList) bool {
	for _, script := range list.TargetScripts {
		slv.Logger.Debug("infrastructure.ScriptListValidator.Validate", "validationTarget", list.Directory+string(script))
		cmd := exec.Command("git", "ls-files", list.Directory+string(script))
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			slv.Logger.Error("ValidateScriptInGitIndex", "Failed to create stdout pipe", err)
			return false
		}

		err = cmd.Start()
		if err != nil {
			slv.Logger.Error("ValidateScriptInGitIndex", "Git ls-files command failed with", err)
			return false
		}

		output, err := ioutil.ReadAll(stdout)
		if err != nil {
			slv.Logger.Error("ValidateScriptInGitIndex", "Failed to read output of git ls-files command", err)
			return false
		}
		slv.Logger.Debug("infrastructure.ScriptListValidator.Validate", "output", string(output))

		err = cmd.Wait()
		if err != nil {
			slv.Logger.Error("ValidateScriptInGitIndex", "Git ls-files command failed with", err)
			return false
		}

		// If the output is not empty, the script exists in the git index
		if strings.TrimSpace(string(output)) != "" {
			return true
		}
	}

	return false
}
