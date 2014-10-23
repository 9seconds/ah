package utils

import (
	"fmt"

	re2 "code.google.com/p/sre2/sre2"
)

type Regexp struct {
	exp re2.Re
}

func (r *Regexp) Match(suspected string) bool {
	return r.exp.Match(suspected)
}

func (r *Regexp) Groups(suspected string) ([]string, error) {
	indexes := r.exp.MatchIndex(suspected)
	if indexes == nil {
		return nil, fmt.Errorf("Cannot get groups for %s", suspected)
	}

	groups := make([]string, len(indexes)/2)
	for idx := 0; idx < len(groups); idx++ {
		groups[idx] = suspected[indexes[2*idx]:indexes[2*idx+1]]
	}

	return groups, nil
}

func CreateRegexp(expression string) *Regexp {
	return &Regexp{exp: re2.MustParse(expression)}
}
