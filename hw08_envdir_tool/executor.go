package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors" //nolint: depguard // import is necessary
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	execCmd.Env = os.Environ()

	for key, envValue := range env {
		index := -1
		for i, s := range execCmd.Env {
			if strings.HasPrefix(s, fmt.Sprintf("%s=", key)) {
				index = i
				break
			}
		}
		if index != -1 {
			execCmd.Env = append(execCmd.Env[:index], execCmd.Env[index+1:]...)
		}

		if !envValue.NeedRemove {
			execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", key, envValue.Value))
		}
	}

	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()
	if err != nil {
		log.Println("failed to execCmd.Run")
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}
