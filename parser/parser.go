package parser

import (
	"regexp"
	"strings"
)

var (
	DefaultRegex = regexp.MustCompile("\\%(([^\\s\\\\%]|[ ])+)\\%")
)

type Parser struct {
	regex *regexp.Regexp
	f     func(string) string
}

func NewParser(regex *regexp.Regexp, f func(string) string) *Parser {
	return &Parser{
		regex: regex,
		f:     f,
	}
}

func (p *Parser) Parse(str string) string {
	str = p.regex.ReplaceAllStringFunc(str, func(s string) string {
		s = strings.ReplaceAll(s, "%", "")
		return p.f(s)
	})

	return minify(str)
}

func minify(str string) string {
	lines := strings.Split(str, "\n")
	empty := false
	for i := 0; i < len(lines); i++ {
		if strings.Trim(lines[i], "\t ") == "" {
			if empty {
				lines[i] = ""
				continue
			}

			empty = true
		} else {
			empty = false
		}
	}

	str = strings.Join(lines, "\n")
	str = regexp.MustCompile("\n{2,}").ReplaceAllString(str, "")
	return str
}
