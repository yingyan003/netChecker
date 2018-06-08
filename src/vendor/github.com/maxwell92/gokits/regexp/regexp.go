package regexp

import (
	"regexp"
)

var Match = regexp.MustCompile

func Expression(res ...*regexp.Regexp) *regexp.Regexp {
	var s string
	for _, re := range res {
		s += re.String()
	}

	return Match(s)
}

func Optional(res ...*regexp.Regexp) *regexp.Regexp {
	return Match(Group(Expression(res...)).String() + `?`)
}

func Repeated(res ...*regexp.Regexp) *regexp.Regexp {
	return Match(Group(Expression(res...)).String() + `+`)
}

func Group(res ...*regexp.Regexp) *regexp.Regexp {
	return Match(`(?:` + Expression(res...).String() + `)`)
}
