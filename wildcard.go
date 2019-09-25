package main

import (
	"path/filepath"
	"strings"
)

func expandPattern(pattern string) string {
	filenames, _ := filepath.Glob(pattern)
	return strings.Join(filenames, " ")
}
