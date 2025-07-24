package services

import (
	"os/exec"
)

func StartServerProcess() (string, error) {
	cmd := exec.Command("echo", "Starting server...")
	output, err := cmd.CombinedOutput()
	return string(output), err
}
