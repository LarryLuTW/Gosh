package main

import (
	"path/filepath"
	"strings"
)

func expandPattern(pattern string) string {
	filenames, _ := filepath.Glob(pattern)
	return strings.Join(filenames, " ")
}

func expandWildcardInCmd(input string) string {
	args := strings.Split(input, " ")
	for i, arg := range args {
		if strings.Contains(arg, "*") || strings.Contains(arg, "?") {
			args[i] = expandPattern(arg)
		}
	}
	return strings.Join(args, " ")
}
