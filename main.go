package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tatsushid/go-fastping"
)

type Response struct {
	Host string
	Up   bool
	RTT  []time.Duration `json:"RTT,omitempty"`
}

var PING_COUNT = 1

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Ping a hostname</h1>"+
		"<form action=\"/ping/\" method=\"GET\">"+
		"<input type=\"text\"  name=\"host\"><br>"+
		"<input type=\"submit\" value=\"Ping\">"+
		"</form>")
	return
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("host")
	if name == "" {
		tokens := strings.Split(r.URL.Path, "/")
		if len(tokens) < 3 || tokens[2] == "" {
			http.Error(w, "You have to specify a hostname to ping", http.StatusBadRequest)
			return
		}
		name = tokens[2]
	}
	times := PING_COUNT
	count := r.URL.Query().Get("count")
	fmt.Printf("Requested %s %s times\n", name, count)

	if count != "" {
		if counts, err := strconv.Atoi(count); err == nil {
			times = counts
		}
	}

	up, rtt, err := CheckHost(name, times)

	response := Response{Host: name, Up: up, RTT: rtt}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !up {
		http.Error(w, string(js), http.StatusNotFound)
		return
	}
	w.Write(js)

}

func CheckHost(name string, times int) (bool, []time.Duration, error) {
	fmt.Printf("Checking %s\n", name)
	adds, err := net.LookupHost(name)
	if err != nil {
		return false, nil, err
	}
	fmt.Println("List of addresses:")
	for _, add := range adds {
		fmt.Println("\t", add)
	}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", name)
	if err != nil {
		return false, nil, err
	}
	rtttimes := make([]time.Duration, 0)
	recv := make(chan time.Duration)
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		recv <- rtt
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Do work
		for t := range recv {
			rtttimes = append(rtttimes, t)
			fmt.Printf("Ping %s: %s\n", name, t)
		}
		wg.Done()
	}()

	for i := 0; i < times; i++ {
		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	close(recv)
	wg.Wait()
	fmt.Println("finished")
	return len(rtttimes) > 0, rtttimes, nil
}

func main() {
	name := flag.String("host", "www.google.es", "Hostname to ping")
	count := flag.Int("c", 3, "Number of ping attempts")
	serve := flag.Bool("server", false, "Start the http server")
	address := flag.String("address", ":8080", "Host and port to start the http server on")
	flag.Parse()

	if !*serve {
		_, _, err := CheckHost(*name, *count)
		if err == nil {
			fmt.Printf("could not find host %s: %s\n", *name, err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ping/", pingHandler)
	log.Fatal(http.ListenAndServe(*address, nil))
}
