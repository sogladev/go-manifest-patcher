package filter

// DefaultFilter initializes a new Filter with predefined patterns
func DefaultFilter() *Filter {
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
			".log",
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
