package main

import (
	"regexp"
	"strings"
)

func minify(css string) string {
	whitespace := regexp.MustCompile("[ \t\n\r\f\v]+")
	css = whitespace.ReplaceAllString(css, " ")

	comments := regexp.MustCompile("/\\*[^*]*\\*/") // TODO
	css = comments.ReplaceAllString(css, " ")

	css = strings.Replace(css, "; ", ";", -1)
	css = strings.Replace(css, ": ", ":", -1)
	css = strings.Replace(css, " {", "{", -1)
	css = strings.Replace(css, "{ ", "{", -1)
	css = strings.Replace(css, ", ", ",", -1)
	css = strings.Replace(css, "} ", "}", -1)
	css = strings.Replace(css, ";}", "}", -1)
	css = strings.TrimSpace(css)
	return css
}
