package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/utils"
)

func main() {
	r := gin.Default()

	r.POST("/start", func(c *gin.Context) {
		if !utils.FileExists("./minecraft-server/server.jar") || !utils.FileExists("./minecraft-server") {
			resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch version manifest",
				})
				return
			}
			defer resp.Body.Close()

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

			if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to decode version manifest",
				})
				return
			}

			latestId := manifest.Latest.Release

			var versionURL string
			for _, version := range manifest.Versions {
				if version.ID == latestId {
					versionURL = version.URL
					break
				}
			}

			if versionURL == "" {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Latest version URL not found",
				})
				return
			}

			res2, err := http.Get(versionURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch latest version details",
				})
				return
			}
			defer res2.Body.Close()

			var versionDetails struct {
				Downloads struct {
					Server struct {
						URL string `json:"url"`
					} `json:"server"`
				} `json:"downloads"`
			}
			if err := json.NewDecoder(res2.Body).Decode(&versionDetails); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to decode version details",
				})
				return
			}

			serverJarURL := versionDetails.Downloads.Server.URL

			fmt.Println("Downloading server.jar from:", serverJarURL)

			err = utils.DownloadFile(serverJarURL, "./minecraft-server/server.jar")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to download server.jar: " + err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Server jar downloaded successfully",
			})
		}

		cmd := exec.Command("echo", "Starting server...")

		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"output": string(output),
		})
	})

	r.Run()
}
