// +build darwin

// Package fontutil .
package fontutil

import (
	"path"

	"github.com/yeqown/infrastructure/pkg/fs"
)

var (
	osxFontPath = path.Join("/System", "Library", "Fonts")
)

// GetSysDefaultFont . return current system default font
func GetSysDefaultFont() string {
	return "STHeiti Light.ttc"
}

// GetSysFontList get font list from curretn system
func GetSysFontList() (fonts []string) {
	files := fs.ListFiles(osxFontPath, fs.IgnoreDirFilter())
	if len(files) != 0 {
		fonts = make([]string, len(files))
		// true: handle files
		for idx, p := range files {
			_, fonts[idx] = path.Split(p)
		}
	}

	return
}

// AssemFontPath .
func AssemFontPath(fontfile string) string {
	return path.Join(osxFontPath, fontfile)
}
