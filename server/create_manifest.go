package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
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

func main() {
	filesDir := flag.String("files", "files", "Directory containing the files to process")
	baseURL := flag.String("url", "http://localhost:8080/", "Base URL for file download links")
	version := flag.String("version", "1.0", "Manifest version")
	flag.Parse()

	var manifest Manifest
	manifest.Version = *version

	// Walk through all files in the directory recursively
	err := filepath.WalkDir(*filesDir, func(path string, d fs.DirEntry, err error) error {
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

		hash, err := calculateHash(path)
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", path, err)
			return nil
		}

		patchFile := PatchFile{
			Path:   relPath,
			Hash:   hash,
			Size:   info.Size(),
			Custom: true,
			URL:    *baseURL + relPath,
		}

		manifest.Files = append(manifest.Files, patchFile)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	// Write the manifest to a file
	outputFile := "manifest.json"
	err = writeManifest(manifest, outputFile)
	if err != nil {
		fmt.Printf("Error writing manifest: %v\n", err)
		return
	}

	fmt.Printf("Manifest created successfully: %s\n", outputFile)
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func writeManifest(manifest Manifest, outputFile string) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
