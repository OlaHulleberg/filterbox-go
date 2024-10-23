package common

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var configName = "filters.json"

type Filter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Configuration struct {
	Filters []Filter `json:"filters"`
}

type configDirProvider func() (string, error)

var ensureConfigDirExists configDirProvider = func() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Couldn't determine user home directory, for filter config.")
	}

	configDirPath := filepath.Join(homeDir, ".local", "share", "filterbox")
	// Handle folder creation
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755) // TODO: Revisit this
		if err != nil {
			return "", errors.New("Error creating directory")
		}
	}

	return configDirPath, nil
}

func LoadOrCreateConfiguration() (Configuration, string, error) {
	var config Configuration

	configDirPath, err := ensureConfigDirExists()
	if err != nil {
		return Configuration{}, "", err
	}

	filePath := filepath.Join(configDirPath, configName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		config = Configuration{
			Filters: []Filter{
				{
					Name: "node_modules",
					Type: "directory",
				},
			},
		}
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return Configuration{}, "", err
		}
		err = os.WriteFile(filePath, data, 0644)
    if err != nil {
		  return Configuration{}, "", err
    }
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Configuration{}, "", err
	}
	err = json.Unmarshal(data, &config)
	return config, filePath, err
}
