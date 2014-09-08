package app

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/weidewang/go-strftime"
)

type HistoryEntry struct {
	Number     int
	Command    string
	Timestamp  time.Time
	HasHistory bool
}

type Parser func(scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error)

func getCommands(parse Parser, filter *regexp.Regexp, env *Environment) ([]HistoryEntry, error) {
	file, err := os.Open(env.HistFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if commands, err := parse(scanner, filter); err == nil {
		return commands, nil
	} else {
		return nil, err
	}
}

func getSliceIndex(index int, length int) int {
	if index >= 0 {
		return index
	} else {
		return length + index
	}
}

func CommandShow(slice *Slice, filter *regexp.Regexp, env *Environment) {
	var parse Parser
	if env.Shell == "bash" {
		parse = parseBash
	} else {
		parse = parseZsh
	}

	commands, err := getCommands(parse, filter, env)
	sliceStart := getSliceIndex(slice.Start, len(commands))
	sliceFinish := getSliceIndex(slice.Finish, len(commands))

	if err != nil || sliceStart < 0 || sliceFinish < 0 || sliceFinish <= sliceStart {
		return
	}

	for _, value := range commands[sliceStart:sliceFinish] {
		fmt.Println(getPrintableCommand(&value, env))
	}
}

func getPrintableCommand(historyEntry *HistoryEntry, env *Environment) string {
	historyMark := ' '
	if historyEntry.HasHistory {
		historyMark = '*'
	}

	timestamp := ""
	if env.HistTimeFormat != "" {
		timestamp = "\t" + strftime.Strftime(&historyEntry.Timestamp,
			env.HistTimeFormat)
	}

	return fmt.Sprintf("!%d %c%s\t%s",
		historyEntry.Number, historyMark, timestamp, historyEntry.Command)
}
