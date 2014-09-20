package app

import (
	"errors"
	"fmt"
	"path/filepath"
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

func (e *Environment) GetTracesDir() string {
	return filepath.Join(e.AppDir, "traces")
}

func (e *Environment) GetBookmarksDir() string {
	return filepath.Join(e.AppDir, "bookmarks")
}

func (e *Environment) GetTraceFileName(number int) string {
	return filepath.Join(e.GetTracesDir(), strconv.Itoa(number))
}

func (e *Environment) GetBookmarkFileName(name string) string {
	return filepath.Join(e.GetBookmarksDir(), name)
}

func ExtractSlice(single interface{}, start interface{}, finish interface{}) (*Slice, error) {
	slice := new(Slice)

	if single == nil && start == nil && finish == nil {
		slice.Finish = -1
		return slice, nil
	}

	if single != nil {
		singleStr := strings.Replace(single.(string), "_", "-", 1)
		if singleInt, err := strconv.Atoi(singleStr); err == nil {
			slice.Finish = singleInt
			return slice, nil
		} else {
			errToReturn := errors.New(fmt.Sprintf("Cannot convert lastNcommands to int: %v", err))
			return slice, errToReturn
		}
	}

	if start == nil || finish == nil {
		err := errors.New("Cannot process slice commands")
		return slice, err
	}
	startStr := strings.Replace(start.(string), "_", "-", 1)
	finishStr := strings.Replace(finish.(string), "_", "-", 1)

	startInt, startErr := strconv.Atoi(startStr)
	if startErr != nil {
		errToReturn := errors.New(fmt.Sprintf("Cannot process startFromNCommand: %v", startErr))
		return slice, errToReturn
	}
	finishInt, finishErr := strconv.Atoi(finishStr)
	if finishErr != nil {
		errToReturn := errors.New(fmt.Sprintf("Cannot process finishByMCommand: %v", startErr))
		return slice, errToReturn
	}

	slice.Start = startInt
	slice.Finish = finishInt

	return slice, nil

}
