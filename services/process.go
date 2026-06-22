package services

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
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

func StopServerProcess() (string, error) {
	log.Printf("executing stop server command")
	if serverCmd == nil || serverStdin == nil {
		return "", fmt.Errorf("server is not running")
	}

 	_, err := serverStdin.Write([]byte("stop\n")); err != nil {
		log.Printf("failed to send stop command: %v", err)
		return "", fmt.Errorf("failed to send stop command: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- serverCmd.Wait()
	}()

	select {
	case <-done:
		log.Printf("server stopped gracefully")
		serverCmd = nil
		serverStdin = nil
		return "server stopped", nil

	case <-time.After(30 * time.Second):
		log.Printf("server did not stop in time, force killing")
		serverCmd.Process.kill()
		serverCmd = nil
		serverStdin = nil
		return "server force-killed after timeout", nil
	}

}
