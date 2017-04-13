// Command goRex provides functionality to break a log down into a JSON structured data format,
// which is described by regular expressions.
// The structure of the JSON is determined by the capture groups used in the regular expression.
// Input will be read from os.Stdin.
// Output will be written to os.Stdout.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	regexpFlag := flag.String("regexp", ".regexp", "File containing the regular expression")
	flag.Parse()

	regexpFile, err := os.Open(*regexpFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer regexpFile.Close()
	regexpr, err := GetRegexp(regexpFile)
	if err != nil {
		log.Fatal(err)
	}

	err = ConvertLog(os.Stdout, os.Stdin, regexpr)
	if err != nil {
		log.Fatal(err)
	}
}

// GetRegexp reads content from r and tries to interpret it as a regular expression and returns it.
func GetRegexp(r io.Reader) (*regexp.Regexp, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	str := string(in)
	if strings.HasSuffix(str, "\r\n") {
		str = string(in[:len(in)-2])
	} else if strings.HasSuffix(str, "\n") {
		str = string(in[:len(in)-1])
	}
	regexpr, err := regexp.Compile(str)
	if err != nil {
		return nil, err
	}
	return regexpr, nil
}

// ConvertLog reads the content from r and splits the content based on the caputure groups of regexp
// and writes the result into w.
// Note that every line is interpreted for itself, which leads to no support for multiline messages.
func ConvertLog(w io.Writer, r io.Reader, regexpr *regexp.Regexp) error {
	dataSets, err := ExtractDataSets(r, regexpr)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(dataSets, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, string(b))
	if err != nil {
		return err
	}
	return nil
}

// ExtractDataSets reads the content of r and splits the content into maps for every line based on the
// capture groups of regexpr.
func ExtractDataSets(r io.Reader, regexpr *regexp.Regexp) ([]map[string]string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var inputLines []string
	// This is not guarenteed to cut of proper newlines.
	// If we run on a proper OS and got Wintrash logs, we might end up with \r
	// at the end of the lines.
	if strings.ToLower(runtime.GOOS) != "windows" {
		inputLines = strings.Split(string(b), "\n")
	} else {
		inputLines = strings.Split(string(b), "\r\n")
	}
	dataSets := make([]map[string]string, len(inputLines)-1)
	for i, line := range inputLines {
		if line != "" {
			dataSet := make(map[string]string)
			matches := regexpr.FindStringSubmatch(line)
			for j, name := range regexpr.SubexpNames() {
				if name != "" {
					dataSet[name] = matches[j]
				}
			}
			dataSets[i] = dataSet
		}
	}
	return dataSets, nil
}
