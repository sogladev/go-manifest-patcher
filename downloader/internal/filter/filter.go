package filter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
)

// Filter struct holds various patterns to ignore files
type Filter struct {
	ExactMatches     []string `json:"exact_matches"`
	ExtensionMatches []string `json:"extension_matches"`
	GlobPatterns     []string `json:"glob_patterns"`
	BaseMatches      []string `json:"base_matches"`
	ExcludePatterns  []string `json:"exclude_patterns"`
}

// LoadFilter loads the filter configuration from a JSON file
func LoadFilter(filename string) (*Filter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var f Filter
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// SaveFilter saves the filter configuration to a JSON file
func SaveFilter(filename string, f *Filter) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON with indentation
	err = encoder.Encode(f)
	if err != nil {
		return err
	}

	return nil
}

// IsIgnored checks if the given file path should be ignored based on the filter patterns
func (f *Filter) IsIgnored(path string) bool {
	// First check if path matches any exclude patterns
	for _, pattern := range f.ExcludePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(path) {
			logger.Debug.Printf("[-] Not ignored (ExcludePattern): %s\n", path)
			return false
		}
	}

	// Check exact matches
	filename := filepath.Base(path)
	for _, exact := range f.ExactMatches {
		if filename == exact {
			logger.Debug.Printf("[+] Ignored (ExactMatch): %s\n", path)
			return true
		}
	}

	// Check extension matches
	for _, ext := range f.ExtensionMatches {
		if strings.HasSuffix(filename, ext) {
			logger.Debug.Printf("[+] Ignored (ExtensionMatch): %s\n", path)
			return true
		}
	}

	// Check base paths (exact matches with path separators normalized)
	normalizedPath := filepath.ToSlash(path)
	for _, basePath := range f.BaseMatches {
		if normalizedPath == filepath.ToSlash(basePath) {
			logger.Debug.Printf("[+] Ignored (BasePath): %s\n", path)
			return true
		}
	}

	// Check glob patterns
	for _, pattern := range f.GlobPatterns {
		g := glob.MustCompile(pattern)
		if g.Match(path) {
			logger.Debug.Printf("[+] Ignored (GlobMatch): %s\n", path)
			return true
		}
	}

	logger.Debug.Printf("[-] Not ignored: %s\n", path)
	return false
}
