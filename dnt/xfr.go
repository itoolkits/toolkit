// dns xfr transport

package dnt

import (
	"log/slog"
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/itoolkits/toolkit/str"
)

type Xfr struct {
	zone   string
	nsInfo string
	mbox   string

	serial uint32

	retryTimes    int
	retryInterval int

	sigAlgorithm string
	sigName      string
	sigSecretKey string

	ns   string
	port string
}

// NewXfr - new xfr structure
func NewXfr(zone string, nsInfo string, mbox string, serial uint32) *Xfr {
	xfr := &Xfr{
		zone:          str.AddSuffix(zone, RootDomain),
		nsInfo:        str.AddSuffix(nsInfo, RootDomain),
		mbox:          str.AddSuffix(mbox, RootDomain),
		serial:        serial,
		retryTimes:    defaultRetryTimes,
		retryInterval: defaultRetryInterval,
	}
	if xfr.mbox == "" {
		xfr.mbox = defaultSOAMBox
	}
	return xfr
}

// SetNS - set ns
func (x *Xfr) SetNS(add string, port string) {
	x.ns = add
	x.port = port
}

// SetRetry - change retry info
func (x *Xfr) SetRetry(retryTimes, retryInterval int) {
	if retryTimes < retryTimesMin {
		retryTimes = retryTimesMin
	}
	if retryTimes > retryTimesMax {
		retryTimes = retryTimesMax
	}
	x.retryTimes = retryTimesMin
	x.retryInterval = retryInterval
}

// SetAlgo - sign transaction
func (x *Xfr) SetAlgo(algo, sigName, secretKey string) {
	x.sigAlgorithm = algo
	x.sigName = str.AddSuffix(sigName, RootDomain)
	x.sigSecretKey = secretKey
}

// Query - query xfr rr list
func (x *Xfr) Query() ([]dns.RR, error) {
	var rst []dns.RR
	var err error
	for i := 0; i < x.retryTimes; i++ {
		if i > 0 && x.retryInterval > 0 {
			time.Sleep(time.Duration(x.retryInterval) * time.Millisecond)
		}

		// call axfr when serial < 1
		if x.serial < 1 {
			rst, err = x.axfr()
		} else {
			rst, err = x.ixfr()
		}

		// for retry
		if err != nil {
			continue
		}

		// break when no error
		break
	}

	return rst, nil
}

// ixfr - ixfr transfer
func (x *Xfr) ixfr() ([]dns.RR, error) {
	msg := &dns.Msg{}
	msg = msg.SetIxfr(x.zone, x.serial, x.nsInfo, x.mbox)
	return x.xfr(msg)
}

// axfr - axfr transfer
func (x *Xfr) axfr() ([]dns.RR, error) {
	msg := &dns.Msg{}
	msg = msg.SetAxfr(x.zone)
	return x.xfr(msg)
}

// xfr - xfr transfer
func (x *Xfr) xfr(msg *dns.Msg) ([]dns.RR, error) {
	tx := &dns.Transfer{
		DialTimeout:  defaultTimeout,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	if x.sigName != "" && x.sigAlgorithm != "" {
		msg.SetTsig(x.sigName, x.sigAlgorithm, 300, time.Now().Unix())
		tx.TsigProvider = sigProvider(x.sigSecretKey)
	}

	ch, err := tx.In(msg, net.JoinHostPort(x.ns, x.port))
	if err != nil {
		slog.Error("dnt xfr transfer error",
			"server", x.ns, "port", x.port, "msg", msg, "error", err)
		return nil, err
	}

	rst := make([]dns.RR, 0)
	for c := range ch {
		if c.Error != nil {
			slog.Error("dnt xfr channel error",
				"server", x.ns, "port", x.port, "msg", msg, "msg", c.Error)
			return nil, err
		}
		rst = append(rst, c.RR...)
	}

	return rst, nil
}

type RROP struct {
	OP     string
	Serial uint32
	RRs    []dns.RR
}

// SeparateXfrRRs - separate xfr record list, record add rrs, delete rrs and soa
func SeparateXfrRRs(rrs []dns.RR) ([]*RROP, *dns.SOA) {
	var tempRRs []dns.RR

	var rrOPs []*RROP

	var lastSoa *dns.SOA
	for _, rr := range rrs {
		switch rr.(type) {
		case *dns.SOA:
			curSoa := rr.(*dns.SOA)
			if lastSoa != nil && len(tempRRs) > 0 {
				rrOP := &RROP{
					RRs:    tempRRs,
					Serial: curSoa.Serial,
				}
				rrOPs = append(rrOPs, rrOP)
				if lastSoa.Serial == curSoa.Serial {
					rrOP.OP = rrOPAdd
					tempRRs = make([]dns.RR, 0)
				} else {
					// ignore reduce condition
					rrOP.OP = rrOPDel
					tempRRs = make([]dns.RR, 0)
				}
			}
			lastSoa = curSoa
		default:
			tempRRs = append(tempRRs, rr)
		}
	}
	return rrOPs, lastSoa
}
