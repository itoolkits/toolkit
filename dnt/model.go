// dns model

package dnt

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type RR struct {
	Zone     string `json:"zone"`
	View     string `json:"view"`
	Domain   string `json:"domain"`
	Hostname string `json:"hostname"`
	TTL      int    `json:"ttl"`
	Class    string `json:"class"`
	RType    string `json:"rtype"`
	RData    string `json:"rdata"`
}

// Marshal marshal to string
func (r *RR) Marshal() string {
	return fmt.Sprintf("%s %d %s %s %s", r.Domain, r.TTL, r.Class, r.RType, r.RData)
}

// Unmarshal to struct
func (r *RR) Unmarshal(zone, view, str string) error {
	seg := strings.Fields(str)
	if len(seg) < 5 {
		return fmt.Errorf("record format error, %s", str)
	}
	d := strings.Trim(seg[0], " ")
	d = strings.Trim(d, ".")
	r.Domain = strings.ToLower(d)

	ttls := strings.Trim(seg[1], " ")
	ttl, err := strconv.Atoi(ttls)
	if err != nil {
		return fmt.Errorf("parse ttl error, %s", str)
	}
	r.TTL = ttl
	r.Class = seg[2]
	r.RType = seg[3]
	r.RData = strings.Join(seg[4:], " ")
	r.View = view

	zone = FixDomain(zone)
	if zone == "" {
		return nil
	}
	hostname, err := AdaptHostname(zone, r.Domain)
	if err != nil {
		return err
	}
	r.Hostname = hostname
	return nil
}

type SOA struct {
	NS      string
	MBox    string
	Serial  int64
	Refresh int64
	Retry   int64
	Expire  int64
	MinTTL  int64
}

// Marshal marshal to string
func (s *SOA) Marshal() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", s.NS, s.MBox, s.Serial, s.Refresh, s.Retry, s.Expire, s.MinTTL)
}

// Unmarshal to struct
func (s *SOA) Unmarshal(str string) error {
	seg := strings.Fields(str)
	if len(seg) != 7 {
		return fmt.Errorf("soa record format error, soa:%s", str)
	}
	s.NS = seg[0]
	s.MBox = seg[1]

	serial, err := strconv.ParseInt(seg[2], 10, 64)
	if err != nil {
		return fmt.Errorf("parse soa serial error, err:%s", err.Error())
	}
	s.Serial = serial

	refresh, err := strconv.ParseInt(seg[3], 10, 64)
	if err != nil {
		return fmt.Errorf("parse soa serial error, err:%s", err.Error())
	}
	s.Refresh = refresh

	retry, err := strconv.ParseInt(seg[4], 10, 64)
	if err != nil {
		return fmt.Errorf("parse soa retry error, err:%s", err.Error())
	}
	s.Retry = retry

	expire, err := strconv.ParseInt(seg[5], 10, 64)
	if err != nil {
		return fmt.Errorf("parse soa expire error, err:%s", err.Error())
	}
	s.Expire = expire

	minTTL, err := strconv.ParseInt(seg[6], 10, 64)
	if err != nil {
		return fmt.Errorf("parse soa min ttl value error, err:%s", err.Error())
	}
	s.MinTTL = minTTL

	return nil
}

// Check soa record check
func (s *SOA) Check() error {
	if s.MinTTL < 1 || s.MinTTL > math.MaxInt {
		return fmt.Errorf("soa min ttl out of range")
	}
	if s.Expire < 1 || s.Expire > math.MaxInt {
		return fmt.Errorf("soa expire out of range")
	}
	if s.Retry < 1 || s.Retry > math.MaxInt {
		return fmt.Errorf("soa retry out of range")
	}
	if s.Refresh < 1 || s.Refresh > math.MaxInt {
		return fmt.Errorf("soa refresh out of range")
	}
	if s.NS == "" {
		return fmt.Errorf("soa ns record can not blank")
	}
	if s.MBox == "" {
		return fmt.Errorf("soa mbox record can not blank")
	}
	return nil
}
