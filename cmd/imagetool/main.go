//go:generate goversioninfo -o resource.syso versioninfo.json

// Package main is the entry point for Image-Tool.
// This is a Windows TUI application for image and PDF processing.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"imagetool/internal/config"
	"imagetool/internal/logging"
	"imagetool/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Set terminal title
	fmt.Print("\033]0;Image-Tool\007")

	// Initialize logging
	logDir := filepath.Join(config.GetConfigDir(), "logs")
	if err := logging.Init(logDir); err != nil {
		// Continue without logging if initialization fails
		fmt.Fprintf(os.Stderr, "Warning: Could not initialize logging: %v\n", err)
	}
	defer logging.Close()

	logging.Info("Application starting", nil)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logging.Warn("Could not load config, using defaults", map[string]interface{}{
			"error": err.Error(),
		})
	}
	_ = cfg // Config available for future use

	// Run TUI
	p := tea.NewProgram(ui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logging.Error("Application error", map[string]interface{}{
			"error": err.Error(),
		})
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	logging.Info("Application exiting normally", nil)
}
