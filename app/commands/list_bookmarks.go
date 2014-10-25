package commands

import (
	"fmt"
	"io/ioutil"
	"strconv"

	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

// ListBookmarks prints the list of bookmarks with their content
func ListBookmarks(env *environments.Environment) {
	bookmarksFileInfos, err := env.GetBookmarkFilenames()
	if err != nil {
		utils.Logger.Panic(err)
	}

	maxLength := 1
	for _, fileInfo := range bookmarksFileInfos {
		name := fileInfo.Name()
		if len(name) > maxLength {
			maxLength = len(name)
		}
	}

	template := "%-" + strconv.Itoa(maxLength) + "s    %s\n"
	utils.Logger.WithFields(logrus.Fields{
		"template": template,
	}).Info("Calculated template to print")

	for _, fileInfo := range bookmarksFileInfos {
		fileName := fileInfo.Name()

		content, err := ioutil.ReadFile(env.GetBookmarkFileName(fileName))
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"filename": fileName,
				"error":    err,
			}).Warn("Cannot read a content of the file so skip")
			continue
		}

		fmt.Printf(template, fileName, string(content))
	}
}
