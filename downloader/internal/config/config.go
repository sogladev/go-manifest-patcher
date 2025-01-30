package config

import (
	"flag"
	"os"

	"github.com/sogladev/go-manifest-patcher/downloader/internal/filter"
)

type Config struct {
	ManifestURL string
	LogLevel    string
	SaveFilter  bool
	SkipUpdate  bool
}

func InitConfig() *Config {
	manifestURL := flag.String("manifest", "manifest.json", "Path to manifest.json file or URL (e.g., http://localhost:8080/manifest.json)")
	logLevel := flag.String("log-level", "info", "Set the log level (debug, info, warning, error)")
	saveFilter := flag.Bool("save-filter", false, "Save the default filter to filter.json and exit")
	skipUpdate := flag.Bool("skip-update", false, "Skip update check (useful for development)")
	flag.Parse()

	if *saveFilter {
		f := filter.DefaultFilter()
		filter.SaveFilter("filter.json", f)
		os.Exit(0)
	}

	return &Config{
		ManifestURL: *manifestURL,
		LogLevel:    *logLevel,
		SaveFilter:  *saveFilter,
		SkipUpdate:  *skipUpdate,
	}
}
