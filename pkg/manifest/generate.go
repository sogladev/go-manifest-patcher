package manifest

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GenerateManifest(filesDir, baseURL, version string) error {
	var m Manifest
	m.Version = version

	// Walk through all files in the directory recursively
	err := filepath.WalkDir(filesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // Skip directories
		}

		// Get path relative to working directory instead of files directory
		relPath, err := filepath.Rel(".", path)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", path, err)
			return nil
		}

		// Replace Windows backslashes with forward slashes
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", path, err)
			return nil
		}

		hash, err := CalculateHash(path)
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", path, err)
			return nil
		}

		patchFile := PatchFile{
			Path:   relPath,
			Hash:   hash,
			Size:   info.Size(),
			Custom: true,
			URL:    baseURL + relPath,
		}

		m.Files = append(m.Files, patchFile)
		return nil
	})

	if err != nil {
		return err
	}

	// Write the manifest to a file
	outputFile := "manifest.json"
	err = writeManifest(m, outputFile)
	return err
}

func writeManifest(manifest Manifest, outputFile string) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
