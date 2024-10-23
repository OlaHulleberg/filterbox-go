package main

import (
	"os"
	"path/filepath"

	"github.com/OlaHulleberg/Filterbox/common"
)

func shouldIgnore(path string, itemType string, config common.Configuration) bool {
	for _, filter := range config.Filters {
		if filter.Type == itemType || filter.Type == "both" {
			matched, err := filepath.Match(filter.Name, filepath.Base(path))
			if err != nil {
				// TODO: Handle the error (e.g., log it, ignore this filter, etc.)
				continue
			}
			if matched {
				return true
			}
		}
	}
	return false
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
