package utils

import (
	"errors"

	pcre "github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// Regexp is just a wrapper to play with different regular expression libraries
// to get the best one.
type Regexp struct {
	exp pcre.Regexp
}

// Match verifies that suspected string matches regular expression.
func (r *Regexp) Match(suspected string) bool {
	return r.exp.MatcherString(suspected, 0).Matches()
}

// Groups returns a slice of a groups of regular expression applied to the
// string.
func (r *Regexp) Groups(suspected string) ([]string, error) {
	matcher := r.exp.MatcherString(suspected, 0)
	if !matcher.Matches() {
		return nil, errors.New("Cannot match")
	}

	groups := make([]string, matcher.Groups())
	for idx := 0; idx < len(groups); idx++ {
		groups[idx] = matcher.GroupString(idx + 1)
	}

	return groups, nil
}

// CreateRegexp creates a regexp with MustCompile (or equialent) method.
func CreateRegexp(expression string) *Regexp {
	return &Regexp{exp: pcre.MustCompile(expression, 0)}
}
