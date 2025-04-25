// ip util

package ipt

import (
	"net"
	"strings"
)

const (
	maxLen   = len("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	hexDigit = "0123456789abcdef"
)

// IPFullFmt - ip full format
func IPFullFmt(ip string) string {
	if !strings.Contains(ip, ":") {
		return ip
	}

	p := net.ParseIP(ip)
	if len(p) != net.IPv6len {
		return ""
	}
	b := make([]byte, 0, maxLen)

	// Print with possible :: in place of run of zeros
	for i := 0; i < net.IPv6len; i += 2 {
		b = appendHex(b, (uint32(p[i])<<8)|uint32(p[i+1]))
		if i != net.IPv6len-2 {
			b = append(b, ':')
		}
	}
	return string(b)
}

// appendHex - append hex
func appendHex(dst []byte, i uint32) []byte {
	if i == 0 {
		return append(dst, '0', '0', '0', '0')
	}

	x := make([]byte, 0)
	for j := 7; j >= 0; j-- {
		v := i >> uint(j*4)
		if v > 0 {
			x = append(x, hexDigit[v&0xf])
		}
	}

	dis := 4 - len(x)
	for j := 0; j < dis; j++ {
		dst = append(dst, '0')
	}
	return append(dst, x...)
}
