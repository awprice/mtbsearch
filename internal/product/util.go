package product

import (
	"github.com/iancoleman/strcase"
	"regexp"
	"strconv"
	"strings"
)

var (
	millimetreRegex = regexp.MustCompile("^([\\d\\.]+)mm$")
	inchRegex = regexp.MustCompile("^([\\d\\.]+)\"$")
)

func cleanHeader(in string) string {
	camel := strcase.ToSnake(in)
	replacer := strings.NewReplacer(".", "", "/", "")
	fixed := replacer.Replace(camel)
	return strings.ReplaceAll(fixed, "__", "_")
}

func extractFloat64(value string, regex *regexp.Regexp) (float64, bool) {
	if !regex.Match([]byte(value)) {
		return 0, false
	}

	matches := regex.FindAllStringSubmatch(value, -1)
	if len(matches) != 1 {
		return 0, false
	}

	if len(matches[0]) != 2 {
		return 0, false
	}

	s := matches[0][1]
	res, err := strconv.ParseFloat(s, 64)
	return res, err == nil
}

func extractMillimetres(value string) (float64, bool) {
	return extractFloat64(value, millimetreRegex)
}

func extractInches(value string) (float64, bool) {
	return extractFloat64(value, inchRegex)
}