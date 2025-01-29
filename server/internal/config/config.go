package config

import (
	"flag"
)

type Config struct {
	Interval       int
	CreateManifest bool
	FilesDir       string
	BaseURL        string
	Version        string
}

func InitConfig() *Config {
	// Add command-line flag for throttle interval
	interval := flag.Int("interval", 1, "ms delay per chunk")
	// Generate a manifest file for the input directory
	createManifest := flag.Bool("create-manifest", false, "Generate manifest.json before starting the server")
	filesDir := flag.String("files", "files", "Directory containing the files to process")
	baseURL := flag.String("url", "http://localhost:8080/", "Base URL for file download links")
	version := flag.String("version", "1.0", "Manifest version")

	flag.Parse()

	return &Config{
		Interval:       *interval,
		CreateManifest: *createManifest,
		FilesDir:       *filesDir,
		BaseURL:        *baseURL,
		Version:        *version,
	}
}
