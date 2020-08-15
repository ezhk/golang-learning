package main

import (
	"os"
)

func main() {
	envDir, cmd := os.Args[1], os.Args[2:]
	environment, err := ReadDir(envDir)
	if err != nil {
		os.Exit(1)
	}

	exitCode := RunCmd(cmd, environment)
	os.Exit(exitCode)
}
