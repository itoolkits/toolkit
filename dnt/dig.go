// dns dig util

package dnt

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/itoolkits/toolkit/str"
)

type Dig struct {
	domain string
	ns     string
	port   string

	rType uint16

	tcp bool

	retryTimes    int
	retryInterval int

	sigAlgorithm string
	sigName      string
	sigSecretKey string

	subnetIPFamily uint16 // 1-ipv4, 2-ipv6
	subnetIP       net.IP
	subnetIPMask   uint8
}

// NewDig - create dig struct
func NewDig(domain, ns string, rType uint16) *Dig {
	return &Dig{
		domain:        domain,
		ns:            ns,
		retryTimes:    defaultRetryTimes,
		retryInterval: defaultRetryInterval,
		tcp:           true,
		rType:         rType,
		port:          defaultDNSServerPort,
	}
}

// SetSubNet set sub net, family(1:ipv4, 2:ipv6)
func (d *Dig) SetSubNet(addr net.IP, ipFamily uint16, mask uint8) {
	d.subnetIP = addr
	d.subnetIPFamily = ipFamily
	d.subnetIPMask = mask
}

// SetRetry - change retry info
func (d *Dig) SetRetry(retryTimes, retryInterval int) {
	if retryTimes < retryTimesMin {
		retryTimes = retryTimesMin
	}
	if retryTimes > retryTimesMax {
		retryTimes = retryTimesMax
	}
	d.retryTimes = retryTimesMin
	d.retryInterval = retryInterval
}

// SetPort - change default port
func (d *Dig) SetPort(port string) *Dig {
	d.port = port
	return d
}

// SetAlgo - sign transaction
func (d *Dig) SetAlgo(algo, sigName, secretKey string) *Dig {
	d.sigAlgorithm = algo
	d.sigName = str.AddSuffix(sigName, RootDomain)
	d.sigSecretKey = secretKey
	return d
}

// Query - dns query
func (d *Dig) Query() ([]dns.RR, error) {
	var rrs []dns.RR
	var err error

	// Retry if error
	for i := 0; i < d.retryTimes; i++ {
		if i > 0 && d.retryInterval > 0 {
			time.Sleep(time.Duration(d.retryInterval) * time.Millisecond)
		}

		rrs, err = d.dnsQuery("udp")
		if err != nil {
			slog.Warn("DNS Query Error. ", "ns", d.ns, "domain", d.domain, "error", err.Error())
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	return rrs, err
}

func (d *Dig) buildSubNetExtra() *dns.OPT {
	return &dns.OPT{
		Hdr: dns.RR_Header{
			Name:   ".",
			Rrtype: dns.TypeOPT,
			Class:  dns.DefaultMsgSize,
		},
		Option: []dns.EDNS0{
			&dns.EDNS0_SUBNET{
				Code:          dns.EDNS0SUBNET,
				Family:        d.subnetIPFamily,
				SourceNetmask: d.subnetIPMask,
				SourceScope:   0,
				Address:       d.subnetIP,
			},
		},
	}
}

// dnsQuery - private function, dns query, support protocol
func (d *Dig) dnsQuery(protocol string) ([]dns.RR, error) {
	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(d.domain), d.rType)

	// edns subnet
	if d.subnetIP != nil && len(d.subnetIP) > 0 {
		msg.Extra = append(msg.Extra, d.buildSubNetExtra())
	}

	// Create client
	client := &dns.Client{
		Net: protocol,
	}

	var secretKeyProvider sigProvider
	// Need sign
	if d.sigAlgorithm != "" {
		msg.SetEdns0(defaultUDPPkgSize, false)
		msg.SetTsig(d.sigName, d.sigAlgorithm, 300, time.Now().Unix())

		secretKeyProvider = sigProvider(d.sigSecretKey)
	}

	if secretKeyProvider != "" {
		client.TsigProvider = secretKeyProvider
	}

	// Net exchange
	rMsg, _, err := client.Exchange(msg, net.JoinHostPort(d.ns, d.port))

	// error
	if err != nil {
		return nil, err
	}

	// msg nil
	if rMsg == nil {
		return nil, fmt.Errorf("DNS Query Response Msg Nil")
	}

	// Use tcp connection
	if rMsg.Truncated && client.Net == "udp" && d.tcp {
		return d.dnsQuery("tcp")
	}

	// NXDomain - Non-Existent Domain
	if rMsg.Rcode != dns.RcodeNameError && rMsg.Rcode != dns.RcodeSuccess {
		return make([]dns.RR, 0), nil
	}

	// Convert RR and return
	if len(rMsg.Answer) > 0 {
		return rMsg.Answer, nil
	}

	if d.rType == dns.TypeNS {
		return rMsg.Ns, nil
	}
	return nil, nil
}
