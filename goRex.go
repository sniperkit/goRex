// Command goRex provides functionality to break a log down into a JSON structured data format,
// which is described by regular expressions.
// The structure of the JSON is determined by the capture groups used in the regular expression.
// Input will be read from os.Stdin. Note, that every line will be interpreted for itself.
// Output will be written to os.Stdout.
// The file containing the regular expression is set with the "regexp" flag. Default: ".regexp"
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// ErrNoMatch describes the failure of matching on a regular expression.
var ErrNoMatch = errors.New("Couldn't match line on regular expression")

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
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var bs []byte
	if bytes.HasSuffix(b, []byte("\r\n")) {
		bs = b[:len(b)-2]
	} else if bytes.HasSuffix(b, []byte("\n")) {
		bs = b[:len(b)-1]
	}
	regexpr, err := regexp.Compile(string(bs))
	if err != nil {
		return nil, err
	}
	return regexpr, nil
}

// ConvertLog reads from r, splits the read data based on the caputure groups of regexpr and writes the result in JSON format into w.
// Note that every line is interpreted for itself, which leads to no support for multiline messages.
func ConvertLog(w io.Writer, r io.Reader, regexpr *regexp.Regexp) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}
		dataSet, err := ExtractDataSet(scanner.Text(), regexpr)
		if err != nil {
			_, err = w.Write([]byte(err.Error() + "\n"))
			if err != nil {
				return err
			}
		} else {
			b, err := json.MarshalIndent(dataSet, "", "\t")
			if err != nil {
				return err
			}
			_, err = w.Write([]byte(string(b) + "\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ExtractDataSet splits s into a map based on the capture groups of regexpr and returns the map.
func ExtractDataSet(s string, regexpr *regexp.Regexp) (map[string]string, error) {
	dataSet := make(map[string]string)
	if !regexpr.MatchString(s) {
		return nil, ErrNoMatch
	}
	matches := regexpr.FindStringSubmatch(s)
	for j, name := range regexpr.SubexpNames() {
		if name != "" {
			dataSet[name] = matches[j]
		}
	}
	return dataSet, nil
}
