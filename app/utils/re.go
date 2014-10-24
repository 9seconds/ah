package utils

import (
	"errors"

	pcre "github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Regexp struct {
	exp pcre.Regexp
}

func (r *Regexp) Match(suspected string) bool {
	return r.exp.MatcherString(suspected, 0).Matches()
}

func (r *Regexp) Groups(suspected string) ([]string, error) {
	matcher := r.exp.MatcherString(suspected, 0)
	if !matcher.Matches() {
		return nil, errors.New("Cannot match")
	}

	groups := make([]string, matcher.Groups())
	for idx := 0; idx < len(groups); idx++ {
		groups[idx] = matcher.GroupString(idx+1)
	}

	return groups, nil
}

func CreateRegexp(expression string) *Regexp {
	return &Regexp{exp: pcre.MustCompile(expression, 0)}
}
