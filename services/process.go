package services

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

var (
	serverCmd   *exec.Cmd
	serverStdin io.WriteCloser
)

func StartServerProcess() (string, error) {
	if serverCmd != nil {
		return "", fmt.Errorf("server already running")
	}

	log.Printf("Starting Server...")
	cmd := exec.Command("java", "-Xms1G", "-Xmx2G", "-jar", "server.jar", "nogui")
	cmd.Dir = "./minecraft-server"

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	serverStdin = stdin

	if err := cmd.Start(); err != nil {
		log.Printf("start server command failed: %v", err)
		return "", fmt.Errorf("failed to start server: %w", err)
	}
	serverCmd = cmd

	go func() {
		cmd.Wait()
		log.Printf("server process exited")
		serverCmd = nil
		serverStdin = nil
	}()

	log.Printf("start server command completed")
	return "server started", nil
}

// TODO: Send "stop" command to Minecraft stdin for graceful shutdown
// TODO: Wait with timeout, then force-kill if process doesn't exit
// TODO: Clear stored process handle/PID after stop
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
