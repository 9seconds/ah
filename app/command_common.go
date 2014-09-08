package app

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var sliceRegexp = regexp.MustCompile(`^(_?\d+)(:(_?\d+))?$`)

type Slice struct {
	Start  int
	Finish int
}

type Environment struct {
	AppDir         string
	HistFile       string
	HistTimeFormat string
	Shell          string
}

func ExtractSlice(argument string) (*Slice, error) {
	slice := new(Slice)
	if argument == "" {
		slice.Finish = -1
		return slice, nil
	}

	matcher := sliceRegexp.FindStringSubmatch(argument)
	if len(matcher) == 0 {
		return nil, errors.New("Incorrect slice format")
	}
	for idx, v := range matcher {
		matcher[idx] = strings.Replace(v, "_", "-", 1)
	}
	if matcher[1] != "" {
		if int_value, result := strconv.Atoi(matcher[1]); result == nil {
			slice.Finish = -1
			slice.Start = -int_value
		} else {
			return nil, result
		}
	}
	if matcher[3] != "" {
		if int_value, result := strconv.Atoi(matcher[3]); result == nil {
			slice.Start = -slice.Start
			slice.Finish = int_value
		} else {
			return nil, result
		}
	}

	return slice, nil
}
