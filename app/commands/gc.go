package commands

import (
	"os"
	"sort"
	"time"

	"../environments"
)

type GcType uint8

const (
	GC_KEEP_LATEST = iota
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
	return fis.content[i].ModTime().Unix() < fis.content[i].ModTime().Unix()
}

func (fis FileInfoSorter) Swap(i, j int) {
	fis.content[i], fis.content[j] = fis.content[j], fis.content[i]
}

func (fis FileInfoSorter) YoungerThan(timestamp int64) FileInfoSorter {
	binarySearchFunc := func(i int) bool {
		return fis.content[i].ModTime().Unix() > timestamp
	}

	index := sort.Search(len(fis.content), binarySearchFunc)
	if index == len(fis.content) {
		return FileInfoSorter{content: []os.FileInfo{}}
	}
	return FileInfoSorter{content: fis.content[index:]}
}

func (fis FileInfoSorter) First(count int) FileInfoSorter {
	if count >= len(fis.content) {
		return FileInfoSorter{content: []os.FileInfo{}}
	}
	return FileInfoSorter{content: fis.content[:count]}
}

func GC(gcType GcType, param int, env *environments.Environment) {
	fileInfos, err := env.GetTraceFilenames()
	if err != nil {
		panic("Cannot fetch the list of trace filenames")
	}
	fileInfoSorter := FileInfoSorter{content: fileInfos}
	sort.Sort(fileInfoSorter)

	if gcType == GC_KEEP_LATEST {
		fileInfoSorter = fileInfoSorter.First(param)
	} else {
		timestamp := time.Now().Unix() - int64(param*SECONDS_IN_DAY)
		fileInfoSorter = fileInfoSorter.YoungerThan(timestamp)
	}

	for _, info := range fileInfoSorter.content {
		os.Remove(info.Name())
	}
}
