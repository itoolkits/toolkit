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

// FixDomain fix domain
func FixDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimRight(domain, ".")
	return domain
}

// AdaptHostname adapt hostname
func AdaptHostname(zone string, domain string) (string, error) {
	zone = strings.ToLower(zone)
	zone = strings.Trim(zone, ".")
	if zone == "" {
		return "", nil
	}
	// parse hostname
	idx := strings.LastIndex(domain, zone)
	if idx < 0 {
		return "", fmt.Errorf("parse hostname error, %s, %s", zone, domain)
	}
	if idx == 0 {
		return "@", nil
	}
	return domain[:idx-1], nil
}
