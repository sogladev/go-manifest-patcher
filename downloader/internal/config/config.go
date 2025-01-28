package config

import (
	"flag"
)

type Config struct {
	ManifestURL string
}

func InitConfig() *Config {
	manifestURL := flag.String("manifest", "http://localhost:8080/manifest.json", "URL to the manifest.json file")
	flag.Parse()

	return &Config{
		ManifestURL: *manifestURL,
	}
}
