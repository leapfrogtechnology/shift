package file

import (
	"mime"
	"os"
	"strings"
)

// GetFileContentType returns the content type of the file.
func GetFileContentType(filename string) string {
	split := strings.Split(filename, ".")

	if len(split) == 0 {
		return "binary/octet-stream"
	}

	ext := split[len(split)-1]

	return mime.TypeByExtension("." + ext)
}

// Exists check existence of file
func Exists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}
