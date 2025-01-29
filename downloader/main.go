package main

import (
	"log"
	"strings"

	"github.com/common-nighthawk/go-figure"

	"github.com/sogladev/golang-terminal-downloader/downloader/internal/config"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/downloader"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
	"github.com/sogladev/golang-terminal-downloader/pkg/manifest"
)

func main() {
	// Print banner
	myFigure := figure.NewFigure("Banner", "slant", true)
	myFigure.Print()
	println("")

	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger()

	// Load manifest from file or URL
	m, err := manifest.LoadManifest(cfg.ManifestURL)
	if err != nil {
		log.Fatalf("Failed to load manifest: %v", err)
	}

	// Verify files and download missing or outdated files
	err = downloader.ProcessManifest(m)
	if err != nil {
		log.Fatalf("Failed to process manifest: %v", err)
	}

	println("\n" + strings.Repeat("-", 80))
	println("All files are up to date or successfully downloaded.")
}
