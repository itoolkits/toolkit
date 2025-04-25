// some domain function

package dnt

import (
	"fmt"
	"strings"
)

// FQD - full query domain
func FQD(domain string) string {
	if len(domain) < 1 {
		return domain
	}
	domain = strings.ToLower(domain)
	if domain[len(domain)-1] == '.' {
		return domain
	}
	return fmt.Sprintf("%s.", domain)
}
