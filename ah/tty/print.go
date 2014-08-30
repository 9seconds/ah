package tty

// --- Imports

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// --- Funcs

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

	lengths := fillLengths(content, shrinkStatus, padding)
	templateParts := createTemplateParts(lengths, padding)
	printTemplateParts(templateParts, content)
}

func fillLengths(content [][]string, shrinkStatus []bool, padding int) []int {
	lengths := make([]int, len(shrinkStatus), len(shrinkStatus))
	shrinkables := 0

	for _, columns := range content {
		for idx, col := range columns {
			columnLength := utf8.RuneCountInString(col)
			if lengths[idx] < columnLength {
				lengths[idx] = columnLength
			}
		}
	}

	for _, status := range shrinkStatus {
		if status {
			shrinkables += 1
		}
	}

	if shrinkables > 0 {
		shrinkableLength := TerminalWidth - padding*len(content[0])
		fmt.Println(shrinkableLength)
		for idx, length := range lengths {
			if !shrinkStatus[idx] {
				shrinkableLength -= length
			}
		}
		shrinkableLength /= shrinkables

		for idx := range lengths {
			if shrinkStatus[idx] {
				lengths[idx] = shrinkableLength
			}
		}
	}

	return lengths
}

func createTemplateParts(lengths []int, padding int) []string {
	emptyArray := make([]string, padding, padding)
	templateParts := make([]string, len(lengths), len(lengths))

	paddedTemplate := strings.Join(emptyArray, " ")

	for idx, value := range lengths {
		width := strconv.FormatInt(int64(value), 10)
		templateParts[idx] = paddedTemplate + "%-" + width + "." + width + "s"
	}
	templateParts[len(templateParts) - 1] += "\n"

	return templateParts
}

func printTemplateParts(templateParts []string, content [][]string) {
	for _, cols := range content {
		for idx, col := range cols {
			fmt.Printf(templateParts[idx], col)
		}
	}
}
