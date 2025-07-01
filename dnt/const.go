// dns const

package dnt

import (
	"time"
)

const (
	HmacSHA1   = "hmac-sha1."
	HmacSHA224 = "hmac-sha224."
	HmacSHA256 = "hmac-sha256."
	HmacSHA384 = "hmac-sha384."
	HmacSHA512 = "hmac-sha512."
	HmacMD5    = "hmac-md5.sig-alg.reg.int."
)

const (
	retryTimesMin        = 1
	retryTimesMax        = 10
	defaultRetryTimes    = 5
	defaultDNSServerPort = "53"
	defaultRetryInterval = 3000 // Millisecond

	defaultUDPPkgSize = 4096 // From dig cmd pkg

	defaultTimeout = time.Second * 10 //time
)

const (
	RootDomain = "."

	defaultSOAMBox = "sa.zone.com."
)

const (
	rrOPAdd = "ADD"
	rrOPDel = "DEL"
)

const (
	TypeA     = "A"
	TypeAAAA  = "AAAA"
	TypeCNAME = "CNAME"
	TypeSOA   = "SOA"
	TypeMX    = "MX"
	TypeTXT   = "TXT"
	TypePTR   = "PTR"
	TypeNS    = "NS"
)
