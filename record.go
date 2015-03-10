package main

import (
	"regexp"
	"fmt"
	"strconv"
	"os"
	"github.com/quipo/statsd"
	"math"
	"errors"
)
var response = regexp.MustCompile(`GET http\:\/\/[^/]+([^\s\?]+).*([0-9]+\.[0-9]+)$`)

type Record struct {
	Message  string
	stats    *statsd.StatsdBuffer
}

type Stat struct {
	responseTime float64
	path         string
}

func (record *Record) Process() {
	stat, err := record.parse()
	if err == nil {
		stat.report(record.stats)
	}
}

func (record *Record) parse() (Stat, error) {
	// fmt.Println(response.FindStringSubmatchIndex(record.Message))
	res := response.FindStringSubmatch(record.Message)
	// fmt.Println(res)
	if len(res) < 3 {
		return Stat{0, ""}, errors.New(`Regex didn't pull out all the parts from the log entry!`)
	}

	time, err := strconv.ParseFloat(res[2], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, `balls!!!`, err)
	}
	time = time * 1000
	// fmt.Println(time)

	return Stat{time, res[1]}, nil
}

func (stat *Stat) report(statsd *statsd.StatsdBuffer) {
	statsd.Timing(`requesttime.` + stat.path, int64(Round(stat.responseTime)) )
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
