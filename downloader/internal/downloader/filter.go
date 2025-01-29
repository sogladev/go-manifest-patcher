package downloader

import (
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/sogladev/golang-terminal-downloader/downloader/internal/logger"
)

// Filter struct holds various patterns to ignore files
type Filter struct {
	ExactMatches     []string
	ExtensionMatches []string
	GlobPatterns     []string
	BaseMatches      []string
	ExcludePatterns  []string
}

// NewFilter initializes a new Filter with predefined patterns
func NewFilter() *Filter {
	return &Filter{
		ExcludePatterns: []string{
			// "Documentation/*",   // Don't ignore documentation
		},
		ExactMatches: []string{
			"README.md",
			"go.sum",
			"go.mod",
			// Add more exact filenames as needed
		},
		ExtensionMatches: []string{
			".gitignore",
			".env",
			// Add more file extensions to ignore
		},
		BaseMatches: []string{
			"manifest.json",
			// Add more base paths as needed
		},
		GlobPatterns: []string{
			"*.log",     // Ignore all .log files
			"temp/*",    // Ignore all files in temp directory
			"*.tmp",     // Ignore all .tmp files
			"*.bak",     // Ignore all .bak files
			"*.go",      // Ignore all .go files
			"docs/*.md", // Ignore markdown files in docs directory
			// Add more glob patterns as needed
		},
	}
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
