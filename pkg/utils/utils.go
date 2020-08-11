package utils

import "strings"

func IsNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}
