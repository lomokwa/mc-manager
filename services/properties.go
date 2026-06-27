package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lomokwa/mc-manager/utils"
)

// readServerProperties parses an existing server.properties file into a map.
// A missing file yields an empty map (no error) so callers can fall back to
// the defaults.
func readServerProperties(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	props := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if key, value, ok := strings.Cut(line, "="); ok {
			props[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return props, nil
}

// UpdateServerProperties merges the given properties over the current
// server.properties (falling back to the defaults for any keys not present
// yet) and writes the file back. The new values take effect on the next
// server start.
func UpdateServerProperties(properties map[string]string) error {
	path := filepath.Join(ServerDir, "server.properties")

	merged := make(map[string]string, len(DefaultServerProperties))
	for k, v := range DefaultServerProperties {
		merged[k] = v
	}

	existing, err := readServerProperties(path)
	if err != nil {
		return fmt.Errorf("failed to read server.properties: %w", err)
	}
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range properties {
		merged[k] = v
	}

	var content strings.Builder
	for k, v := range merged {
		fmt.Fprintf(&content, "%s=%s\n", k, v)
	}

	if err := utils.WriteFile(path, []byte(content.String())); err != nil {
		return fmt.Errorf("failed to write server.properties: %w", err)
	}
	return nil
}

// GetServerProperties returns the effective server.properties: the defaults
// overlaid with whatever is currently in the file (if any).
func GetServerProperties() (map[string]string, error) {
	merged := make(map[string]string, len(DefaultServerProperties))
	for k, v := range DefaultServerProperties {
		merged[k] = v
	}

	existing, err := readServerProperties(filepath.Join(ServerDir, "server.properties"))
	if err != nil {
		return nil, fmt.Errorf("failed to read server.properties: %w", err)
	}
	for k, v := range existing {
		merged[k] = v
	}
	return merged, nil
}
