package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	// Directory containing the files
	filesDir := "files"
	baseURL := "http://localhost:8080/" // Base URL for file download links

	// Read all files in the directory
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	var manifest Manifest
	manifest.Version = "1.0"

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		filePath := filepath.Join(filesDir, file.Name())
		hash, err := calculateHash(filePath)
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", file.Name(), err)
			continue
		}

		patchFile := PatchFile{
			Path:   filePath,
			Hash:   hash,
			Size:   file.Size(),
			Custom: true,
			URL:    baseURL + filePath,
		}

		manifest.Files = append(manifest.Files, patchFile)
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

	return ioutil.WriteFile(outputFile, data, 0644)
}
