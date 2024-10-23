package main

import (
	"testing"

	"github.com/OlaHulleberg/Filterbox/common"
	"github.com/stretchr/testify/assert"
)

func TestShouldIgnore(t *testing.T) {
	config := common.Configuration{
		Filters: []common.Filter{
			{Name: "*.txt", Type: "file"},
			{Name: "logs", Type: "directory"},
			{Name: "common", Type: "both"},
		},
	}

	assert.True(t, shouldIgnore("test.txt", "file", config))
	assert.False(t, shouldIgnore("test.jpg", "file", config))
	assert.True(t, shouldIgnore("logs", "directory", config))
	assert.True(t, shouldIgnore("common", "file", config))
	assert.True(t, shouldIgnore("common", "directory", config))
}
