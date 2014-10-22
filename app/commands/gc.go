package commands

import (
	"os"
	"sort"
	"time"

	logrus "github.com/Sirupsen/logrus"

	"../environments"
)

type GcType uint8

const (
	GC_ALL GcType = iota
	GC_KEEP_LATEST
	GC_OLDER_THAN
)

const SECONDS_IN_DAY = 60 * 60 * 24

type FileInfoSorter struct {
	content []os.FileInfo
}

func (fis FileInfoSorter) Len() int {
	return len(fis.content)
}

func (fis FileInfoSorter) Less(i, j int) bool {
	return fis.content[i].ModTime().Unix() < fis.content[j].ModTime().Unix()
}

func (fis FileInfoSorter) Swap(i, j int) {
	fis.content[i], fis.content[j] = fis.content[j], fis.content[i]
}

func (fis FileInfoSorter) YoungerThan(timestamp int64) []os.FileInfo {
	binarySearchFunc := func(i int) bool {
		return fis.content[i].ModTime().Unix() > timestamp
	}
	index := sort.Search(len(fis.content), binarySearchFunc)
	return fis.content[:index]
}

func (fis FileInfoSorter) Tail(first int) []os.FileInfo {
	if first >= len(fis.content) {
		return fis.content
	}
	return fis.content[len(fis.content)-first:]
}

func GC(gcType GcType, param int, env *environments.Environment) {
	logger, _ := env.GetLogger()
	fileInfos, err := env.GetTraceFilenames()
	if err != nil {
		panic("Cannot fetch the list of trace filenames")
	}
	fileInfoSorter := FileInfoSorter{content: fileInfos}
	sort.Sort(fileInfoSorter)

	switch gcType {
	case GC_KEEP_LATEST:
		fileInfos = fileInfoSorter.Tail(param)
	case GC_OLDER_THAN:
		timestamp := time.Now().Unix() - SECONDS_IN_DAY*int64(param)
		fileInfos = fileInfoSorter.YoungerThan(timestamp)
	default:
		fileInfos = fileInfoSorter.content
	}

	for _, info := range fileInfos {
		logger.WithFields(logrus.Fields{
			"filename": info.Name(),
		}).Info("Remove file")
		os.Remove(env.GetTraceFileName(info.Name()))
	}
}
