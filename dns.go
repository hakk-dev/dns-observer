package main

import (
	"fmt"
	"log"
    "net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

type Provider struct {
	Name string
}

var (
	dns_port   int
	responseIP string
	dns_addr   string
	db         = map[string]Provider{}
	rangeMap   = make(map[string][]string)
)

func parseQuery(m *dns.Msg, remote_addr string) {

	for _, q := range m.Question {
		var recType string = ""
		if q.Qtype == 1 {
			recType = "A"
		}
		log.Printf("%s query %s records for %s\n", remote_addr, recType, q.Name)
		switch q.Qtype {
		case dns.TypeA:
			ipAddr, _, err := net.SplitHostPort(remote_addr)
            if err != nil {
                log.Println(err)
            }
			log.Println("remote ip:", ipAddr)
			// add to the db map
			id := getId(q.Name)
			provider := checkProvider(ipAddr)
			log.Println(provider)
			db[id] = Provider{provider}
			rr, err := dns.NewRR(fmt.Sprintf("%s 300 A %s", q.Name, responseIP))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	remote_addr := w.RemoteAddr().String()

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, remote_addr)
	}

	w.WriteMsg(m)
}

func dnsStart() {

	dns.HandleFunc(".", handleDnsRequest)

	go func() {
		// start UDP server
		serverUDP := &dns.Server{Addr: dns_addr + ":" + strconv.Itoa(dns_port), Net: "udp"}
		log.Printf("Starting udp server at %s:%d\n", dns_addr, dns_port)
		err := serverUDP.ListenAndServe()
		defer serverUDP.Shutdown()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		// start TCP server
		serverTCP := &dns.Server{Addr: dns_addr + ":" + strconv.Itoa(dns_port), Net: "tcp"}
		log.Printf("Starting tcp at %s:%d\n", dns_addr, dns_port)
		err := serverTCP.ListenAndServe()
		defer serverTCP.Shutdown()

		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()
}
