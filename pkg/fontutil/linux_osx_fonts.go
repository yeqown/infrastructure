// +build darwin linux

// Package fontutil .
package fontutil

import "path"

var (
	osxFontPath = path.Join("System", "Library", "Fonts")
)

// GetSysTextList get font list from curretn system
func GetSysTextList() []string {
	return nil
}

// AssemFontPath .
func AssemFontPath(fontfile string) string {
	return path.Join(osxFontPath, fontfile)
}
