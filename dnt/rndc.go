// rndc func

package dnt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

// ParseDumpDB parse dump db
func ParseDumpDB(path string) (map[string]map[string][]*RR, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rst := make(map[string]map[string][]*RR)

	var z string
	var v string
	var ok bool
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		line = strings.Trim(line, " ")
		if len(line) < 1 {
			continue
		}
		zone, view, h := extractZone(line)
		if h {
			z = zone
			v = view
			ok = h
			continue
		}
		if !ok {
			continue
		}

		if line[0] == ';' || line[0] == ',' || line[0] == '#' {
			continue
		}
		rr := &RR{}
		err := rr.Unmarshal(z, v, line)
		if err != nil {
			return nil, err
		}

		vrs, h := rst[z]
		if !h {
			rst[z] = map[string][]*RR{
				v: {rr},
			}
		} else {
			rs, h := vrs[v]
			if !h {
				vrs[v] = []*RR{rr}
			} else {
				vrs[v] = append(rs, rr)
			}
		}
	}
	return rst, nil
}

// extractZone extract zone from line
func extractZone(line string) (string, string, bool) {
	seg := strings.Fields(line)
	if len(seg) < 4 {
		return "", "", false
	}
	for i := 0; i < len(seg); i++ {
		if seg[i] == ";" || seg[i] == "," || seg[i] == "#" {
			continue
		}
		if strings.ToLower(seg[i]) != "zone" {
			continue
		}
		if i >= len(seg)-3 {
			continue
		}
		if seg[i+1] != "dump" || seg[i+2] != "of" {
			continue
		}
		s := seg[i+3]
		s = strings.Trim(s, "'")
		seg = strings.Split(s, "/")
		if len(seg) < 2 {
			continue
		}
		if seg[1] != "IN" {
			return "", "", false
		}
		zone := strings.ToLower(seg[0])
		var view string
		if len(seg) == 3 {
			view = seg[2]
		}
		return zone, view, true
	}

	return "", "", false
}

// ParseTextZoneFile parse text zone file
func ParseTextZoneFile(zone string, view string, path string) ([]*RR, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rrs := make([]*RR, 0)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}
		rr := &RR{}
		err = rr.Unmarshal(zone, view, line)
		if err != nil {
			return nil, err
		}
		rrs = append(rrs, rr)
	}
	return rrs, nil
}

// ParseZoneFile parse bind zone file
func ParseZoneFile(zone string, view string, path string) ([]*RR, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rrList := make([]*RR, 0)
	zp := dns.NewZoneParser(f, "", "")
	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		rrConv := &RecordConv{
			Zone:   zone,
			View:   view,
			Record: rr,
		}
		record, err := rrConv.ConvRR()
		if err != nil {
			return nil, err
		}
		rrList = append(rrList, record)
	}
	err = zp.Err()
	if err != nil {
		return nil, err
	}
	return rrList, nil
}

type RecordConv struct {
	Record dns.RR
	Zone   string
	View   string
}

// ConvRR convert dns.RR to RR
func (d *RecordConv) ConvRR() (*RR, error) {
	header := d.Record.Header()
	if header == nil {
		return nil, fmt.Errorf("record header is nil, %s", d.Record.String())
	}

	rr, err := d.newRR()
	if err != nil {
		return nil, err
	}

	rVal, err := d.fromRecordVal()
	if err != nil {
		return nil, err
	}

	rr.RData = rVal
	return rr, nil
}

// newRR create RR
func (d *RecordConv) newRR() (*RR, error) {
	header := d.Record.Header()
	name := FixDomain(header.Name)
	ttl := header.Ttl

	hostname, err := AdaptHostname(d.Zone, name)
	if err != nil {
		return nil, err
	}

	return &RR{
		Zone:     d.Zone,
		View:     d.View,
		TTL:      int(ttl),
		Class:    dns.Class(header.Class).String(),
		Domain:   name,
		Hostname: hostname,
		RType:    dns.TypeToString[header.Rrtype],
	}, nil
}

// fromRecordVal from record val
func (d *RecordConv) fromRecordVal() (string, error) {
	rrType := d.Record.Header().Rrtype
	switch rrType {
	case dns.TypeA:
		if a, ok := d.Record.(*dns.A); ok {
			return a.A.String(), nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeAAAA:
		if a, ok := d.Record.(*dns.AAAA); ok {
			return a.AAAA.String(), nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeCNAME:
		if a, ok := d.Record.(*dns.CNAME); ok {
			return a.Target, nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeNS:
		if a, ok := d.Record.(*dns.NS); ok {
			return a.Ns, nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeMX:
		if a, ok := d.Record.(*dns.MX); ok {
			return strconv.Itoa(int(a.Preference)) + " " + a.Mx, nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypePTR:
		if a, ok := d.Record.(*dns.PTR); ok {
			return a.Ptr, nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeSOA:
		if a, ok := d.Record.(*dns.SOA); ok {
			return a.String(), nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeSRV:
		if a, ok := d.Record.(*dns.SRV); ok {
			return strconv.Itoa(int(a.Priority)) + " " +
				strconv.Itoa(int(a.Weight)) + " " +
				strconv.Itoa(int(a.Port)) + " " + a.Target, nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	case dns.TypeTXT:
		if a, ok := d.Record.(*dns.TXT); ok {
			return strings.Join(a.Txt, "\n"), nil
		}
		return "", fmt.Errorf("record %s convert RR error", d.Record.String())
	default:
		return "", fmt.Errorf("record %s type convert not support", d.Record.String())
	}
}
