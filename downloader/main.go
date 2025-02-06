package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"

	"github.com/sogladev/go-manifest-patcher/downloader/internal/config"
	"github.com/sogladev/go-manifest-patcher/downloader/internal/filter"
	"github.com/sogladev/go-manifest-patcher/downloader/internal/logger"
	"github.com/sogladev/go-manifest-patcher/downloader/internal/transaction"
	"github.com/sogladev/go-manifest-patcher/downloader/launch"
	"github.com/sogladev/go-manifest-patcher/downloader/updater"
	"github.com/sogladev/go-manifest-patcher/pkg/manifest"
	"github.com/sogladev/go-manifest-patcher/pkg/prompt"
)

const currentVersion = "v1.0.1-epoch"

func main() {
	// Print banner
	myFigure := figure.NewFigure("Project Epoch", "slant", true)
	myFigure.Print()
	println("unofficial patch download utility - Sogladev")
	println("Bugs or issues: https://github.com/sogladev/go-manifest-patcher/")
	println(strings.Repeat("-", 96))
	println("")

	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Check for updates
	if cfg.SkipUpdate {
		logger.Debug.Println("Skipping update check as per configuration")
	} else {
		var release *updater.LatestRelease
		var err error
		release, err = updater.Fetch(currentVersion)
		if err != nil {
			logger.Debug.Printf("Failed to check for updates: %v", err)
		} else if release == nil {
			logger.Debug.Println("No new version available")
		} else {
			fmt.Printf("Current version : %s\n", currentVersion)
			fmt.Printf("New version available: %s\n", release.Version)
			if err := prompt.PromptyN("Do you want to update? [y/N]: "); err == nil {
				if err := release.Download(); err != nil {
					logger.Warning.Printf("Failed to update: %v", err)
				}
			}
		}
	}

	if err := run(cfg); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func run(cfg *config.Config) error {
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

	// Load local files
	localFiles, err := filter.CollectExtraFiles(f)
	if err != nil {
		return fmt.Errorf("error reading local files: %v", err)
	}

	// Create transaction and prompt user
	transaction := transaction.CreateTransaction(m)
	if err := transaction.Print(m, localFiles); err != nil {
		return err
	}
	if err := prompt.PromptyN("Is this ok [y/N]: "); err != nil {
		return err
	}

	// Verify files and download missing or outdated files
	if err := transaction.Download(m, localFiles); err != nil {
		logger.Error.Fatalf("Failed to process manifest: %v", err)
	}

	println("\n" + strings.Repeat("-", 96))
	println("All files are up to date or successfully downloaded.")

	// Launch the game client
	err = prompt.PromptyN("\nLaunch WoW Client? [y/N]: ")
	if err == prompt.ErrUserCancelled {
		return nil
	}

	println("")
	launch.LaunchGameClient()
	return nil
}
