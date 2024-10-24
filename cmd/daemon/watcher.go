package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/OlaHulleberg/Filterbox/common"
	"github.com/OlaHulleberg/Filterbox/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/xattr"
)

type DropboxInfo struct {
	Personal struct {
		Path string `json:"path"`
	} `json:"personal"`
}

func getDropboxPath() (string, error) {
	var dbPath string
	if runtime.GOOS == "windows" {
		dbPath = os.Getenv("LOCALAPPDATA") + "\\Dropbox\\info.json"
	} else {
		dbPath = os.Getenv("HOME") + "/.dropbox/info.json"
	}

	file, err := os.ReadFile(dbPath)
	if err != nil {
		return "", err
	}

	var info DropboxInfo
	if err := json.Unmarshal(file, &info); err != nil {
		return "", err
	}

	return info.Personal.Path, nil
}

var config common.Configuration

func watchConfig() {
	var configFilePath string
	var err error

	// Load initial configuration
	config, configFilePath, err = common.LoadOrCreateConfiguration()
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Error with configuration: %s", err)
		os.Exit(1)
	}

	AppLogger.Println(logger.LevelInfo, "Config loaded:", configFilePath)
	AppLogger.Println(logger.LevelDebug, config)

	configWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Error creating file watcher: %s", err)
		os.Exit(1)
	}
	defer configWatcher.Close()

	configDirPath := filepath.Dir(configFilePath)
	err = configWatcher.Add(configDirPath)
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Error adding file to watcher: %s", err)
		os.Exit(1)
	}

	for {
		select {
		case event, ok := <-configWatcher.Events:
			if !ok {
				return
			}
			// Check if file event soure is our config
			if event.Name == configFilePath {
				// Handle write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					AppLogger.Printf(logger.LevelInfo, "Modified file: %s", event.Name)

					config, _, err = common.LoadOrCreateConfiguration()
					if err != nil {
						AppLogger.Printf(logger.LevelError, "Error reloading configuration: %s", err)
					} else {
						AppLogger.Println(logger.LevelInfo, "Configuration reloaded successfully.")
						AppLogger.Println(logger.LevelDebug, config)
					}
				}
			}

		case err, ok := <-configWatcher.Errors:
			if !ok {
				return
			}
			AppLogger.Printf(logger.LevelError, "Watcher error: %s", err)
		}
	}
}

var watcher *fsnotify.Watcher

func startWatching() {
	dropboxPath, err := getDropboxPath()
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Failed to get Dropbox path: %s", err)
		os.Exit(1)
	}

	AppLogger.Println(logger.LevelInfo, "Dropbox path found:", dropboxPath)

	// Create a config Watcher
	go watchConfig()

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		AppLogger.Println(logger.LevelError, err)
	}
	defer watcher.Close()

	AppLogger.Println(logger.LevelInfo, "Started watching...")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				AppLogger.Println(logger.LevelDebug, "New file event:", event)
				handleEvent(event, config)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				AppLogger.Println(logger.LevelError, "error:", err)
			}
		}
	}()

	err = filepath.Walk(dropboxPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				AppLogger.Printf(logger.LevelError, "Failed to add directory to watcher: %s, error: %v", path, err)
				return err
			} else {
				AppLogger.Printf(logger.LevelDebug, "Adding directory to watcher: %s", path)
			}
		}

		if forceScan {
			if info.IsDir() && shouldIgnore(path, "directory", config) {
				AppLogger.Println(logger.LevelInfo, "Adding ignore xattr to directory and skipping:", path)
				ignoreFile(path)

				// Skip walking through this directory's contents
				return filepath.SkipDir
			} else if !info.IsDir() && shouldIgnore(path, "file", config) {
				AppLogger.Println(logger.LevelInfo, "Adding ignore xattr to file:", path)
				ignoreFile(path)
			}
		}

		return nil
	})

	if err != nil {
		AppLogger.Println(logger.LevelWarn, err) // Path may have been deleted
	}

	<-done
}

func addPathToWatcher(event fsnotify.Event) {
	err := watcher.Add(event.Name)
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Failed to add path to watcher: %s, error: %v", event.Name, err)
	} else {
		AppLogger.Printf(logger.LevelDebug, "Adding path to watcher: %s", event.Name)
	}
}

func handleEvent(event fsnotify.Event, config common.Configuration) {
    if event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
        entityType := "directory"
        if isFile(event.Name) {
            entityType = "file"
        } else {
            addPathToWatcher(event)
        }

        if shouldIgnore(event.Name, entityType, config) {
            AppLogger.Printf(logger.LevelInfo, "Adding ignore xattr to %s: %s", entityType, event.Name)
            ignoreFile(event.Name)
        }
    }
}

func ignoreFile(filePath string) {
	err := xattr.Set(filePath, "com.dropbox.ignored", []byte("1"))
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Failed to set xattr for %s: %v", filepath.Base(filePath), err)
	}
}
