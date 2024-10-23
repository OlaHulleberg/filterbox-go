package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
    LevelError = iota
    LevelWarn
    LevelInfo
    LevelDebug
)

type Logger struct {
    level int
}

func (l *Logger) SetLevel(level int) {
    l.level = level
}

func (l *Logger) output(level int, v ...interface{}) {
    if level <= l.level {
        _, file, line, ok := runtime.Caller(2) // 2 levels up in call stack
        if !ok {
            file = "???"
            line = 0
        }
        logPrefix := fmt.Sprintf("[%s] %s:%d: ", LevelToString(level), filepath.Base(file), line)
        log.Output(3, logPrefix+fmt.Sprintln(v...)) // 3 to offset the call to Output
    }
}

func (l *Logger) Println(level int, v ...interface{}) {
    format := strings.Repeat("%v ", len(v))
    format = strings.TrimRight(format, " ")
    l.output(level, fmt.Sprintf(format, v...))
}

func (l *Logger) Printf(level int, format string, v ...interface{}) {
    l.output(level, fmt.Sprintf(format, v...))
}

func LevelToString(level int) string {
    switch level {
    case LevelError:
        return "ERROR"
    case LevelWarn:
        return "WARN"
    case LevelInfo:
        return "INFO"
    case LevelDebug:
        return "DEBUG"
    default:
        return "UNKNOWN"
    }
}

func StringToLogLevel(logLevel string) (int, error) {
    switch strings.ToLower(logLevel) {
    case "error":
        return LevelError, nil
    case "warn":
        return LevelWarn, nil
    case "info":
        return LevelInfo, nil
    case "debug":
        return LevelDebug, nil
    }

    return 0, fmt.Errorf("invalid log level: %s", logLevel)
}

func init() {
    log.SetFlags(0) // Turn off flags since we are handling it
}

func CreateLogger(logLevelParameter string, fileName string) (*Logger, error) {
    logLevel, err := StringToLogLevel(logLevelParameter)
    if err != nil {
        return nil, err
    }

    var logDir string
    if runtime.GOOS == "windows" {
        localAppData := os.Getenv("LOCALAPPDATA")
        if localAppData == "" {
            return nil, fmt.Errorf("failed to get LOCALAPPDATA environment variable")
        }
        logDir = filepath.Join(localAppData, "FilterBox")
    } else {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil, fmt.Errorf("failed to get user home directory: %v", err)
        }
        logDir = filepath.Join(homeDir, ".local", "share", "FilterBox")
    }

    // Create the log directory if it doesn't exist
    if _, err := os.Stat(logDir); os.IsNotExist(err) {
        err = os.MkdirAll(logDir, 0755)
        if err != nil {
            return nil, fmt.Errorf("failed to create log directory: %v", err)
        }
    }

    logFilePath := filepath.Join(logDir, fileName)

    logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %v", err)
    }

    log.SetOutput(logFile)

    appLogger := &Logger{logLevel}
    appLogger.Println(LevelDebug, "Logger initialized")

    return appLogger, nil
}