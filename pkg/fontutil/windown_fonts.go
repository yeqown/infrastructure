// +build windows

// Package fontutil .
package fontutil

import (
	"os"
	"path"
)

var (
	winFontPath string
)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	winFontPath = path.Join(homedir, "AppData", "Local", "Microsoft", "Windows", "Fonts")
}

// GetSysTextList get font list from curretn system
func GetSysTextList() []string {
	return nil
}

// AssemFontPath .
func AssemFontPath(fontfile string) string {
	return path.Join(winFontPath, fontfile)
}
