package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/sogladev/go-manifest-patcher/pkg/util"
)

const (
	repoOwner             = "sogladev"
	repoName              = "go-manifest-patcher"
	apiURL                = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases"
	executableNameLinux   = "patcher-epoch-linux-amd64"
	executableNameWindows = "patcher-epoch-windows-amd64.exe"
)

type releases struct {
	TagName    string `json:"tag_name"`
	PreRelease bool   `json:"prerelease"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type LatestRelease struct {
	Version string
	url     string
}

func (r *LatestRelease) Download() error {
	tempFile := getExecutableName() + ".new"
	if err := download(r.url, tempFile); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("update failed: %w", err)
	}

	if runtime.GOOS == "windows" {
		// On Windows we cannot replace the running executable, so inform the user.
		fmt.Printf("Update downloaded as: %s.\nPlease rename the new executable to %s and restart the application!\n", tempFile, getExecutableName())
		return nil
	}

	if err := replaceExecutable(tempFile); err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	if err := setExecutablePermission(); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	fmt.Println("Update successful. Please restart the application.")
	return nil
}

func Fetch(currentVersion string) (*LatestRelease, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases: %s", resp.Status)
	}

	var releases []releases
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	// Loop through releases and find the best special edition
	bestVersion := ""
	bestURL := ""
	for _, rel := range releases {
		if rel.PreRelease {
			continue // Skip pre-releases
		}
		if MatchSpecialEdition(rel.TagName) {
			if bestVersion == "" || CompareVersions(rel.TagName, bestVersion) > 0 {
				// Update best version
				bestVersion = rel.TagName
				// Choose first asset or refine for each release
				if len(rel.Assets) > 0 {
					bestURL = rel.Assets[0].BrowserDownloadURL
				}
			}
		}
	}

	// Compare our best with current version
	if bestVersion != "" && CompareVersions(bestVersion, currentVersion) > 0 {
		return &LatestRelease{
			Version: bestVersion,
			url:     bestURL,
		}, nil
	}

	return nil, nil
}

func download(url, dest string) error {
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

func getExecutableName() string {
	if runtime.GOOS == "windows" {
		return executableNameWindows
	}
	return executableNameLinux
}

func replaceExecutable(newPath string) error {
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

func setExecutablePermission() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	return os.Chmod(execPath, 0755)
}
