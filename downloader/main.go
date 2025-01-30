package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/common-nighthawk/go-figure"

	"github.com/sogladev/golang-terminal-downloader/downloader/internal/config"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/downloader"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/filter"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
	"github.com/sogladev/golang-terminal-downloader/downloader/updater"
	"github.com/sogladev/golang-terminal-downloader/pkg/manifest"
	"github.com/sogladev/golang-terminal-downloader/pkg/prompt"
)

const currentVersion = "v0.0.1"

func main() {
	// Print banner
	myFigure := figure.NewFigure("Banner", "slant", true)
	myFigure.Print()
	println("")

	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Check for updates
	newVersion, downloadURL, err := updater.CheckForUpdate(currentVersion)
	if err != nil {
		log.Fatalf("Failed to check for updates: %v", err)
	}
	if newVersion != "" {
		fmt.Printf("New version available: %s\n", newVersion)
		err = prompt.PromptyN("Do you want to update? [y/N]: ")
		if err == nil {
			tempFile := "new-" + updater.GetExecutableName()
			err := updater.Download(downloadURL, tempFile)
			if err != nil {
				log.Fatalf("Failed to download update: %v", err)
			}
			if runtime.GOOS == "windows" {
				// On Windows we cannot replace the running executable, so we need to inform the user
				fmt.Printf("Update downloaded as: %v.\nPlease rename the new executable to %v and restart the application!\n", tempFile, updater.GetExecutableName())
			} else {
				err = updater.ReplaceExecutable(tempFile)
				if err != nil {
					log.Fatalf("Failed to replace executable: %v", err)
				}
				fmt.Println("Update successful. Please restart the application.")
			}
			return
		}
	}

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
			return
		} else {
			logger.Error.Fatalf("Failed to process manifest: %v", err)
		}
	}

	println("\n" + strings.Repeat("-", 80))
	println("All files are up to date or successfully downloaded.")
}
