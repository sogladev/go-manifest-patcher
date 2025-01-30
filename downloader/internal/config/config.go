package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/sogladev/go-manifest-patcher/downloader/internal/filter"
	"github.com/sogladev/go-manifest-patcher/downloader/internal/logger"
)

type Config struct {
	ManifestURL string
	LogLevel    string
	SaveFilter  bool
	SkipUpdate  bool
}

func InitConfig() *Config {
	manifestURL := flag.String("manifest", "https://updater.project-epoch.net/api/manifest?environment=production&internal_key=", "Path to manifest.json file or URL (e.g., http://localhost:8080/manifest.json)")
	logLevel := flag.String("log-level", "info", "Set the log level (debug, info, warning, error)")
	saveFilter := flag.Bool("save-filter", false, "Save the default filter to filter.json and exit")
	skipUpdate := flag.Bool("skip-update", false, "Skip update check (useful for development)")
	flag.Parse()

	if *saveFilter {
		f := filter.DefaultFilter()
		err := filter.SaveFilter("filter.json", f)
		if err != nil {
			logger.Error.Fatalf("Failed to save filter.json: %v", err)
		}
		fmt.Println("Saved default filter to filter.json")
		os.Exit(0)
	}

	return &Config{
		ManifestURL: *manifestURL,
		LogLevel:    *logLevel,
		SaveFilter:  *saveFilter,
		SkipUpdate:  *skipUpdate,
	}
}
