package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/sogladev/go-manifest-patcher/pkg/prompt"
	"github.com/sogladev/go-manifest-patcher/pkg/util"
)

const (
	repoOwner = "sogladev"
	repoName  = "go-manifest-patcher"
	apiURL    = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var specialVersionRegex = regexp.MustCompile(`-\w+$`)

func matchVersion(tagName string) bool {
	return !specialVersionRegex.MatchString(tagName)
}

func CheckForUpdate(currentVersion string) (string, string, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch releases: %s", resp.Status)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", "", err
	}

	// Find the latest matching release
	var latestRelease *Release
	for _, release := range releases {
		if matchVersion(release.TagName) {
			latestRelease = &release
			break
		}
	}

	if latestRelease == nil {
		return "", "", nil
	}

	if latestRelease.TagName != currentVersion {
		for _, asset := range latestRelease.Assets {
			if asset.Name == GetExecutableName() {
				return latestRelease.TagName, asset.BrowserDownloadURL, nil
			}
		}
	}

	return "", "", nil
}

func UpdateWithProgress(currentVersion string) error {
	newVersion, downloadURL, err := CheckForUpdate(currentVersion)
	if err != nil {
		return err
	}

	if newVersion == "" {
		fmt.Println("No updates available. You're running the latest version.")
		return nil
	}

	fmt.Printf("\nNew version available: %s -> %s\n", currentVersion, newVersion)
	err = prompt.PromptyN("Do you want to update? [y/N]: ")
	if err != nil {
		return nil // User declined update
	}

	tempFile := GetExecutableName() + ".new"
	err = Download(downloadURL, tempFile)
	if err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Println("Update completed successfully!")
	return nil
}

func Download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	start := time.Now()
	var downloaded int64

	reader := io.TeeReader(resp.Body, &progressWriter{
		out:        out,
		downloaded: &downloaded,
		total:      resp.ContentLength,
		startTime:  start,
	})

	_, err = io.Copy(io.Discard, reader)
	return err
}

type progressWriter struct {
	out        io.Writer
	downloaded *int64
	total      int64
	startTime  time.Time
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.out.Write(p)
	*pw.downloaded += int64(n)
	elapsed := time.Since(pw.startTime)
	speed := float64(*pw.downloaded) / elapsed.Seconds()

	util.PrintProgress(util.ProgressInfo{
		Current:    int(*pw.downloaded),
		Total:      int(pw.total),
		FileIndex:  1,
		TotalFiles: 1,
		Speed:      speed,
		FileSize:   pw.total,
		Elapsed:    elapsed,
		FileName:   "Updating...",
	})

	return n, err
}

func GetExecutableName() string {
	if runtime.GOOS == "windows" {
		return "patcher-windows-amd64.exe"
	}
	return "patcher-linux-amd64"
}

func ReplaceExecutable(newPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Rename(newPath, execPath)
	if err != nil {
		return err
	}

	return nil
}
