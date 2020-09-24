package main

import (
	"errors"
	"regexp"
	"strconv"
)

// Coder should use following naming convention for Benchmark functions
// Naming convention: Benchmark[Function_name]/[Function_Series]_[Function_argument](b *testing.B)
var re *regexp.Regexp = regexp.MustCompile(`Benchmark([a-zA-Z0-9]+)/([a-zA-Z0-9]+)_([_a-zA-Z0-9]+)-([0-9]+)$`)

// Storage for Func(Arg)=Result relations
type BenchArgSet map[string]float64
type BenchNameSet map[string]BenchArgSet

// parseNameArgThread parses function name, argument and number of threads from benchmark output.
func parseNameArgThread(line string) (name string, series string, arg string, c int, err error) {

	arr := re.FindStringSubmatch(line)

	// we expect 4 columns
	if len(arr) != 5 {
		return "", "", "", 0, errors.New("Can't parse benchmark result")
	}

	name, series, arg = arr[1], arr[2], arr[3]

	c, err = strconv.Atoi(arr[4])
	if err != nil {
		return "", "", "", 0, errors.New("Can't parse benchmark result")
	}

	return name, series, arg, c, nil
}
