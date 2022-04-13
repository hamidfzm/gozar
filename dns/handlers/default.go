package handlers

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/bogdanovich/dns_resolver"
	roundrobin "github.com/hlts2/round-robin"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"

	ldns "gozar/dns"
)

type DefaultHandler struct {
	server *ldns.Server

	proxies roundrobin.RoundRobin
}

func (h *DefaultHandler) testIs(s1 string) bool {
	for _, s := range h.server.Config.DNS.Domains {
		index := strings.Index(s1, s+".")
		if 0 <= index {
			return true
		}
	}
	return false
}

func otherDNS(s string) string {
	resolver := dns_resolver.New([]string{"8.8.8.8", "8.8.4.4"})
	resolver.RetryTimes = 5

	// c := dns.Client{}
	// m := dns.Msg{}
	// m.SetQuestion(s, dns.TypeA)
	// r, _, err := c.Exchange(&m, "8.8.8.8:53")
	ip, err := resolver.LookupHost(s[0 : strings.Count(s, "")-2])
	if err != nil {
		logrus.Error(err)
	}

	// szIp1 := ip
	// logrus.Info(r)
	// if len(r.Answer) > 0 {
	// 	for _, ans := range r.Answer {
	// 		Arecord := ans.(*dns.A)
	// 		szIp1 = fmt.Sprintf(`%s`, Arecord.A)
	// 		// logrus.Info(Arecord.A, szIp1)
	// 	}
	// }
	// // logrus.Info(t)
	if 0 < len(ip) {
		return fmt.Sprintf(`%s`, ip[0])
	} else {
		return ""
	}
}

func (h *DefaultHandler) parseQuery(m *dns.Msg, addressOfRequester net.Addr) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA:
			if h.testIs(q.Name) {
				proxy := h.proxies.Next().Host
				logrus.Infof("✅ %s ➡ %s", q.Name, proxy)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, proxy))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				return
			}

			szIP1 := otherDNS(q.Name)
			if 0 < len(szIP1) {
				logrus.Infof("❌ %s ➡ %s", q.Name, szIP1)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, szIP1))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func (h *DefaultHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	addressOfRequester := w.RemoteAddr()
	msg := new(dns.Msg)
	msg.SetReply(r)
	switch r.Opcode {
	case dns.OpcodeQuery:
		h.parseQuery(msg, addressOfRequester)
	}

	err := w.WriteMsg(msg)
	if err != nil {
		logrus.Error(err)
	}
}

func Configure(server *ldns.Server) error {
	urls := make([]*url.URL, len(server.Config.DNS.Proxies))
	for i, p := range server.Config.DNS.Proxies {
		urls[i] = &url.URL{Host: p}
	}
	rr, err := roundrobin.New(urls...)
	if err != nil {
		return err
	}

	handler := &DefaultHandler{server: server, proxies: rr}
	server.SetHandler(handler)

	return nil
}
