package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PatchFile struct {
	Path   string `json:"Path"`
	Hash   string `json:"Hash"`
	Size   int64  `json:"Size"`
	Custom bool   `json:"Custom"`
	URL    string `json:"URL"`
}

type Manifest struct {
	Version string      `json:"Version"`
	Files   []PatchFile `json:"Files"`
}

func LoadManifest(source string) (*Manifest, error) {
	var data []byte
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		fmt.Printf("Downloading manifest from: %s\n", source)
		data, err = downloadManifestData(source)
	} else {
		fmt.Printf("Loading manifest from local file: %s\n", source)
		data, err = os.ReadFile(source)
	}

	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error parsing manifest: %v", err)
	}

	// Convert Windows-style paths to cross-platform paths
	for i := range manifest.Files {
		manifest.Files[i].Path = filepath.ToSlash(manifest.Files[i].Path)
	}

	return &manifest, nil
}

func downloadManifestData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching manifest: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest, status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
