package filter

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func CollectExtraFiles(f *Filter) (map[string]bool, error) {
	localFiles := map[string]bool{}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if !f.IsIgnored(path) {
				localFiles[path] = true
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading local files: %v", err)
	}
	return localFiles, nil
}
