package cmd

import (
	"os"
	"os/exec"
)

func createCommand(executable string, args ...string) *exec.Cmd {
	command := exec.Command(executable, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	return command
}
