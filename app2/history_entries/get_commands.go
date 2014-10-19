package history_entries

import (
	"bufio"
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"

	"../environments"
	"../utils"
)

func GetCommands(filter *regexp.Regexp, env *environments.Environment) ([]HistoryEntry, error) {
	if !env.OK() {
		return nil, errors.New("Environment is not prepared")
	}

	historyChan := getHistoryEntriesChan(env)

	histFile, _ := env.GetHistFile()
	file := utils.Open(histFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if commands, err := getParser(env)(scanner, filter); err == nil {
		histories := <-historyChan
		for _, number := range histories {
			if len(commands) > number {
				commands[number-1].hasHistory = true
			}
		}
		return commands, nil
	} else {
		return nil, err
	}
}

func getHistoryEntriesChan(env *environments.Environment) chan []int {
	historyChan := make(chan []int, 1)

	go func() {
		entries := make([]int, 0, 16)

		files, err := ioutil.ReadDir(env.GetTracesDir())
		if err != nil {
			historyChan <- entries
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if number, err := strconv.Atoi(file.Name()); err == nil && number >= 0 {
				entries = append(entries, number)
			}
		}

		historyChan <- entries
	}()

	return historyChan
}
