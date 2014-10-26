package slices

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Slice is just a pair of start/finish numbers. Used for show command.
type Slice struct {
	Start  int
	Finish int
}

// GetSliceIndex returns fixed index based on a given length. It is required if you
// want to support negative indexing. I want to support it.
func GetSliceIndex(index int, length int) int {
	if index >= 0 {
		return index
	}
	return length + index
}

// ExtractSlice extracts slice from the given arguments. Please check how show command
// looks like to get an idea what will it construct.
func ExtractSlice(single interface{}, start interface{}, finish interface{}) (*Slice, error) {
	slice := new(Slice)

	if single == nil && start == nil && finish == nil {
		return extractNils()
	}

	if single != nil {
		return extractSingle(single.(string))
	}

	if start == nil || finish == nil {
		err := errors.New("Cannot process slice commands")
		return slice, err
	}

	return extractStartFinish(start.(string), finish.(string))
}

func extractNils() (*Slice, error) {
	return &Slice{Start: 0, Finish: -1}, nil
}

func extractSingle(single string) (slice *Slice, err error) {
	converted, err := strconv.Atoi(convertSubstituteToMinus(single))
	if err == nil {
		slice = &Slice{Start: -converted - 1, Finish: -1}
	} else {
		err = fmt.Errorf("Cannot convert %v to integer", err)
	}
	return
}

func extractStartFinish(start string, finish string) (slice *Slice, err error) {
	slice = new(Slice)

	converted, err := parseInt(convertSubstituteToMinus(start))
	if err != nil {
		return
	}
	slice.Start = converted

	converted, err = parseInt(convertSubstituteToMinus(finish))
	if err != nil {
		return
	}
	slice.Finish = converted

	return
}

func parseInt(str string) (converted int, err error) {
	converted, err = strconv.Atoi(str)
	if err != nil {
		err = fmt.Errorf("Cannot convert %v to integer", err)
	}
	return
}

func convertSubstituteToMinus(str string) string {
	return strings.Replace(str, "_", "-", 1)
}
