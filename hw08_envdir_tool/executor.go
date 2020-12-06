package main

import (
	"fmt"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	// RunCmd runs a command + arguments (cmd) with environment variables from env

	execEnv := make([]string, len(env))
	for envVariable, value := range env {
		// skip empty values
		if value == "" {
			continue
		}
		execEnv = append(execEnv, fmt.Sprintf("%s=%s", envVariable, value))
	}

	command, args := cmd[0], cmd[1:]
	cmdExec := exec.Command(command, args...)
	cmdExec.Env = execEnv

	stdoutStderr, err := cmdExec.CombinedOutput()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			return exitError.ExitCode()
		}

		// os.PathError and other non caugth errors
		return 1
	}

	fmt.Printf("%s", stdoutStderr)

	return 0
}
