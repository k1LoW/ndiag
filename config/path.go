package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func ImagePath(prefix string, vals []string, format string) string {
	return safeFilename(fmt.Sprintf("%s-%s.%s", prefix, strings.Join(vals, "-"), format))
}

func MdPath(prefix string, vals []string) string {
	return safeFilename(fmt.Sprintf("%s-%s.md", prefix, strings.Join(vals, "-")))
}

var unsafeCharRe = regexp.MustCompile(`[\\\/*:?"<>|]`)

func safeFilename(f string) string {
	f = filepath.Clean(filepath.Base(strings.ToLower(unsafeCharRe.ReplaceAllString(f, "_"))))
	return f
}
