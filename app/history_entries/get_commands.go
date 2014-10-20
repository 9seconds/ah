package history_entries

import (
	"bufio"
	"errors"
	"io/ioutil"
	"regexp"

	logrus "github.com/Sirupsen/logrus"

	"../environments"
	"../utils"
)

func GetCommands(filter *regexp.Regexp, env *environments.Environment) ([]*HistoryEntry, error) {
	if !env.OK() {
		return nil, errors.New("Environment is not prepared")
	}

	resultChan, consumeChan := processHistories(env)

	histFile, _ := env.GetHistFile()
	file := utils.Open(histFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if commands, err := getParser(env)(env, scanner, filter, consumeChan); err == nil {
		<-resultChan
		return commands, nil
	} else {
		return nil, err
	}
}

func processHistories(env *environments.Environment) (chan bool, chan *HistoryEntry) {
	resultChan := make(chan bool, 1)
	consumeChan := make(chan *HistoryEntry, historyEventsCapacity)

	go func() {
		entries := make(map[string]bool)
		logger, _ := env.GetLogger()

		files, err := ioutil.ReadDir(env.GetTracesDir())
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Error on traces directory listing")
			resultChan <- true
			return
		}

		for _, file := range files {
			if file.IsDir() {
				logger.WithFields(logrus.Fields{
					"filename": file.Name(),
				}).Info("Skip file because it is directory")
				continue
			}
			entries[file.Name()] = true
		}

		for {
			entry, ok := <-consumeChan
			if ok {
				if _, found := entries[entry.GetTraceName()]; found {
					entry.hasHistory = true
				}
			} else {
				resultChan <- true
			}
		}
	}()

	return resultChan, consumeChan
}
