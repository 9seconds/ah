package utils


import (
	"fmt"
	"unicode/utf8"
	"strings"
	"strconv"
)


// Gentle TTY printer which undestands how to deal with different
// terminal widths

// Basic usage is rather simple: just give it an array of rows (which is
// an array of strings), set padding and determine which columns may be
// shrinked.

// Example:

// 	content := [][]string{
// 		[]string{"1", "/usr/bin/ls", "12309"},
// 		[]string{"12", "/usr/bin/rsync", "12309"},
// 		[]string{"13", "/usr/bin/python", "12308"},
// 		[]string{"14", "/usr/bin/env", "12307"},
// 	}
// 	shrinkStatus := []bool{false, true, false}

// 	utils.PrintFormatted(content[:], shrinkStatus, 4)
func PrintFormatted(content [][]string, shrinkStatus []bool, padding int) {
	if len(content) == 0 {
		return
	}

	maxLengths := make([]int, len(shrinkStatus), len(shrinkStatus))
	shrinkables := 0

	for _, columns := range content {
		for idx, col := range columns {
			columnLength := utf8.RuneCountInString(col)
			if maxLengths[idx] < columnLength {
				maxLengths[idx] = columnLength
			}
		}
	}
	for _, status := range shrinkStatus {
		if status {
			shrinkables += 1
		}
	}
	if shrinkables > 0 {
		shrinkableLength := TerminalWidth - padding * len(content[0])
		fmt.Println(shrinkableLength)
		for idx, length := range maxLengths {
			if !shrinkStatus[idx] {
				shrinkableLength -= length
			}
		}
        shrinkableLength /= shrinkables

		for idx := range maxLengths {
			if shrinkStatus[idx] {
				maxLengths[idx] = shrinkableLength
			}
		}
	}

	emptyArray := make([]string, padding, padding)
	paddedTemplate := strings.Join(emptyArray, " ")

	templateParts := make([]string, len(maxLengths), len(maxLengths))
	for idx, value := range maxLengths {
		width := strconv.FormatInt(int64(value), 10)
		templateParts[idx] = paddedTemplate + "%-" + width + "." + width + "s"
	}
	templateParts[len(templateParts)-1] += "\n"

	for _, cols := range content {
		for idx, col := range cols {
			fmt.Printf(templateParts[idx], col)
		}
	}
}
