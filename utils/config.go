package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func getConfigFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".config", "shipo-cli", "config"), nil
}

func LoadConfig() (map[string]string, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("config file not found; please create ~/.config/shipo-cli/config with your handle, password, and limit")
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}
