package infrastructure

import (
	"io/ioutil"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/Rindrics/execute-script-with-merge/domain"
)

type TargetScriptListValidator struct {
	Logger *slog.Logger
}

func (slv *TargetScriptListValidator) Validate(list domain.TargetScriptList) bool {
	if len(list.TargetScripts) == 0 {
		return false
	}

	for _, script := range list.TargetScripts {
		_, ok := getScriptType(string(script))
		if !ok {
			slv.Logger.Error("Unsupported script extension", "script", string(script))
			return false
		}

		slv.Logger.Debug("infrastructure.TargetScriptListValidator.Validate", "validationTarget", list.Directory+string(script))
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
		slv.Logger.Debug("infrastructure.TargetScriptListValidator.Validate", "output", string(output))

		err = cmd.Wait()
		if err != nil {
			slv.Logger.Error("ValidateScriptInGitIndex", "Git ls-files command failed with", err)
			return false
		}

		// If the output is empty, the script does not exist in the git index
		if strings.TrimSpace(string(output)) == "" {
			return false
		}
	}

	return true
}

func getScriptType(fileName string) (domain.ScriptType, bool) {
	for ext, t := range domain.ScriptExtensionMapping {
		if strings.HasSuffix(fileName, ext) {
			return t, true
		}
	}
	return 0, false // unsupported extension
}
