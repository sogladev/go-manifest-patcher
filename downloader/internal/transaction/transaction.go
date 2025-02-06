package transaction

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sogladev/go-manifest-patcher/downloader/internal/logger"
	"github.com/sogladev/go-manifest-patcher/pkg/manifest"
	"github.com/sogladev/go-manifest-patcher/pkg/util"
)

type Status int

const (
	UpToDate Status = iota
	Missing
	OutOfDate
)

type FileOperation struct {
	Path   string
	Size   int64
	Hash   string
	File   *manifest.PatchFile
	Status Status
}

type Transaction struct {
	Operations []*FileOperation
}

func newTransaction() *Transaction {
	return &Transaction{
		Operations: make([]*FileOperation, 0),
	}
}

func CreateTransaction(m *manifest.Manifest) *Transaction {
	transaction := newTransaction()
	for i, file := range m.Files {
		var status Status
		hash, err := manifest.CalculateHashMD5(file.Path)
		if err == nil {
			if hash == file.Hash {
				status = UpToDate
			} else {
				status = OutOfDate
			}
		} else {
			status = Missing
		}
		operation := &FileOperation{
			Path:   file.Path,
			Size:   file.Size,
			Hash:   hash,
			File:   &m.Files[i],
			Status: status,
		}
		transaction.Operations = append(transaction.Operations, operation)
	}
	return transaction
}

func (t *Transaction) Print(m *manifest.Manifest, localFiles map[string]bool) error {
	var totalDownloadSize int64
	var totalDiskChange int64

	// Split operations into categories based on Status
	// Remove Manifest files from localFiles, result are ExtraFiles not in manifest
	filteredOps := map[Status][]*FileOperation{
		UpToDate:  {},
		OutOfDate: {},
		Missing:   {},
	}
	for _, op := range t.Operations {
		filteredOps[op.Status] = append(filteredOps[op.Status], op)
		delete(localFiles, op.Path)
	}

	fmt.Println("\nManifest Overview:")
	fmt.Printf(" Version: %s\n", m.Version)
	fmt.Printf(" %s\n", util.ColorGreen("Up-to-date files:"))
	for _, op := range filteredOps[UpToDate] {
		fmt.Printf("  %s (Size: %s)\n",
			util.ColorGreen(op.File.Path),
			humanize.Bytes(uint64(op.Size)),
		)
	}

	fmt.Printf("\n %s\n", util.ColorYellow("Outdated files (will be updated):"))
	for _, op := range filteredOps[OutOfDate] {
		totalDownloadSize += op.File.Size
		totalDiskChange += op.File.Size - op.Size

		fmt.Printf("  %s (Current Size: %s, New Size: %s)\n",
			util.ColorYellow(op.File.Path),
			humanize.Bytes(uint64(op.File.Size)),
			humanize.Bytes(uint64(op.Size)),
		)
		logger.Debug.Printf("File: %s, Current Hash: %s, New Hash: %s", op.File.URL, op.Hash, op.File.Hash)
	}

	fmt.Printf("\n %s\n", util.ColorRed("Missing files (will be downloaded):"))
	for _, op := range filteredOps[Missing] {
		totalDownloadSize += op.File.Size
		totalDiskChange += op.File.Size
		fmt.Printf("  %s (New Size: %s)\n",
			util.ColorRed(op.File.Path),
			humanize.Bytes(uint64(op.Size)),
		)
	}

	fmt.Printf("\n %s\n", util.ColorCyan("Extra files (not in manifest):"))
	extraFilesCount := 0
	for file := range localFiles {
		if extraFilesCount < 10 {
			info, _ := os.Stat(file)
			fmt.Printf("  %s (Size: %s)\n",
				util.ColorCyan(file),
				humanize.Bytes(uint64(info.Size())),
			)
		}
		extraFilesCount++
	}
	if extraFilesCount > 10 {
		fmt.Printf("  ...and %d more files\n", extraFilesCount-10)
	}

	if len(t.Operations) > 0 {
		fmt.Printf("\nTransaction Summary:\n")
		fmt.Printf(" Installing/Updating: %d files\n\n", len(t.Operations))

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
	}
	return nil
}

func (t *Transaction) Download(m *manifest.Manifest, localFiles map[string]bool) error {
	totalFiles := len(t.Operations)
	currentFile := 0

	for _, op := range t.Operations {
		if op.Status == Missing || op.Status == OutOfDate {
			currentFile++
			err := downloadFile(op.File.URL, op.Path, currentFile, totalFiles)
			if err != nil {
				return fmt.Errorf("error downloading file %s: %v", op.Path, err)
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

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

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
