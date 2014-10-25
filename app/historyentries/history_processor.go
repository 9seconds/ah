package historyentries

import (
	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

func processHistories(env *environments.Environment) (chan bool, chan *HistoryEntry) {
	resultChan := make(chan bool, 1)
	consumeChan := make(chan *HistoryEntry, historyEventsCapacity)

	go func() {
		entries := make(map[string]bool)

		files, err := env.GetTraceFilenames()
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Error on traces directory listing")
			resultChan <- true
			return
		}

		for _, file := range files {
			entries[file.Name()] = true
		}
		utils.Logger.WithField("filenames", entries).Info("Parsed filenames")

		for entry := range consumeChan {
			if _, found := entries[entry.GetTraceName()]; found {
				entry.hasHistory = true
			}
		}
		resultChan <- true
	}()

	return resultChan, consumeChan
}
