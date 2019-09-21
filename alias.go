package main

import "strings"

var aliasTable = map[string]string{}

func setAlias(key, value string) {
	aliasTable[key] = value
}

func unsetAlias(key string) {
	delete(aliasTable, key)
}

func expandAlias(input string) string {
	args := strings.SplitN(input, " ", 2)
	cmd := args[0]

	if expandedCmd, ok := aliasTable[cmd]; ok {
		return strings.Replace(input, cmd, expandedCmd, 1)
	}

	return input
}
