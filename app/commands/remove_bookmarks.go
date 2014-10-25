package commands

import (
	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

// RemoveBookmarks removes the list of bookmarks from the storage.
func RemoveBookmarks(bookmarks []string, env *environments.Environment) {
	for _, bookmark := range bookmarks {
		utils.RemoveWithLogging(env.GetBookmarkFileName(bookmark))
	}
}
