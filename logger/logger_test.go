package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupLogLevel(t *testing.T) {
    {
        logLevel, err := StringToLogLevel("debug")
        if (assert.NoError(t, err)) {
            assert.Equal(t, LevelDebug, logLevel, "log level should be set to debug")
        }
    }

    {
        logLevel, err := StringToLogLevel("info")
        if (assert.NoError(t, err)) {
            assert.Equal(t, LevelInfo, logLevel, "log level should be set to info")
        }
    }

    {
        logLevel, err := StringToLogLevel("warn")
        if (assert.NoError(t, err)) {
            assert.Equal(t, LevelWarn, logLevel, "log level should be set to warn")
        }
    }

    {
        logLevel, err := StringToLogLevel("error")
        if (assert.NoError(t, err)) {
            assert.Equal(t, LevelError, logLevel, "log level should be set to error")
        }
    }
}
