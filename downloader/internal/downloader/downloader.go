package downloader

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sogladev/golang-terminal-downloader/downloader/internal/filter"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
	"github.com/sogladev/golang-terminal-downloader/pkg/manifest"
	"github.com/sogladev/golang-terminal-downloader/pkg/prompt"
	"github.com/sogladev/golang-terminal-downloader/pkg/util"

	"github.com/dustin/go-humanize"
)

type FileOperation struct {
	Path      string
	Size      int64
	NewSize   int64
	LocalHash string
	NewHash   string
}

func findOperationIndex(operations []FileOperation, path string) int {
	for i, op := range operations {
		if op.Path == path {
			return i
		}
	}
	return -1
}

func ProcessManifest(m *manifest.Manifest, f *filter.Filter) error {

	// Track files in the local folder for "extra files" detection
	localFiles := map[string]bool{}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if !f.IsIgnored(path) {
				localFiles[path] = true
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading local files: %v", err)
	}

	// Track categorized files
	upToDate := []string{}
	missing := []string{}
	outdated := []string{}

	// Track file operations and sizes
	operations := make([]FileOperation, 0)
	var totalDownloadSize int64
	var totalDiskChange int64

	// Process each file in the manifest
	for _, file := range m.Files {
		op := FileOperation{Path: file.Path}

		// Get remote file size
		resp, err := http.Head(file.URL)
		if err == nil {
			op.NewSize = resp.ContentLength
		}

		// Check if the file exists locally
		if info, err := os.Stat(file.Path); err == nil {
			op.Size = info.Size()
		}

		localHash, err := manifest.CalculateHashMD5(file.Path)
		if err == nil && localHash == file.Hash {
			upToDate = append(upToDate, file.Path)
		} else if err == nil && localHash != file.Hash {
			op.LocalHash = localHash
			op.NewHash = file.Hash
			outdated = append(outdated, file.Path)
			totalDownloadSize += op.NewSize
			totalDiskChange += op.NewSize - op.Size
			operations = append(operations, op)
		} else if os.IsNotExist(err) {
			missing = append(missing, file.Path)
			totalDownloadSize += op.NewSize
			totalDiskChange += op.NewSize
			operations = append(operations, op)
		}

		delete(localFiles, file.Path)
	}

	// Display categorized files
	fmt.Println("\nManifest Overview:")
	fmt.Printf(" Version: %s\n", m.Version)
	fmt.Printf(" %s\n", util.ColorGreen("Up-to-date files:"))
	for _, file := range upToDate {
		info, _ := os.Stat(file)
		fmt.Printf("  %s (Size: %s)\n",
			util.ColorGreen(file),
			humanize.Bytes(uint64(info.Size())),
		)
	}

	fmt.Printf("\n %s\n", util.ColorYellow("Outdated files (will be updated):"))
	for _, file := range outdated {
		info, _ := os.Stat(file)
		op := operations[findOperationIndex(operations, file)]
		fmt.Printf("  %s (Current Size: %s, New Size: %s)\n",
			util.ColorYellow(file),
			humanize.Bytes(uint64(info.Size())),
			humanize.Bytes(uint64(op.NewSize)),
		)
		logger.Debug.Printf("File: %s, Current Hash: %s, New Hash: %s", file, op.LocalHash, op.NewHash)
	}

	fmt.Printf("\n %s\n", util.ColorRed("Missing files (will be downloaded):"))
	for _, file := range missing {
		newSize := operations[findOperationIndex(operations, file)].NewSize
		fmt.Printf("  %s (New Size: %s)\n",
			util.ColorRed(file),
			humanize.Bytes(uint64(newSize)),
		)
	}

	fmt.Printf("\n %s\n", util.ColorCyan("Extra files (not in manifest):"))
	for file := range localFiles {
		info, _ := os.Stat(file)
		fmt.Printf("  %s (Size: %s)\n",
			util.ColorCyan(file),
			humanize.Bytes(uint64(info.Size())),
		)
	}

	// Display transaction summary
	if len(operations) > 0 {
		fmt.Printf("\nTransaction Summary:\n")
		fmt.Printf(" Installing/Updating: %d files\n\n", len(operations))

		fmt.Printf("Total size of inbound files is %s. Need to download %s.\n",
			humanize.Bytes(uint64(totalDownloadSize)),
			humanize.Bytes(uint64(totalDownloadSize)))

		if totalDiskChange > 0 {
			fmt.Printf("After this operation, %s of additional disk space will be used.\n",
				humanize.Bytes(uint64(totalDiskChange)))
		} else {
			fmt.Printf("After this operation, %s of disk space will be freed.\n",
				humanize.Bytes(uint64(-totalDiskChange)))
		}

		err = prompt.PromptyN("Is this ok [y/N]: ")
		if err != nil {
			return err
		}
	}

	// Handle downloads for missing and outdated files
	totalFiles := len(missing) + len(outdated)
	currentFile := 0

	for _, file := range append(missing, outdated...) {
		currentFile++
		// Find the corresponding file in the manifest
		for _, mf := range m.Files {
			if mf.Path == file {
				err := downloadFile(mf.URL, mf.Path, currentFile, totalFiles)
				if err != nil {
					return fmt.Errorf("error downloading file %s: %v", mf.Path, err)
				}
			}
		}
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
	speed := float64(*pw.Downloaded) / elapsed.Seconds()

	util.PrintProgress(util.ProgressInfo{
		Current:    int(*pw.Downloaded),
		Total:      int(pw.Total),
		FileIndex:  pw.FileIndex,
		TotalFiles: pw.TotalFiles,
		Speed:      speed,
		FileSize:   pw.Total,
		Elapsed:    elapsed,
		FileName:   pw.FileName,
	})
	return n, err
}
