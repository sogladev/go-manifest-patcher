package main

import (
	"downloader/internal/config"
	"downloader/internal/downloader"
	"downloader/internal/logger"
	"downloader/internal/manifest"
	"log"
	"strings"
)

func main() {
	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger()

	// Download and parse manifest
	manifestURL := cfg.ManifestURL
	m, err := manifest.DownloadManifest(manifestURL)
	if err != nil {
		log.Fatalf("Failed to download manifest: %v", err)
	}

	// Verify files and download missing or outdated files
	err = downloader.ProcessManifest(m)
	if err != nil {
		log.Fatalf("Failed to process manifest: %v", err)
	}

	println("\n" + strings.Repeat("-", 80))
	println("All files are up to date or successfully downloaded.")
}
