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
	var stat, err = record.parse()
        if err == `nil` {
		stat.report(record.stats)
	}
	
}

func (record *Record) parse() (Stat, string) {
        var response = regexp.MustCompile(record.recordRegex)
	res := response.FindStringSubmatch(record.Message)
	if len(res) == 3 {

		time, err := strconv.ParseFloat(res[2], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, `balls!!!`, err)
		}
		time = time * 1000

		return Stat{time, res[1]}, `nil`
	} else {
		return Stat{0, `null`}, `Stat did not match`
	}
        
}

func (stat *Stat) report(statsd *statsd.StatsdBuffer) {
	statsd.Timing(`requesttime.` + stat.path, int64(Round(stat.responseTime)) )
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
