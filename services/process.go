package services

import (
	"log"
	"os/exec"
)

// TODO: Replace placeholder echo with actual java command: java -Xms1G -Xmx2G -jar server.jar nogui
// TODO: Set cmd.Dir to the minecraft-server directory so java finds server.jar
// TODO: Run process async (cmd.Start) instead of blocking (CombinedOutput) - Minecraft is long-running
// TODO: Store process handle/PID to support stop and prevent duplicate starts
// TODO: Pipe stdin so we can send "stop" command for graceful shutdown
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
