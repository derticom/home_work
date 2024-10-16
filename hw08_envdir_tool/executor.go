package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors" //nolint: depguard // import is necessary
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, envValue := range env {
		if env[key].NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Printf("failed to os.Unsetenv for key %v", key)
				return 1
			}
			continue
		}
		err := os.Setenv(key, envValue.Value)
		if err != nil {
			log.Printf("failed to os.Setenv for key %v", key)
			return 1
		}
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

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
