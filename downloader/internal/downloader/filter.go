package downloader

import (
	"path/filepath"
	"strings"
)

// Filter struct holds various patterns to ignore files
type Filter struct {
	ExactMatches     map[string]struct{}
	ExtensionMatches map[string]struct{}
	GlobPatterns     []string
}

// NewFilter initializes a new Filter with predefined patterns
func NewFilter() *Filter {
	return &Filter{
		ExactMatches: map[string]struct{}{
			"README.md": {},
			"go.sum":    {},
			"go.mod":    {},
			// Add more exact filenames as needed
		},
		ExtensionMatches: map[string]struct{}{
			".gitignore": {},
			".env":       {},
			// Add more file extensions to ignore
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
	// Get the base name of the file
	base := filepath.Base(path)

	// Check for exact matches
	if _, exists := f.ExactMatches[base]; exists {
		return true
	}

	// Check for extension matches
	ext := filepath.Ext(base)
	if _, exists := f.ExtensionMatches[ext]; exists {
		return true
	}

	// Check for glob patterns
	for _, pattern := range f.GlobPatterns {
		matched, err := filepath.Match(pattern, base)
		if err == nil && matched {
			return true
		}

		// For patterns with directories, use Match on the full path
		if strings.Contains(pattern, "/") {
			matched, err := filepath.Match(pattern, path)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}
