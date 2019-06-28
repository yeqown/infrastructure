package fs

import (
	"os"
	"path/filepath"
	"strings"
)

// Filter to filter files to be records and walked
type Filter func(info os.FileInfo) bool

// ListFiles .
func ListFiles(root string, filter Filter) []string {
	filenames := make([]string, 0)
	root, _ = filepath.Abs(root)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if filter(info) {
			return nil
		}
		filenames = append(filenames, path)
		return nil
	})

	return filenames
}

// IgnoreDirFilter .
func IgnoreDirFilter() Filter {
	return func(info os.FileInfo) bool {
		return info.IsDir()
	}
}

// IgnoreFiletype .
func IgnoreFiletype(filetyps []string) Filter {
	return func(info os.FileInfo) bool {
		var (
			name = info.Name()
		)

		for _, v := range filetyps {
			if strings.HasSuffix(name, v) {
				return true
			}
		}
		return false
	}
}

// MergeFilter .
func MergeFilter(filters ...Filter) Filter {
	return func(info os.FileInfo) bool {
		for _, filter := range filters {
			if ok := filter(info); ok {
				return true
			}
		}

		return false
	}
}
