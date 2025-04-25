// string util

package str

import (
	"strings"
)

// AddSuffix - add suffix
func AddSuffix(s string, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// RemoveSuffix - remove suffix
func RemoveSuffix(s string, m string) string {
	if strings.HasSuffix(s, m) {
		return s[:len(s)-len(m)]
	}
	return s
}

// RemovePrefix - remove prefix
func RemovePrefix(s string, m string) string {
	if strings.HasPrefix(s, m) {
		return s[len(m):]
	}
	return s
}
