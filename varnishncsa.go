package main

import (
	"fmt"
	"log"
	"os/exec"
	"bufio"
	"os"
	"github.com/quipo/statsd"
	"time"
)

type Varnishncsa struct {
	StatsdServer        string
        StatsdPrefix        string
	StatsdSendInterval  int
}

func (varnishncsa *Varnishncsa) Connect() {
	statsdclient := statsd.NewStatsdClient(varnishncsa.StatsdServer, varnishncsa.StatsdPrefix)
	statsdclient.CreateSocket()
	interval := time.Duration(int(time.Millisecond) * varnishncsa.StatsdSendInterval)
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	cmd := exec.Command("varnishncsa", "-F", `%h %l %u %t "%r" %s %b "%{Referer}i" "%{User-agent}i" %{Varnish:time_firstbyte}x`)
	stdout, outerr := cmd.StdoutPipe()
	if outerr != nil {
		log.Fatal(outerr)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
	    r := Record{scanner.Text(), stats}
	    r.Process()
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "There was an error with the scanner attached to varnishncsa", err)
        }
	cmd.Wait()
}

