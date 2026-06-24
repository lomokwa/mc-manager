package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/lomokwa/mc-manager/types"
)

var (
	serverCmd   *exec.Cmd
	serverStdin io.WriteCloser
	logHub      *types.LogHub
	stdinMu     sync.Mutex
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

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("start server command failed: %v", err)
		return "", fmt.Errorf("failed to start server: %w", err)
	}
	serverCmd = cmd

	logHub = types.NewLogHub()

	scanner := bufio.NewScanner(stdout)
	ready := make(chan string, 1)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
			logHub.Broadcast(line)
			if strings.Contains(line, "Done") {
				ready <- line
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error reading server output: %v", err)
		}
		ready <- ""
	}()

	go func() {
		err := cmd.Wait()
		log.Printf("server process exited: %v", err)
		if logHub != nil {
			logHub.Close()
		}
		serverCmd = nil
		serverStdin = nil
		logHub = nil
	}()

	select {
	case line := <-ready:
		if line == "" {
			return "", fmt.Errorf("server process exited before becoming ready")
		}
		return line, nil
	case <-time.After(120 * time.Second):
		cmd.Process.Kill()
		serverCmd = nil
		serverStdin = nil
		return "", fmt.Errorf("server failed to start within 120 seconds")
	}
}

func StopServerProcess() (string, error) {
	log.Printf("executing stop server command")
	if serverCmd == nil || serverStdin == nil {
		return "", fmt.Errorf("server is not running")
	}

	if err := SendCommand("stop"); err != nil {
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
		serverCmd.Process.Kill()
		serverCmd = nil
		serverStdin = nil
		return "server force-killed after timeout", nil
	}
}

func IsServerRunning() bool {
	return serverCmd != nil
}

func SendCommand(cmd string) error {
	stdinMu.Lock()
	defer stdinMu.Unlock()
	if serverStdin == nil {
		return fmt.Errorf("server is not running")
	}
	_, err := serverStdin.Write([]byte(cmd + "\n"))
	return err
}

func GetLogHub() *types.LogHub {
	return logHub
}
