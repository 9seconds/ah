package slices

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Slice struct {
	Start  int
	Finish int
}

func GetSliceIndex(index int, length int) int {
	if index >= 0 {
		return index
	} else {
		return length + index
	}
}

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

func extractSingle(single string) (*Slice, error) {
	if converted, err := strconv.Atoi(convertSubstituteToMinus(single)); err == nil {
		slice := Slice{Start: -converted - 1, Finish: -1}
		return &slice, nil
	} else {
		return nil, fmt.Errorf("Cannot convert %v to integer", err)
	}
}

func extractStartFinish(start string, finish string) (*Slice, error) {
	slice := new(Slice)

	if converted, err := parseInt(convertSubstituteToMinus(start)); err != nil {
		return nil, err
	} else {
		slice.Start = converted
	}

	if converted, err := parseInt(convertSubstituteToMinus(finish)); err != nil {
		return nil, err
	} else {
		slice.Finish = converted
	}

	return slice, nil
}

func parseInt(str string) (int, error) {
	if converted, err := strconv.Atoi(str); err == nil {
		return converted, nil
	} else {
		return 0, fmt.Errorf("Cannot convert %v to integer", err)
	}
}

func convertSubstituteToMinus(str string) string {
	return strings.Replace(str, "_", "-", 1)
}
