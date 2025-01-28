package downloader

import (
	"downloader/internal/manifest"
	"downloader/pkg/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ProcessManifest(m *manifest.Manifest) error {
	for i, file := range m.Files {
		// Check if the file exists and matches the hash
		localHash, err := manifest.CalculateHash(file.Path)
		if err == nil && localHash == file.Hash {
			fmt.Printf("File %s is up to date.\n", file.Path)
			continue
		}

		// Download the file
		err = downloadFile(file.URL, file.Path, i+1, len(m.Files))
		if err != nil {
			return fmt.Errorf("error downloading file %s: %v", file.Path, err)
		}
		fmt.Printf("File %s downloaded successfully.\n", file.Path)
	}

	return nil
}

func downloadFile(url, filePath string, fileIndex, totalFiles int) error {
	start := time.Now()
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

	// Track download progress
	total := resp.ContentLength
	var downloaded int64

	reader := io.TeeReader(resp.Body, &progressWriter{
		Out:        out,
		Downloaded: &downloaded,
		Total:      total,
		StartTime:  start,
		FileIndex:  fileIndex,
		TotalFiles: totalFiles,
		FileName:   filepath.Base(filePath),
	})

	_, err = io.Copy(io.Discard, reader)
	return err
}

// progressWriter tracks progress and prints it
type progressWriter struct {
	Out        *os.File
	Downloaded *int64
	Total      int64
	StartTime  time.Time
	FileIndex  int
	TotalFiles int
	FileName   string
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.Out.Write(p)
	*pw.Downloaded += int64(n)
	elapsed := time.Since(pw.StartTime)
	mbDownloaded := float64(*pw.Downloaded) / (1024.0 * 1024.0)
	speed := mbDownloaded / elapsed.Seconds()
	fileSizeMB := float64(pw.Total) / (1024.0 * 1024.0)
	util.PrintProgress(util.ProgressInfo{
		Current:    int(*pw.Downloaded),
		Total:      int(pw.Total),
		FileIndex:  pw.FileIndex,
		TotalFiles: pw.TotalFiles,
		Speed:      speed,
		FileSizeMB: fileSizeMB,
		Elapsed:    elapsed,
		FileName:   pw.FileName,
	})
	return n, err
}
