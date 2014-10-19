package history_entries

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	logrus "github.com/Sirupsen/logrus"

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

	if commands, err := getParser(env)(env, scanner, filter); err == nil {
		histories := <-historyChan
		fmt.Println(histories)
		for idx, _ := range commands {
			if _, ok := histories[commands[idx].number]; ok {
				commands[idx].hasHistory = true
			}
		}
		return commands, nil
	} else {
		return nil, err
	}
}

func getHistoryEntriesChan(env *environments.Environment) chan map[uint]bool {
	historyChan := make(chan map[uint]bool, 1)

	go func() {
		entries := make(map[uint]bool)
		logger, _ := env.GetLogger()

		files, err := ioutil.ReadDir(env.GetTracesDir())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Error on traces directory listing")
			historyChan <- entries
		}

		for _, file := range files {
			if file.IsDir() {
				logger.WithFields(logrus.Fields{
					"filename": file.Name(),
				}).Info("Skip file because it is directory")
				continue
			}
			if number, err := strconv.Atoi(file.Name()); err == nil && number >= 0 {
				logger.WithFields(logrus.Fields{
					"number": number,
				}).Debug("Add history trace to the list of entries")
				entries[uint(number)] = true
			} else {
				logger.WithFields(logrus.Fields{
					"error":  err,
					"number": number,
				}).Warn("Cannot add trace to the list of entries")
			}
		}

		historyChan <- entries
	}()

	return historyChan
}
