package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func MakeDiagramFilename(prefix, id, format string) string {
	return safeFilename(fmt.Sprintf("%s-%s.%s", prefix, id, format))
}

func MakeMdFilename(prefix, id string) string {
	if id == "" {
		return safeFilename(fmt.Sprintf("%s.md", prefix))
	}
	return safeFilename(fmt.Sprintf("%s-%s.md", prefix, id))
}

var unsafeCharRe = regexp.MustCompile(`[\\\/*:?"<>|\s]`)

func safeFilename(f string) string {
	f = filepath.Clean(filepath.Base(strings.ToLower(unsafeCharRe.ReplaceAllString(f, "_"))))
	return f
}
