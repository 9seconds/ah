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
	zshLineRegexp       = utils.CreateRegexp(`^: (\d+):\d;(.*)$`)
)

type (
	// Parser is a signature for a function which parses file and returns a Keeper.
	Parser func(Keeper, *bufio.Scanner, *utils.Regexp, chan *HistoryEntry) (Keeper, error)

	// ShellSpecificParser is a signature for a function which implements shell specific logic for parsing.
	ShellSpecificParser func(Keeper, string, uint, *HistoryEntry, *utils.Regexp, chan *HistoryEntry) (bool, uint, *HistoryEntry)
)

func getParser(env *environments.Environment) Parser {
	shell, _ := env.GetShell()
	logger, _ := env.GetLogger()

	shellSpecific := parseBash
	if shell == environments.ShellZsh {
		shellSpecific = parseZsh
	}

	return func(keeper Keeper, scanner *bufio.Scanner, filter *utils.Regexp, historyChan chan *HistoryEntry) (Keeper, error) {
		defer close(historyChan)

		var currentNumber uint
		continueToConsume := false
		currentEvent := keeper.Init()

		for scanner.Scan() && keeper.Continue() {
			text := scanner.Text()

			logger.WithFields(logrus.Fields{
				"text": text,
			}).Info("Parse history line")

			if continueToConsume {
				logger.Info("Attach the line to the previous command")

				currentEvent.command += "\n" + text
				if strings.HasSuffix(text, `\`) {
					continue
				}
				continueToConsume = false
				currentEvent = keeper.Commit(currentEvent, historyChan)
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
		if converted, err := strconv.Atoi(text[1:]); err == nil {
			currentEvent.timestamp = converted
		}
	} else {
		if filter == nil || filter.Match(text) {
			currentEvent.command = text
			currentEvent.number = currentNumber

			continueToConsume = strings.HasSuffix(text, "\\")
			if !continueToConsume {
				currentEvent = keeper.Commit(currentEvent, historyChan)
			}
		}
		currentNumber++
	}

	return continueToConsume, currentNumber, currentEvent
}

func parseZsh(keeper Keeper, text string, currentNumber uint, currentEvent *HistoryEntry, filter *utils.Regexp, historyChan chan *HistoryEntry) (bool, uint, *HistoryEntry) {
	continueToConsume := false
	groups, err := zshLineRegexp.Groups(text)

	if err != nil {
		return continueToConsume, currentNumber, currentEvent
	}
	timestamp, command := groups[0], groups[1]
	currentNumber++

	if filter != nil && !filter.Match(command) {
		return continueToConsume, currentNumber, currentEvent
	}

	converted, _ := strconv.Atoi(timestamp)
	currentEvent.command = command
	currentEvent.number = currentNumber
	currentEvent.timestamp = converted

	continueToConsume = strings.HasSuffix(text, "\\")
	if !continueToConsume {
		currentEvent = keeper.Commit(currentEvent, historyChan)
	}

	return continueToConsume, currentNumber, currentEvent
}
