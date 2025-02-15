package utils

import "strings"

func UCFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
