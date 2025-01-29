package config

import (
	"flag"
)

type Config struct {
	ManifestURL string
}

func InitConfig() *Config {
	manifestURL := flag.String("manifest", "manifest.json", "Path to manifest.json file or URL (e.g., http://localhost:8080/manifest.json)")
	flag.Parse()

	return &Config{
		ManifestURL: *manifestURL,
	}
}
