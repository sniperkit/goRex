package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestGetRegexp(t *testing.T) {
	testStr := "AaBbCcZz 1234567890"
	exprStr := `(?P<Letters>[a-zA-Z]+)[\s\t]*(?P<Numbers>[0-9]+)`
	r := strings.NewReader(exprStr)
	expr, err := GetRegexp(r)
	if err != nil {
		t.Error(err.Error())
	}
	if !expr.MatchString(testStr) {
		t.Error("Regular expression is incorrect")
	}
}

func TestExtractDataSet(t *testing.T) {
	testStr := "AaBbCcZz 1234567890"
	exprStr := `(?P<Letters>[a-zA-Z]+)[\s\t]*(?P<Numbers>[0-9]+)`
	r := strings.NewReader(exprStr)
	expr, err := GetRegexp(r)
	if err != nil {
		t.Error(err.Error())
	}
	m, err := ExtractDataSet(testStr, expr)
	if err != nil {
		t.Error(err.Error())
	}
	if m["Letters"] != "AaBbCcZz" {
		t.Error("Letters were not parsed correctly")
	}
	if m["Numbers"] != "1234567890" {
		t.Error("Numbers were not parsed correctly")
	}
}

func TestConvertLog(t *testing.T) {
	testRegexp := regexp.MustCompile(`(?P<Letters>[a-zA-Z]+)[\s\t]*(?P<Numbers>[0-9]+)`)
	testStr := "a1\n"
	expectedStr := `{
	"Letters": "a",
	"Numbers": "1"
}
`
	buf := new(bytes.Buffer)
	r := strings.NewReader(testStr)
	ConvertLog(buf, r, testRegexp)
	if buf.String() != expectedStr {
		t.Errorf("got:\n%v\nwant:\n%v\n", buf.String(), expectedStr)
	}
}
