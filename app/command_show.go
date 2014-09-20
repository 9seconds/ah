package app

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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

func getCommands(filter *regexp.Regexp, env *Environment) ([]HistoryEntry, error) {
	historyChan := make(chan []int, 1)
	go getHistoryEntries(env, historyChan)

	var parse Parser
	if env.Shell == "bash" {
		parse = parseBash
	} else {
		parse = parseZsh
	}

	file, err := os.Open(env.HistFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if commands, err := parse(scanner, filter); err == nil {
		histories := <-historyChan
		for _, value := range histories {
			if len(commands) > value {
				commands[value-1].HasHistory = true
			}
		}
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

func getHistoryEntries(env *Environment, resultChan chan []int) {
	entries := make([]int, 0, 16)

	files, err := ioutil.ReadDir(env.GetTracesDir())
	if err != nil {
		resultChan <- entries
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if number, err := strconv.Atoi(file.Name()); err == nil && number >= 0 {
			entries = append(entries, number)
		}
	}

	resultChan <- entries
}

func CommandShow(slice *Slice, filter *regexp.Regexp, env *Environment) {
	commands, err := getCommands(filter, env)
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
