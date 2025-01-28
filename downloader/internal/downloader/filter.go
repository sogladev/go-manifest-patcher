package downloader

type Filter struct {
	IgnoredFiles map[string]struct{}
}

func NewFilter() *Filter {
	return &Filter{
		IgnoredFiles: map[string]struct{}{
			"README.md": {},
			"main.go":   {},
			"go.sum":    {},
		},
	}
}

func (f *Filter) IsIgnored(file string) bool {
	_, exists := f.IgnoredFiles[file]
	return exists
}
