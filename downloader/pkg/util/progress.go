package util

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
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

func PrintProgress(info ProgressInfo) {
	percent := (float64(info.Current) / float64(info.Total)) * 100
	fmt.Printf(
		"\r[%d/%d] %s %.2f%% | %s/s | %s | elapsed %s",
		info.FileIndex, info.TotalFiles, info.FileName,
		percent,
		humanize.Bytes(uint64(info.Speed)),
		humanize.Bytes(uint64(info.FileSize)),
		info.Elapsed.Truncate(time.Second).String(),
	)
}
