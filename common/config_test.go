package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadOrCreateConfiguration(t *testing.T) {
    tmpDir := t.TempDir()
    configFile := filepath.Join(tmpDir, "filters.json")

    originalEnsureConfigDirExists := ensureConfigDirExists
    ensureConfigDirExists = func() (string, error) {
      return tmpDir, nil
    }
    defer func() {
      ensureConfigDirExists = originalEnsureConfigDirExists
    }()

    config, configPath, err := LoadOrCreateConfiguration()
    assert.NoError(t, err, "should not error when creating new configuration")
    assert.Len(t, config.Filters, 1, "should have one default filter")
    assert.Equal(t, "node_modules", config.Filters[0].Name, "default filter should be node_modules")
    assert.Equal(t, configFile, configPath, "Loaded config path should point to the control config.")

    // Test loading existing configuration
    newConfig := Configuration{
        Filters: []Filter{{Name: "test", Type: "file"}},
    }
    data, _ := json.MarshalIndent(newConfig, "", "  ")
    _ = os.WriteFile(configFile, data, 0644)

    loadedConfig, _, err := LoadOrCreateConfiguration()
    assert.NoError(t, err, "should not error when loading existing configuration")
    assert.Equal(t, newConfig, loadedConfig, "loaded configuration should match saved")
}
