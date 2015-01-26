package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"time"
)

var VERSION = "1.0.0"

var (
	opendns_servers = map[string]int{
		"208.67.222.222": 53,
		"208.67.220.220": 53,
	}

	googledns_servers = map[string]int{
		"8.8.8.8": 53,
		"8.8.4.4": 53,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "iptool"
	app.Usage = "Opinionated tool to perform common queries on connected hosts"
	app.Author = "Andres Villarroel"
	app.Email = "andres.via@gmail.com"
	app.Version = VERSION
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "router",
			Usage:  "Do a DNS request to myip.opendns.com to get your router IP address, using the same technique as in the command `dig +short myip.opendns.com @208.67.222.222` but using GO code.",
			Action: router_action,
		},
		cli.Command{
			Name:   "ip",
			Usage:  "Creates a simple UDP connection to Google or OpenDNS DNS servers and returns the source IP address",
			Action: ip_action,
		},
	}
	app.Run(os.Args)
}

func putinfo(info string) {
	fmt.Printf("%v", info)
	os.Exit(0)
}

func dns_servers() map[string]int {
	m := make(map[string]int)
	for k, v := range opendns_servers {
		m[k] = v
	}
	for k, v := range googledns_servers {
		m[k] = v
	}
	return m
}

func ip_action(ctx *cli.Context) {
	for k, v := range dns_servers() {
		c, err := net.DialTimeout("udp", fmt.Sprintf("%v:%v", k, v), time.Second)
		if ok(err, log.Print, "Error during dial") {
			h, _, err := net.SplitHostPort(c.LocalAddr().String())
			ok(err, log.Print, "Not able to determine our address")
			putinfo(h)
		}
	}
	log.Fatal("Pool of remote servers exhausted")
}

func router_action(ctx *cli.Context) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = []dns.Question{
		dns.Question{
			Name:   "myip.opendns.com.",
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		},
	}
	for k, v := range opendns_servers {
		in, err := dns.Exchange(msg, fmt.Sprintf("%v:%v", k, v))
		if ok(err, log.Print, "Not able to determine router address") {
			for _, ans := range in.Answer {
				if a, ok := ans.(*dns.A); ok {
					putinfo(a.A.String())
				}
			}
		}
	}
	log.Fatal("Pool of remote servers exhausted")
}

func ok(err error, f func(...interface{}), s string) bool {
	if err == nil {
		return true
	} else {
		f(fmt.Sprintf("%v => %v", s, err))
		return false
	}
}
