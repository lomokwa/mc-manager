package services

import (
	"log"
	"os/exec"
)

func StartServerProcess() (string, error) {
	log.Printf("executing start server command")
	cmd := exec.Command("echo", "Starting server...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("start server command failed: %v", err)
		return string(output), err
	}

	log.Printf("start server command completed")
	return string(output), err
}

func StopServerProcess() (string, error) {
	log.Printf("executing stop server command")
	cmd := exec.Command("echo", "Stopping server...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("stop server command failed: %v", err)
		return string(output), err
	}

	log.Printf("stop server command completed")
	return string(output), err
}
