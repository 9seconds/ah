package commands

import (
	"os"
	"sort"
	"time"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

// GcType is the type of GC to execute.
type GcType uint8

// Types of garbage collecting.
const (
	GcAll GcType = iota
	GcKeepLatest
	GcOlderThan
)

// GcDir is the directory where GC has to be performed.
type GcDir uint8

// Directories for garbage collecting.
const (
	GcTracesDir GcDir = iota
	GcBookmarksDir
)

const secondsInDay = 60 * 60 * 24

type fileInfoSorter struct {
	content []os.FileInfo
}

func (fis fileInfoSorter) Len() int {
	return len(fis.content)
}

func (fis fileInfoSorter) Less(i, j int) bool {
	return fis.content[i].ModTime().Unix() < fis.content[j].ModTime().Unix()
}

func (fis fileInfoSorter) Swap(i, j int) {
	fis.content[i], fis.content[j] = fis.content[j], fis.content[i]
}

func (fis fileInfoSorter) YoungerThan(timestamp int64) []os.FileInfo {
	binarySearchFunc := func(i int) bool {
		return fis.content[i].ModTime().Unix() > timestamp
	}
	index := sort.Search(len(fis.content), binarySearchFunc)
	return fis.content[:index]
}

func (fis fileInfoSorter) Tail(first int) []os.FileInfo {
	if first >= len(fis.content) {
		return fis.content
	}
	return fis.content[len(fis.content)-first:]
}

// GC implements g (garbage collecting) command.
func GC(gcType GcType, gcDir GcDir, param int, env *environments.Environment) {
	listFunction := env.GetTracesFileInfos
	fileNameFunction := env.GetTraceFileName
	if gcDir == GcBookmarksDir {
		listFunction = env.GetBookmarksFileInfos
		fileNameFunction = env.GetBookmarkFileName
	}

	fileInfos, err := listFunction()
	if err != nil {
		utils.Logger.Panic("Cannot fetch the list of trace filenames")
	}
	infoSorter := fileInfoSorter{content: fileInfos}
	sort.Sort(infoSorter)

	switch gcType {
	case GcKeepLatest:
		fileInfos = infoSorter.Tail(param)
	case GcOlderThan:
		timestamp := time.Now().Unix() - secondsInDay*int64(param)
		fileInfos = infoSorter.YoungerThan(timestamp)
	default:
		fileInfos = infoSorter.content
	}

	for _, info := range fileInfos {
		utils.RemoveWithLogging(fileNameFunction(info.Name()))
	}
}
