package historyentries

import (
	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
)

func processHistories(env *environments.Environment) (chan bool, chan *HistoryEntry) {
	resultChan := make(chan bool, 1)
	consumeChan := make(chan *HistoryEntry, historyEventsCapacity)

	go func() {
		entries := make(map[string]bool)
		logger, _ := env.GetLogger()

		files, err := env.GetTraceFilenames()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Error on traces directory listing")
			resultChan <- true
			return
		}

		for _, file := range files {
			entries[file.Name()] = true
		}

		for entry := range consumeChan {
			if _, found := entries[entry.GetTraceName()]; found {
				entry.hasHistory = true
			}
		}
		resultChan <- true
	}()

	return resultChan, consumeChan
}
