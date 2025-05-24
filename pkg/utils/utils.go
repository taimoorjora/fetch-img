package utils

import (
	"path/filepath"
	"strings"
)

// SanitizeFilename removes invalid characters from filename
func SanitizeFilename(filename string) string {
	// Remove invalid characters
	filename = strings.Map(func(r rune) rune {
		if r < 32 || r == '\\' || r == '/' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, filename)

	// Ensure the filename isn't too long
	if len(filename) > 255 {
		filename = filename[:255]
	}

	return filepath.Clean(filename)
}
