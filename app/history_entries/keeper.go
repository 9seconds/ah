package history_entries

import (
	"os"
	"strconv"
)

var historyEventsCapacity = 5000

func init() {
	histFileSize := os.Getenv("HISTFILESIZE")
	if histFileSize != "" {
		if converted, err := strconv.Atoi(histFileSize); err == nil && converted > 0 {
			historyEventsCapacity = converted
		}
	}
}

type Keeper interface {
	Init() *HistoryEntry
	Commit(*HistoryEntry, chan *HistoryEntry) *HistoryEntry
	Continue() bool
	Result() interface{}
}

type singleKeeper struct {
	current *HistoryEntry
}

type preciseNumberKeeper struct {
	singleKeeper
	currentIndex  int
	preciseNumber int
}

type allKeeper struct {
	currentIndex int
	entries      []HistoryEntry
}

type rangeKeeper struct {
	current      *HistoryEntry
	currentIndex int
	start        int
	finish       int
	entries      []HistoryEntry
}

func (sk *singleKeeper) Init() *HistoryEntry {
	sk.current = new(HistoryEntry)
	return sk.current
}

func (sk *singleKeeper) Commit(event *HistoryEntry, historyChannel chan *HistoryEntry) *HistoryEntry {
	return sk.current
}

func (sk *singleKeeper) Continue() bool {
	return true
}

func (sk *singleKeeper) Result() interface{} {
	return *sk.current
}

func (pnk *preciseNumberKeeper) Commit(event *HistoryEntry, historyChannel chan *HistoryEntry) *HistoryEntry {
	pnk.currentIndex++
	return pnk.singleKeeper.Commit(event, historyChannel)
}

func (pnk *preciseNumberKeeper) Continue() bool {
	return pnk.currentIndex <= pnk.preciseNumber
}

func (pnk *preciseNumberKeeper) SetNumber(number int) {
	pnk.preciseNumber = number
}

func (ak *allKeeper) Init() *HistoryEntry {
	ak.entries = make([]HistoryEntry, historyEventsCapacity)
	return &ak.entries[0]
}

func (ak *allKeeper) Commit(event *HistoryEntry, historyChannel chan *HistoryEntry) *HistoryEntry {
	historyChannel <- event
	ak.currentIndex++
	if ak.currentIndex == len(ak.entries) {
		ak.entries = append(ak.entries, HistoryEntry{})
	}
	return &ak.entries[ak.currentIndex]
}

func (ak *allKeeper) Continue() bool {
	return true
}

func (ak *allKeeper) Result() interface{} {
	return ak.entries[:ak.currentIndex]
}

func (rk *rangeKeeper) SetLimits(start, finish int) {
	rk.start = start
	rk.finish = finish
}

func (rk *rangeKeeper) Init() *HistoryEntry {
	rk.entries = make([]HistoryEntry, rk.finish-rk.start)
	rk.current = new(HistoryEntry)
	return rk.current
}

func (rk *rangeKeeper) Commit(event *HistoryEntry, historyChannel chan *HistoryEntry) *HistoryEntry {
	historyChannel <- event
	if rk.start <= rk.currentIndex && rk.currentIndex < rk.finish {
		rk.entries[rk.currentIndex-rk.start] = *rk.current
		rk.current = new(HistoryEntry)
		rk.currentIndex++
	}
	return rk.current
}

func (rk *rangeKeeper) Continue() bool {
	return rk.currentIndex < rk.finish
}

func (rk *rangeKeeper) Result() interface{} {
	return rk.entries[rk.currentIndex-rk.start]
}

func getKeeper(mode GetCommandsMode, varargs ...int) Keeper {
	switch mode {
	case GET_COMMANDS_ALL:
		return new(allKeeper)
	case GET_COMMANDS_RANGE:
		keeper := new(rangeKeeper)
		keeper.SetLimits(varargs[0], varargs[1])
		return keeper
	case GET_COMMANDS_SINGLE:
		return new(singleKeeper)
	case GET_COMMANDS_PRECISE:
		keeper := new(preciseNumberKeeper)
		keeper.SetNumber(varargs[0])
		return keeper
	}
	panic("Unknown GetCommandsMode")
}
