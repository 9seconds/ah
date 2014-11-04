package historyentries

import (
	"bufio"
	"strconv"
	"strings"

	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

var (
	bashTimestampRegexp = utils.CreateRegexp(`^#\s*\d+$`)

	zshLineRegexp = utils.CreateRegexp(`^: (\d+):\d;(.*?)$`)

	fishCmdRegexp  = utils.CreateRegexp(`^- cmd:\s*(.*?)$`)
	fishWhenRegexp = utils.CreateRegexp(`\s*when:\s*(\d+)$`)
)

type (
	// Parser is a signature for a function which parses file and returns a Keeper.
	Parser func(Keeper, *bufio.Scanner, *utils.Regexp, chan *HistoryEntry) (Keeper, error)

	// ShellSpecificParser is a signature for a function which implements shell specific logic for parsing.
	ShellSpecificParser func(Keeper, string, uint, *HistoryEntry, *utils.Regexp, chan *HistoryEntry) (bool, uint, *HistoryEntry)
)

func getParser(env *environments.Environment) Parser {
	var currentNumber uint
	var shellSpecific ShellSpecificParser

	switch env.GetShell() {
	case environments.ShellBash:
		shellSpecific = parseBash
		currentNumber = 1
	case environments.ShellZsh:
		shellSpecific = parseZsh
		currentNumber = 0
	case environments.ShellFish:
		shellSpecific = parseFish
		currentNumber = 0
	default:
		utils.Logger.Panicf("Unknown shell %v", env.GetShell())
	}

	return func(keeper Keeper, scanner *bufio.Scanner, filter *utils.Regexp, historyChan chan *HistoryEntry) (Keeper, error) {
		defer close(historyChan)

		continueToConsume := false
		currentEvent := keeper.Init()

		for keeper.Continue() && scanner.Scan() {
			text := scanner.Text()

			utils.Logger.WithFields(logrus.Fields{
				"text":              text,
				"continueToConsume": continueToConsume,
				"currentEvent":      currentEvent,
			}).Info("Parse history line")

			if continueToConsume {
				utils.Logger.Info("Attach the line to the previous command")

				currentEvent.command += "\n" + text
				if strings.HasSuffix(text, `\`) {
					continue
				}
				continueToConsume = false
				utils.Logger.WithFields(logrus.Fields{
					"event": currentEvent,
				}).Info("Commit event")
				currentEvent = keeper.Commit(currentEvent, historyChan)
			}

			if text == "" {
				utils.Logger.Info("Skip empty line")
				continue
			}

			continueToConsume, currentNumber, currentEvent = shellSpecific(
				keeper,
				text,
				currentNumber,
				currentEvent,
				filter,
				historyChan)
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return keeper, nil
	}
}

func parseBash(keeper Keeper, text string, currentNumber uint, currentEvent *HistoryEntry, filter *utils.Regexp, historyChan chan *HistoryEntry) (bool, uint, *HistoryEntry) {
	continueToConsume := false
	if bashTimestampRegexp.Match(text) {
		if converted, err := strconv.ParseInt(text[1:], 10, 64); err == nil {
			utils.Logger.WithFields(logrus.Fields{
				"timestamp": converted,
			}).Info("Parse timestamp")
			currentEvent.timestamp = converted
		} else {
			utils.Logger.WithFields(logrus.Fields{
				"text":  text,
				"error": err,
			}).Warn("Cannot parse timestamp")
		}
	} else {
		if filter == nil || filter.Match(text) {
			currentEvent.command = text
			currentEvent.number = currentNumber

			continueToConsume = strings.HasSuffix(text, "\\")
			if !continueToConsume {
				utils.Logger.WithFields(logrus.Fields{
					"event": currentEvent,
				}).Info("Commit event")
				currentEvent = keeper.Commit(currentEvent, historyChan)
			}
		} else {
			utils.Logger.Info("Skip text line because of the filter.")
		}
		currentNumber++
	}

	return continueToConsume, currentNumber, currentEvent
}

func parseZsh(keeper Keeper, text string, currentNumber uint, currentEvent *HistoryEntry, filter *utils.Regexp, historyChan chan *HistoryEntry) (bool, uint, *HistoryEntry) {
	continueToConsume := false
	groups, err := zshLineRegexp.Groups(text)

	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Cannot parse current line, skip.")
		return continueToConsume, currentNumber, currentEvent
	}
	timestamp, command := groups[0], groups[1]
	currentNumber++

	if filter != nil && !filter.Match(command) {
		utils.Logger.Info("Skip text line because of the filter.")
		return continueToConsume, currentNumber, currentEvent
	}

	converted, _ := strconv.ParseInt(timestamp, 10, 64)
	currentEvent.command = command
	currentEvent.number = currentNumber
	currentEvent.timestamp = converted

	continueToConsume = strings.HasSuffix(text, `\`)
	if !continueToConsume {
		utils.Logger.WithFields(logrus.Fields{
			"event": currentEvent,
		}).Info("Commit event")
		currentEvent = keeper.Commit(currentEvent, historyChan)
	}

	return continueToConsume, currentNumber, currentEvent
}

func parseFish(keeper Keeper, text string, currentNumber uint, currentEvent *HistoryEntry, filter *utils.Regexp, historyChan chan *HistoryEntry) (bool, uint, *HistoryEntry) {
	if groups, err := fishCmdRegexp.Groups(text); err == nil {
		currentEvent.command = groups[0]
		currentEvent.number = currentNumber
	} else if groups, err := fishWhenRegexp.Groups(text); err == nil {
		converted, _ := strconv.ParseInt(groups[0], 10, 64)
		currentEvent.timestamp = converted
		currentEvent = keeper.Commit(currentEvent, historyChan)
		currentNumber++
	}

	return false, currentNumber, currentEvent
}
