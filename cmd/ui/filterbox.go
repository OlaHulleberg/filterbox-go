package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/OlaHulleberg/Filterbox/common"
	"github.com/OlaHulleberg/Filterbox/logger"
)

var AppLogger *logger.Logger
var config common.Configuration
var configFilePath string
var fyneWindow fyne.Window
var selectedIndex int = -1 // Initially, no item is selected

func main() {
	// Load AppLogger
	var logLevelParameter string
	flag.StringVar(&logLevelParameter, "loglevel", "info", "Set the logging level (error, warn, info, debug)")
	flag.Parse()

	var err error
	AppLogger, err = logger.CreateLogger(logLevelParameter)
	if err != nil {
		log.Printf("Failed to create logger: %s", err)
		os.Exit(1)
	}

	// Start App and read config
	myApp := app.New()
	fyneWindow = myApp.NewWindow("FilterBox Configuration")

	config, configFilePath, err = common.LoadOrCreateConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	filterList := widget.NewList(
		func() int {
			return len(config.Filters)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Filters")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(config.Filters[id].Name)
		},
	)

	filterList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
	}

	addFilterButton := widget.NewButton("Add Filter", func() {
		showAddFilterDialog()
		refreshWindowLayout()
		saveConfig()
	})

	removeFilterButton := widget.NewButton("Remove Selected Filter", func() {
		if selectedIndex >= 0 && selectedIndex < len(config.Filters) {
			config.Filters = append(config.Filters[:selectedIndex], config.Filters[selectedIndex+1:]...)
			selectedIndex = -1                 // Reset the selection
			filterList.Unselect(selectedIndex) // Unselect the item in the list
			filterList.Refresh()
			refreshWindowLayout()
			saveConfig()
		}
	})

	buttonsContainer := container.NewHBox(addFilterButton, removeFilterButton)

	fyneWindow.SetContent(container.NewBorder(nil, buttonsContainer, nil, nil, filterList))

	// Set the window size. Maintain the current width and change only the height
	currentWidth := fyneWindow.Content().Size().Width
	fyneWindow.Resize(fyne.NewSize(currentWidth, 200))

	fyneWindow.ShowAndRun()
}

func refreshWindowLayout() {
	fyneWindow.Resize(fyneWindow.Content().MinSize())
}

func showAddFilterDialog() {
	// Create a new window for the dialog
	dialogWindow := fyne.CurrentApp().NewWindow("Add Filter")

	nameEntry := widget.NewEntry()
	typeEntry := widget.NewSelect([]string{"file", "directory", "both"}, nil)

	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Type", typeEntry),
	)

	form.OnSubmit = func() {
		config.Filters = append(config.Filters, common.Filter{
			Name: nameEntry.Text,
			Type: typeEntry.Selected,
		})
		refreshWindowLayout() // Update layout to reflect new filter
		dialogWindow.Close()
	}
	form.OnCancel = func() {
		dialogWindow.Close()
	}

	// Set the form as the content of the dialog window
	dialogWindow.SetContent(form)
	dialogWindow.Show()
}

func saveConfig() {
	configData, err := json.Marshal(config)
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Failed to marshal config: %s", err)
		return
	}

	err = os.WriteFile(configFilePath, configData, 0644)
	if err != nil {
		AppLogger.Printf(logger.LevelError, "Failed to write config file: %s", err)
	}
}

