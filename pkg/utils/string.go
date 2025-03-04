package utils

import "strings"

func UCFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func LCFirst(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func LessString(a, b string) bool {
	return a < b
}
