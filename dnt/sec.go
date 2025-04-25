// dns transport sec

package dnt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/miekg/dns"
)

// AlgoProvider - implements interface dns.TsigProvider
type sigProvider string

// Generate - implements interface dns.TsigProvider
func (algo sigProvider) Generate(msg []byte, t *dns.TSIG) ([]byte, error) {
	rawSecret, err := fromBase64([]byte(algo))
	if err != nil {
		return nil, err
	}

	algoName := strings.ToLower(t.Algorithm)

	var h hash.Hash
	switch algoName {
	case HmacSHA1:
		h = hmac.New(sha1.New, rawSecret)
	case HmacSHA224:
		h = hmac.New(sha256.New224, rawSecret)
	case HmacSHA256:
		h = hmac.New(sha256.New, rawSecret)
	case HmacSHA384:
		h = hmac.New(sha512.New384, rawSecret)
	case HmacSHA512:
		h = hmac.New(sha512.New, rawSecret)
	case HmacMD5:
		h = hmac.New(md5.New, rawSecret)
	default:
		return nil, fmt.Errorf("algorithm not support %s", algoName)
	}

	h.Write(msg)

	return h.Sum(nil), nil
}

// Verify - copy from dns lib, implements interface dns.TsigProvider
func (algo sigProvider) Verify(msg []byte, t *dns.TSIG) error {
	b, err := algo.Generate(msg, t)
	if err != nil {
		return err
	}
	mac, err := hex.DecodeString(t.MAC)
	if err != nil {
		return err
	}
	if !hmac.Equal(b, mac) {
		return fmt.Errorf("encrypted msg verification failed")
	}
	return nil
}

// fromBase64 - copy from dns lib
func fromBase64(s []byte) (buf []byte, err error) {
	bufLen := base64.StdEncoding.DecodedLen(len(s))
	buf = make([]byte, bufLen)
	n, err := base64.StdEncoding.Decode(buf, s)
	buf = buf[:n]
	return
}
