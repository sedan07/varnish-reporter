package main

import (
	"flag"
	"os"
	"fmt"
	"strings"
)

var (
	statsdServer = flag.String("H", "localhost:8125", "Hostname and port of the statsd server")
	statsdPrefix = flag.String("p", "", "Prefix to add to all statsd keys")
	statsdInterval = flag.Int("i", 1000, "Number of milliseconds between flushes to statsd")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: varnish-timer [flags]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()

	if len(*statsdPrefix) > 0 && ! strings.HasSuffix(*statsdPrefix, `.`) {
		fmt.Fprintf(os.Stderr, "warning: statsd prefix must end in a '.'\n")
		os.Exit(2)
	}
	varnishncsa := Varnishncsa{*statsdServer, *statsdPrefix, *statsdInterval}
	varnishncsa.Connect()
}

