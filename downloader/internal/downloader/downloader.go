package downloader

import (
	"downloader/internal/manifest"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ProcessManifest(m *manifest.Manifest) error {
	for _, file := range m.Files {
		// Check if the file exists and matches the hash
		localHash, err := manifest.CalculateHash(file.Path)
		if err == nil && localHash == file.Hash {
			fmt.Printf("File %s is up to date.\n", file.Path)
			continue
		}

		// Download the file
		err = downloadFile(file.URL, file.Path)
		if err != nil {
			return fmt.Errorf("error downloading file %s: %v", file.Path, err)
		}
		fmt.Printf("File %s downloaded successfully.\n", file.Path)
	}

	return nil
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	// Create directories if they don't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
