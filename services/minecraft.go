package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lomokwa/mc-manager/types"
	"github.com/lomokwa/mc-manager/utils"
)

func DownloadLatestServerJar(destPath string) error {
	log.Printf("downloading version manifest")
	res, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return fmt.Errorf("failed to fetch version manifest")
	}

	defer res.Body.Close()

	var manifest struct {
		Latest struct {
			Release  string `json:"release"`
			Snapshot string `json:"snapshot"`
		} `json:"latest"`
		Versions []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			URL  string `json:"url"`
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&manifest); err != nil {
		return fmt.Errorf("failed to decode version manifest")
	}

	latestId := manifest.Latest.Release

	var versionUrl string
	for _, version := range manifest.Versions {
		if version.ID == latestId {
			versionUrl = version.URL
			break
		}
	}

	if versionUrl == "" {
		return fmt.Errorf("latest version URL not found")
	}

	log.Printf("downloading latest version details")
	versionRes, err := http.Get(versionUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch latest version details")
	}
	defer versionRes.Body.Close()

	var versionDetails struct {
		Downloads struct {
			Server struct {
				URL string `json:"url"`
			} `json:"server"`
		} `json:"downloads"`
	}

	if err := json.NewDecoder(versionRes.Body).Decode(&versionDetails); err != nil {
		return fmt.Errorf("failed to decode version details")
	}

	serverJarUrl := versionDetails.Downloads.Server.URL

	log.Printf("downloading server jar to %s", destPath)
	err = utils.DownloadFile(serverJarUrl, "./minecraft-server/server.jar")
	if err != nil {
		return fmt.Errorf("failed to download server.jar: %s", err)
	}

	log.Printf("server jar download complete")

	return nil
}

func PrepareServerFiles(serverDir string, createLaunchScript bool, configureProperties bool, requestProperties map[string]string) error {
	log.Printf("preparing server files in %s", serverDir)

	// Reject control characters in caller-supplied properties before writing
	// anything: a newline in a key or value would inject extra
	// server.properties lines (e.g. enabling RCON or setting a resource pack).
	for k, v := range requestProperties {
		if strings.ContainsAny(k, "\r\n") || strings.ContainsAny(v, "\r\n") {
			return fmt.Errorf("invalid server property %q: keys and values must not contain newlines", k)
		}
	}

	if err := utils.WriteFile(filepath.Join(serverDir, "eula.txt"), []byte("eula=true")); err != nil {
		return err
	}

	// Create server.properties file content.
	properties := make(map[string]string, len(DefaultServerProperties))
	for k, v := range DefaultServerProperties {
		properties[k] = v
	}

	for k, v := range requestProperties {
		properties[k] = v
	}

	var content strings.Builder
	for k, v := range properties {
		fmt.Fprintf(&content, "%s=%s\n", k, v)
	}

	propertiesContent := []byte(content.String())
	if configureProperties {
		log.Printf("writing server.properties")
		if err := utils.WriteFile(filepath.Join(serverDir, "server.properties"), propertiesContent); err != nil {
			return err
		}
	}

	if createLaunchScript {
		log.Printf("writing launch scripts")
		shellScriptPath := filepath.Join(serverDir, "start-server.sh")
		batScriptPath := filepath.Join(serverDir, "start-server.bat")

		if err := utils.WriteFile(shellScriptPath, []byte(DefaultStartServerShellScript)); err != nil {
			return fmt.Errorf("failed to write start-server.sh: %w", err)
		}

		if err := os.Chmod(shellScriptPath, 0755); err != nil {
			return fmt.Errorf("failed to set executable permission on start-server.sh: %w", err)
		}

		if err := utils.WriteFile(batScriptPath, []byte(DefaultStartServerBatchScript)); err != nil {
			return fmt.Errorf("failed to write start-server.bat: %w", err)
		}
	}

	log.Printf("server file preparation complete")

	return nil
}

func loadUUIDs(filename string) (map[string]bool, error) {
	data, err := os.ReadFile(filepath.Join(ServerDir, filename))
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, err
	}

	var entries []struct {
		UUID string `json:"uuid"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", filename, err)
	}

	set := make(map[string]bool, len(entries))
	for _, e := range entries {
		set[e.UUID] = true
	}

	return set, nil
}

func GetOnlinePlayers() ([]string, error) {
	hub := GetLogHub()
	if hub == nil {
		return nil, fmt.Errorf("log hub not available")
	}

	ch := hub.Subscribe()
	defer hub.Unsubscribe(ch)

draining:
	for {
		select {
		case <-ch:
		default:
			break draining
		}
	}

	if err := SendCommand("list"); err != nil {
		return nil, err
	}

	for {
		select {
		case line := <-ch:
			if strings.Contains(line, "players online:") {
				parts := strings.SplitN(line, "players online: ", 2)

				if len(parts) < 2 || parts[1] == "" {
					return []string{}, nil
				}

				names := strings.Split(parts[1], ", ")
				for i := range names {
					names[i] = strings.TrimSpace(names[i])
				}

				return names, nil
			}

		case <-time.After(5 * time.Second):
			return nil, fmt.Errorf("timed out waiting for player list")
		}
	}
}

func ListPlayers() ([]types.Player, error) {
	data, err := os.ReadFile(filepath.Join(ServerDir, "usercache.json"))
	if err != nil {
		return nil, err
	}

	var userCache []types.UserCacheEntry
	if err := json.Unmarshal(data, &userCache); err != nil {
		return nil, fmt.Errorf("failed to decode usercache.json: %w", err)
	}

	// Load status set
	opSet, err := loadUUIDs("ops.json")
	if err != nil {
		return nil, err
	}

	whitelistSet, err := loadUUIDs("whitelist.json")
	if err != nil {
		return nil, err
	}

	bannedSet, err := loadUUIDs("banned-players.json")
	if err != nil {
		return nil, err
	}

	// Get online players
	onlineSet := make(map[string]bool)
	if IsServerRunning() {
		names, err := GetOnlinePlayers()
		if err != nil {
			log.Printf("could not find online players")
			for _, n := range names {
				onlineSet[n] = false
			}
		} else {
			for _, name := range names {
				onlineSet[name] = true
			}
		}
	}

	players := make([]types.Player, 0, len(userCache))
	for _, u := range userCache {
		players = append(players, types.Player{
			UUID:          u.UUID,
			Name:          u.Name,
			Online:        onlineSet[u.Name],
			IsOp:          opSet[u.UUID],
			IsBanned:      bannedSet[u.UUID],
			IsWhitelisted: whitelistSet[u.UUID],
		})
	}
	return players, nil
}
