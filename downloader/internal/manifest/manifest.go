package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func DownloadManifest(url string) (*Manifest, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching manifest: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %v", err)
	}

	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error parsing manifest: %v", err)
	}

	// Convert Windows-style paths to cross-platform paths
	for i := range manifest.Files {
		manifest.Files[i].Path = strings.ReplaceAll(manifest.Files[i].Path, "\\", "/")
	}

	return &manifest, nil
}
