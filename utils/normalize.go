package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	rxNormalizeUsername = regexp.MustCompile(`[^\w\.]`)
)

func NormalizeUsername(input string) string {
	return rxNormalizeUsername.ReplaceAllString(
		strings.ToLowerSpecial(unicode.TurkishCase, input),
		"",
	)
}

func RemoveDots(input string) string {
	return strings.Replace(input, ".", "", -1)
}
