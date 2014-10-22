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
	}
	return length + index
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
	converted, err := strconv.Atoi(convertSubstituteToMinus(single))
	if err == nil {
		slice := Slice{Start: -converted - 1, Finish: -1}
		return &slice, nil
	}
	return nil, fmt.Errorf("Cannot convert %v to integer", err)
}

func extractStartFinish(start string, finish string) (*Slice, error) {
	slice := new(Slice)

	converted, err := parseInt(convertSubstituteToMinus(start))
	if err != nil {
		return nil, err
	}
	slice.Start = converted

	converted, err = parseInt(convertSubstituteToMinus(finish))
	if err != nil {
		return nil, err
	}
	slice.Finish = converted

	return slice, nil
}

func parseInt(str string) (int, error) {
	converted, err := strconv.Atoi(str)
	if err == nil {
		return converted, nil
	}
	return 0, fmt.Errorf("Cannot convert %v to integer", err)
}

func convertSubstituteToMinus(str string) string {
	return strings.Replace(str, "_", "-", 1)
}
