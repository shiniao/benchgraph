package main

import (
	"testing"

	"golang.org/x/tools/benchmark/parse"
)

var bTests = []struct {
	line    string // input
	name    string // expected result
	series  string
	arg     string
	nsperop float64
}{
	{"BenchmarkQueue/lock-queue_8-8                  9541816               105 ns/op", "Queue", "lock-queue", "8", 105},
	{"BenchmarkQueue/lock-free-queue_16-8            6674638               180 ns/op", "Queue", "lock-free-queue", "16", 180},
	{"BenchmarkQueue/lock-queue_16-8                19866838               399 ns/op", "Queue", "lock-queue", "16", 399},
}

func TestParser(t *testing.T) {
	for _, tt := range bTests {
		b, _ := parse.ParseLine(tt.line)
		name, series, arg, _, _ := parseNameArgThread(b.Name)
		if name != tt.name {
			t.Errorf("parseNameArgThread(%s): expected %s, actual %s", b.Name, tt.name, name)
		}
		if series != tt.series {
			t.Errorf("parseNameArgThread(%s): expected %s, actual %s", b.Name, tt.series, series)
		}
		if arg != tt.arg {
			t.Errorf("parseNameArgThread(%s): expected %s, actual %s", b.Name, tt.arg, arg)
		}
		if b.NsPerOp != tt.nsperop {
			t.Errorf("parseNameArgThread(%s): expected %f, actual %f", b.Name, tt.nsperop, b.NsPerOp)
		}
	}
}
