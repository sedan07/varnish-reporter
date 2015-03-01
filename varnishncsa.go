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
}

func (varnishncsa *Varnishncsa) Connect() {
	prefix := "varnish."
	statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
	statsdclient.CreateSocket()
	interval := time.Second * 2 
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
            fmt.Printf("%s \n", scanner.Text())
	    r := Record{scanner.Text(), stats}
	    go r.Process()
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
        }
	cmd.Wait()
}

