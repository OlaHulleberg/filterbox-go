package main

import (
	"flag"
	"log"
	"os"

	"github.com/OlaHulleberg/Filterbox/logger"
	"github.com/kardianos/service"
)

var AppLogger *logger.Logger
var forceScan = false

type program struct{}

func (p *program) Start(s service.Service) error {
    // Start should not block. Do the actual work async.
    go p.run()
    return nil
}

func (p *program) run() {
    startWatching()
}

func (p *program) Stop(s service.Service) error {
    // Stop should not block. Return with a few seconds.
    return nil
}

func main() {
    var logLevelParameter string
    flag.StringVar(&logLevelParameter, "loglevel", "info", "Set the logging level (error, warn, info, debug)")
    flag.BoolVar(&forceScan, "forcescan", false, "Force a recursive scan of the Dropbox folder, and reapply tags.")
    flag.Parse()

    var err error
    AppLogger, err = logger.CreateLogger(logLevelParameter)
    if (err != nil) {
        log.Printf("Failed to create logger: %s", err)
        os.Exit(1)
    }

    svcConfig := &service.Config{
        Name:        "FilterboxDaemon",
        DisplayName: "Filterbox Daemon",
        Description: "A daemon for watching and tagging Dropbox folders.",
		EnvVars: map[string]string {
			"LOCALAPPDATA": os.Getenv("LOCALAPPDATA"),
			"HOME": os.Getenv("HOME"),
		},
    }

    prg := &program{}
    s, err := service.New(prg, svcConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Handle service commands
    if len(os.Args) > 1 {
        err := service.Control(s, os.Args[1])
        if err != nil {
            log.Printf("Valid actions: %q\n", service.ControlAction)
            log.Fatal(err)
        }
        return
    }

    if err := s.Run(); err != nil {
        AppLogger.Println(logger.LevelError, err)
    }
}