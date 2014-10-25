package commands

import (
	"os"

	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
)

// RemoveBookmarks removes the list of bookmarks from the storage.
func RemoveBookmarks(bookmarks []string, env *environments.Environment) {
	logger, _ := env.GetLogger()

	for _, bookmark := range bookmarks {
		fileName := env.GetBookmarkFileName(bookmark)
		err := os.Remove(fileName)

		if err == nil {
			logger.WithFields(logrus.Fields{
				"filename": fileName,
			}).Info("File was deleted")
		} else {
			logger.WithFields(logrus.Fields{
				"filename": fileName,
				"error":    err,
			}).Warn("File was not deleted")
		}
	}
}
