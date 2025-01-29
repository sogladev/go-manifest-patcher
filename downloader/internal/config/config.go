package config

import (
	"flag"
)

type Config struct {
	ManifestURL string
	LogLevel    string
}

func InitConfig() *Config {
	manifestURL := flag.String("manifest", "manifest.json", "Path to manifest.json file or URL (e.g., http://localhost:8080/manifest.json)")
	logLevel := flag.String("log-level", "info", "Set the log level (debug, info, warning, error)")
	flag.Parse()

	flag.Parse()

	return &Config{
		ManifestURL: *manifestURL,
		LogLevel:    *logLevel,
	}
}
