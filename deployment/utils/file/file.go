package file

import (
	"mime"
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
