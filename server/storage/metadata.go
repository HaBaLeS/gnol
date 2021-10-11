package storage

import (
	"path/filepath"
	"regexp"
	"strings"
)

func isCoverImage(f string) bool {
	if !isImageFile(f) {
		return false
	}
	if strings.Contains(strings.ToLower(f), "banner") {
		return false
	}
	return true
}

func isImageFile(f string) bool {
	ext := strings.ToLower(filepath.Ext(f))
	if strings.HasPrefix(f, ".") {
		return false //no . (hidden files)
	}
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func findName(filename string) string {
	out := filename
	//Remove things between ( )
	re := regexp.MustCompile("\\(.*\\)")

	for _, r := range re.FindAll([]byte(filename), -1) {
		out = strings.ReplaceAll(out, string(r), "")
	}

	//clean separators
	out = strings.ReplaceAll(out, "_", " ")
	out = strings.ReplaceAll(out, "-", " ")
	out = strings.ReplaceAll(out, filepath.Ext(filename), "")
	return out
}
