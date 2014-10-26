package utils

import (
	"fmt"
	"regexp"

	//re2 "code.google.com/p/sre2/sre2"
)

// Regexp is just a wrapper to play with different regular expression libraries
// to get the best one.
type Regexp struct {
	exp *regexp.Regexp
}

// Match verifies that suspected string matches regular expression.
func (r *Regexp) Match(suspected string) bool {
	return r.exp.MatchString(suspected)
}

// Groups returns a slice of a groups of regular expression applied to the
// string.
func (r *Regexp) Groups(suspected string) (groups []string, err error) {
	indexes := r.exp.FindStringSubmatchIndex(suspected)
	if indexes == nil {
		err = fmt.Errorf("Cannot get groups for %s", suspected)
	} else {
		groups = make([]string, len(indexes)/2-1)
		for idx := 0; idx < len(groups); idx++ {
			groups[idx] = suspected[indexes[2*idx+2]:indexes[2*idx+3]]
		}
	}
	return
}

// CreateRegexp creates a regexp with MustCompile (or equialent) method.
func CreateRegexp(expression string) *Regexp {
	return &Regexp{exp: regexp.MustCompile(expression)}
}
