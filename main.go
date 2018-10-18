package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

func main() {
	name := flag.String("host", "www.google.es", "Hostname to ping")
	count := flag.Int("c", 3, "Number of times to wait for the ")
	flag.Parse()

	adds, err := net.LookupHost(*name)
	if err != nil {
		panic(err)
	}
	fmt.Println("List of addresses:")
	for _, add := range adds {
		fmt.Println("\t", add)
	}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", *name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	for i := 0; i < *count; i++ {
		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("finished")

}
