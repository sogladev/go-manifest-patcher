package main

import (
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"

	"github.com/sogladev/golang-terminal-downloader/downloader/internal/config"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/downloader"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/filter"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
	"github.com/sogladev/golang-terminal-downloader/pkg/manifest"
	"github.com/sogladev/golang-terminal-downloader/pkg/prompt"
)

func main() {
	// Print banner
	myFigure := figure.NewFigure("Banner", "slant", true)
	myFigure.Print()
	println("")

	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Load manifest from file or URL
	m, err := manifest.LoadManifest(cfg.ManifestURL)
	if err != nil {
		logger.Error.Fatalf("Failed to load manifest: %v", err)
	}

	// Load filter configuration
	f, err := filter.LoadFilter("filter.json")
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug.Println("No custom filter config found, using default filter")
			f = filter.DefaultFilter()
		} else {
			logger.Error.Fatalf("Failed to parse filter.json: %v", err)
		}
	} else {
		println("\nUsing custom filter config")
	}

	// Verify files and download missing or outdated files
	err = downloader.ProcessManifest(m, f)
	if err != nil {
		if err == prompt.ErrUserCancelled {
			os.Exit(0) // Exit gracefully if user cancelled
		} else {
			logger.Error.Fatalf("Failed to process manifest: %v", err)
		}
	}

	println("\n" + strings.Repeat("-", 80))
	println("All files are up to date or successfully downloaded.")
}
