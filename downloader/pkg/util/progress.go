package util

import (
	"fmt"
	"time"
)

type ProgressInfo struct {
	Current    int
	Total      int
	FileIndex  int
	TotalFiles int
	Speed      float64
	FileSizeMB float64
	Elapsed    time.Duration
	FileName   string
}

func PrintProgress(info ProgressInfo) {
	percent := (float64(info.Current) / float64(info.Total)) * 100
	fmt.Printf(
		"\r[%d/%d] %s %.2f%% | %.2f MiB/s | %.2f MiB | elapsed %s",
		info.FileIndex, info.TotalFiles, info.FileName, percent, info.Speed, info.FileSizeMB, info.Elapsed.String(),
	)
}
