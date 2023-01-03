package main

import (
	"regexp"
	"strings"
)

// Take a source and tokenize it.
var replace_re = regexp.MustCompile(`\s\s+`)

func StripSource(source string) string {
	x := make([]string, 0)
	for _, c := range source {
		switch c {
		case ' ', '\t', '\r', '\n', '\\', '$', '(', ')', '{', '}', '[', ']', ',':
			x = append(x, " ")
		case '+', '-', '*', '/', '>', '<':
			x = append(x, " "+string(c)+" ")
		default:
			x = append(x, string(c))
		}
	}
	return replace_re.ReplaceAllString(strings.Join(x, ""), " ")
}
