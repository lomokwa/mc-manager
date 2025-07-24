package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lomokwa/mc-manager/utils"
)

func DownloadLatestServerJar(destPath string) error {
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

	err = utils.DownloadFile(serverJarUrl, "./minecraft-server/server.jar")
	if err != nil {
		return fmt.Errorf("failed to download server.jar: %s", err)
	}

	return nil
}
