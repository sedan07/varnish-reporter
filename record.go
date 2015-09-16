package main

import (
	"regexp"
	"fmt"
	"strconv"
	"os"
	"github.com/quipo/statsd"
	"math"
)

type Record struct {
	Message      string
	stats        *statsd.StatsdBuffer
        recordRegex  string
}

type Stat struct {
	responseTime float64
	path         string
}

func (record *Record) Process() {
	stat := record.parse()
	stat.report(record.stats)
}

func (record *Record) parse() Stat {
        var response = regexp.MustCompile(record.recordRegex)
	fmt.Println(response.FindStringSubmatchIndex(record.Message))
	res := response.FindStringSubmatch(record.Message)
	fmt.Println(res)
	if len(res) < 3 {
		panic(`record.recordRegex Regex didn't pull out all the parts from the log entry!`)
	}

	time, err := strconv.ParseFloat(res[2], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, `balls!!!`, err)
	}
	time = time * 1000
	fmt.Println(time)

	return Stat{time, res[1]}
}

func (stat *Stat) report(statsd *statsd.StatsdBuffer) {
	statsd.Timing(`requesttime.` + stat.path, int64(Round(stat.responseTime)) )
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
