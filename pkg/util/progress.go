package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	maxFileNameLength = 25 // reduced from 40 to keep line length under 80
	progressBarWidth  = 20
	totalLineWidth    = 80
)

var (
	completedFiles = make(map[int]bool)
)

type ProgressInfo struct {
	Current    int
	Total      int
	FileIndex  int
	TotalFiles int
	Speed      float64 // bytes per second
	FileSize   int64   // total bytes
	Elapsed    time.Duration
	FileName   string
}

func truncateFileName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name + strings.Repeat(" ", maxLen-len(name))
	}
	return name[:maxLen-3] + "..." + strings.Repeat(" ", 0)
}

func createProgressBar(current, total, width int) string {
	progress := float64(current) / float64(total)
	filled := int(progress * float64(width))
	return "[" + strings.Repeat("-", filled) + strings.Repeat(" ", width-filled) + "]"
}

func PrintProgress(info ProgressInfo) {
	// Format common elements
	percent := (float64(info.Current) / float64(info.Total)) * 100
	progressBar := createProgressBar(info.Current, info.Total, progressBarWidth)
	fileName := truncateFileName(info.FileName, maxFileNameLength)
	speed := humanize.Bytes(uint64(info.Speed))
	size := humanize.Bytes(uint64(info.FileSize))

	totalFilesWidth := len(fmt.Sprintf("%d", info.TotalFiles))

	if info.Current >= info.Total && !completedFiles[info.FileIndex] {
		completedFiles[info.FileIndex] = true
		fmt.Printf("\r[%*d/%d] %-*s %s 100%% (complete)\n",
			totalFilesWidth, info.FileIndex, info.TotalFiles,
			maxFileNameLength-1, fileName,
			createProgressBar(1, 1, progressBarWidth))
	} else if !completedFiles[info.FileIndex] {
		fmt.Printf("\r[%*d/%d] %-*s %s %5.1f%% %-8s %5s",
			totalFilesWidth, info.FileIndex, info.TotalFiles,
			maxFileNameLength-1, fileName,
			progressBar,
			percent,
			speed,
			size)
	}
}
