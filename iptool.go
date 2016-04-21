package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/miekg/dns"
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

var VERSION = "1.0.3"

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
	app.Usage = "Opinionated tool to perform common IP queries on connected hosts"
	app.Author = "Andres Villarroel"
	app.Email = "andres.via@gmail.com"
	app.Version = VERSION
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "router",
			Usage:  "Do a DNS request to myip.opendns.com to get your router IP address",
			Action: router_action,
		},
		cli.Command{
			Name:   "ip",
			Usage:  "Creates a simple UDP/53 connection to Google or OpenDNS and returns the source IP address",
			Action: ip_action,
		},
		cli.Command{
			Name:   "lan",
			Usage:  "alias of 'ip' command",
			Action: ip_action,
		},
		cli.Command{
			Name:   "docker",
			Usage:  "Obtain address from $DOCKER_HOST, docker.local or local.docker. If /var/run/docker.sock is a socket uses 'lan' or 127.0.0.1 if something answers on 2375 or 2376",
			Action: docker_action,
		},
		cli.Command{
			Name:   "version",
			Usage:  "Version number only",
			Action: func(ctx *cli.Context) { putinfo(app.Version) },
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

func docker_action(ctx *cli.Context) {
	resolve_func_map := map[string]func(chan string, chan bool){
		"DOCKER_HOST":  resolve_from_env,
		"docker.local": func(s chan string, b chan bool) { resolve_from_lookup("docker.local", s, b) },
		"local.docker": func(s chan string, b chan bool) { resolve_from_lookup("local.docker", s, b) },
	}
	port_func_map := map[string]func(chan bool, chan bool){
		"2375": func(b1, b2 chan bool) { tcp_check("2375", b1, b2) },
		"2376": func(b1, b2 chan bool) { tcp_check("2376", b1, b2) },
	}
	resolve_resolve := make(chan string, len(resolve_func_map))
	resolve_done := make(chan bool, len(resolve_func_map))
	resolve_all_done := make(chan bool)
	port_resolve := make(chan bool, len(port_func_map))
	port_done := make(chan bool, len(port_func_map))
	port_all_done := make(chan bool)
	for _, fun := range resolve_func_map {
		go fun(resolve_resolve, resolve_done)
	}
	go func() {
		for i := len(resolve_func_map); i > 0; i-- {
			<-resolve_done
		}
		resolve_all_done <- true
	}()
	for _, fun := range port_func_map {
		go fun(port_resolve, port_done)
	}
	go func() {
		for i := len(port_func_map); i > 0; i-- {
			<-port_done
		}
		port_all_done <- true
	}()
	select {
	case info := <-resolve_resolve:
		putinfo(info)
	case <-resolve_all_done:
		if file_info, err := os.Stat("/var/run/docker.sock"); err == nil {
			if file_info.Mode()&os.ModeSocket == os.ModeSocket {
				ip_action(ctx)
				return
			}
		}
		select {
		case <-port_resolve:
			putinfo("127.0.0.1")
		case <-port_all_done:
			log.Fatal("All options to get docker address exhausted")
		}
	}
}

func tcp_check(port string, resolve, done chan bool) {
	if conn, err := net.Dial("tcp", "127.0.0.1:"+port); err == nil {
		defer conn.Close()
		resolve <- true
	}
	done <- true
}

func resolve_from_env(resolve chan string, done chan bool) {
	docker_host := os.Getenv("DOCKER_HOST")
	if docker_host != "" {
		if docker_url, err := url.Parse(docker_host); err == nil {
			if host, _, err := net.SplitHostPort(docker_url.Host); err == nil {
				resolve <- host
			}
		}
	}
	done <- true
}

func resolve_from_lookup(lookup_host string, resolve chan string, done chan bool) {
	if names, err := net.LookupHost(lookup_host); err == nil {
		resolve <- names[0]
	}
	done <- true
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
