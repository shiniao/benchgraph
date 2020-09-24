package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/fatih/color"
	"golang.org/x/tools/benchmark/parse"
)

// uploadData sends data to server and expects graph url.
func uploadData(apiUrl, data, title string) (string, error) {

	resp, err := http.PostForm(apiUrl, url.Values{"data": {data}, "title": {title}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("Server din't return graph URL")
	}

	return string(body), nil
}

func main() {

	var oBenchNames, oBenchArgs stringList

	// graph elements will be ordered as in benchmark output by default - unless the order was specified here
	flag.Var(&oBenchNames, "obn", "comma-separated list of benchmark names")
	flag.Var(&oBenchArgs, "oba", "comma-separated list of benchmark arguments")
	title := flag.String("title", "Benchmark results in ns/op", "title of a graph")
	apiUrl := flag.String("apiurl", "http://benchgraph.codingberg.com", "url to server api")
	flag.Parse()

	var skipBenchNamesParsing, skipBenchArgsParsing bool

	if oBenchNames.Len() > 0 {
		skipBenchNamesParsing = true
	}
	if oBenchArgs.Len() > 0 {
		skipBenchArgsParsing = true
	}

	benchResults := make(BenchNameSet)

	// parse Golang benchmark results, line by line
	scan := bufio.NewScanner(os.Stdin)
	green := color.New(color.FgGreen).SprintfFunc()
	red := color.New(color.FgRed).SprintFunc()
	for scan.Scan() {
		line := scan.Text()

		mark := green("√")

		b, err := parse.ParseLine(line)
		if err != nil {
			mark = red("?")
		}

		// read bench name and arguments
		if b != nil {
			name, series, arg, _, err := parseNameArgThread(b.Name)
			fmt.Printf("series: %s, arg: %s", series, arg)
			if err != nil {
				mark = red("!")
				fmt.Printf("%s %s\n", mark, line)
				continue
			}

			if !skipBenchNamesParsing && !oBenchNames.stringInList(series) {
				oBenchNames.Add(series)
			}

			if !skipBenchArgsParsing && !oBenchArgs.stringInList(arg) {
				oBenchArgs.Add(arg)
			}

			if _, ok := benchResults[series]; !ok {
				benchResults[series] = make(BenchArgSet)
			}

			benchResults[series][arg] = b.NsPerOp
		}

		fmt.Printf("%s %s\n", mark, line)
	}

	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v", err)
		os.Exit(1)
	}

	if len(benchResults) == 0 {
		fmt.Fprintf(os.Stderr, "no data to show.\n\n")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Waiting for server response ...")

	data := graphData(benchResults, oBenchNames, oBenchArgs)

	graphUrl, err := uploadData(*apiUrl, string(data), *title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "uploading data: %v", err)
		os.Exit(1)
	}

	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println(graphUrl)
	fmt.Println()
	fmt.Println("=========================================")

}
