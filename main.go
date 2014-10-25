package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	logrus "github.com/Sirupsen/logrus"
	docopt "github.com/docopt/docopt-go"

	"github.com/9seconds/ah/app/commands"
	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/slices"
	"github.com/9seconds/ah/app/utils"
)

const docoptOptions = `ah - A better history.

Ah is a better way to traverse the history of your shell prompts. Right now it
supports only 3 additional possibilities you are probably have dreamt about:
    1. Good searching;
    2. Better memoizing of entries;
    3. Storing an output of the commands.

Searching is done using 's' command. You may filter output or fetch a history
slice you are wondering about. Filtering uses regular expressions or fuzzy
matching out of box.

Memoizing means that you have an ability to bookmark some favourite commands
and use ah to store some sort of short snippets or ad-hoc shell scripts.

And ah gives you a possibility to have a persistent storage of an output
of any command you are executing. And you can return back to it any time
you want.

Usage:
    ah [options] s [-z] [-g PATTERN] [<lastNcommands> | <startFromNCommand> <finishByMCommand>]
    ah [options] b <commandNumber> <bookmarkAs>
    ah [options] e <commandNumberOrBookMarkName>
    ah [options] t [-y] [--] <command>...
    ah [options] l <numberOfCommandYouWantToCheck>
    ah [options] lb
    ah [options] rb <bookmarkToRemove>...
    ah [options] (gt | gb) (--keepLatest <keepLatest> | --olderThan <olderThan> | --all)
    ah (-h | --help)
    ah --version

Options:
    -s SHELL, --shell=SHELL                               Shell flavour you are using. By default, ah will do some shallow investigations.
    -f HISTFILE, --histfile=HISTFILE                      The path to a history file. By default ah will try to use default history file of your shell
    -t HISTTIMEFORMAT, --histtimeformat=HISTTIMEFORMAT    A time format for history output. Will use $HISTTIMEFORMAT by default.
    -d APPDIR, --appdir=APPDIR                            A place where ah has to store its data.
    -g PATTERN, --grep PATTERN                            A pattern to filter command lines. It is regular expression if no -f option is set.
    -y, --tty                                             Allocates pseudo-tty is necessary
    -z, --fuzzy                                           Interpret -g pattern as fuzzy match string.
    -v, --debug                                           Shows a debug log of command execution.`

const version = "ah 0.7"

var validateBookmarkName = utils.CreateRegexp(`^[A-Za-z_]\w*$`)

type executor func(map[string]interface{}, *environments.Environment)

func main() {
	defer func() {
		if exc := recover(); exc != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v\n", exc))
			os.Exit(1)
		}
	}()

	arguments, err := docopt.Parse(docoptOptions, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	env := new(environments.Environment)

	argShell := arguments["--shell"]
	if argShell == nil {
		env.DiscoverShell()
	} else {
		env.SetShell(argShell.(string))
	}

	argHistFile := arguments["--histfile"]
	if argHistFile == nil {
		env.DiscoverHistFile()
	} else {
		env.SetHistFile(argHistFile.(string))
	}

	argHistTimeFormat := arguments["--histtimeformat"]
	if argHistTimeFormat == nil {
		env.DiscoverHistTimeFormat()
	} else {
		env.SetHistTimeFormat(argHistTimeFormat.(string))
	}

	argAppDir := arguments["--appdir"]
	if argAppDir == nil {
		env.DiscoverAppDir()
	} else {
		env.SetAppDir(argAppDir.(string))
	}

	if arguments["--debug"].(bool) {
		env.EnableDebugLog()
	} else {
		env.DisableDebugLog()
	}

	logger, _ := env.GetLogger()
	logger.WithFields(arguments).Debug("Arguments")
	logger.Debug("Environment ", env)

	logger.WithFields(logrus.Fields{
		"error": os.MkdirAll(env.GetTracesDir(), 0777),
	}).Info("Create traces dir")
	logger.WithFields(logrus.Fields{
		"error": os.MkdirAll(env.GetBookmarksDir(), 0777),
	}).Info("Create bookmarks dir")

	var exec executor
	if arguments["t"].(bool) {
		logger.Info("Execute command 'tee'")
		exec = executeTee
	} else if arguments["s"].(bool) {
		logger.Info("Execute command 'show'")
		exec = executeShow
	} else if arguments["l"].(bool) {
		logger.Info("Execute command 'listTrace'")
		exec = executeListTrace
	} else if arguments["b"].(bool) {
		logger.Info("Execute command 'bookmark'")
		exec = executeBookmark
	} else if arguments["e"].(bool) {
		logger.Info("Execute command 'execute'")
		exec = executeExec
	} else if arguments["gt"].(bool) || arguments["gb"].(bool) {
		logger.Info("Execute command 'gc'")
		exec = executeGC
	} else if arguments["lb"].(bool) {
		logger.Info("Execute command 'listBookmarks'")
		exec = executeListBookmarks
	} else if arguments["rb"].(bool) {
		logger.Info("Execute command 'removeBookmarks'")
		exec = executeRemoveBookmarks
	} else {
		logger.Info("No valid choices", arguments)
		panic("Unknown command. Please be more precise")
	}
	exec(arguments, env)
}

func executeTee(arguments map[string]interface{}, env *environments.Environment) {
	cmds := arguments["<command>"].([]string)
	tty := arguments["--tty"].(bool)

	logger, _ := env.GetLogger()
	logger.WithFields(logrus.Fields{
		"command arguments": cmds,
		"pseudo-tty":        tty,
	}).Info("Arguments of 'tee'")

	commands.Tee(cmds, tty, env)
}

func executeShow(arguments map[string]interface{}, env *environments.Environment) {
	logger, _ := env.GetLogger()

	slice, err := slices.ExtractSlice(
		arguments["<lastNcommands>"],
		arguments["<startFromNCommand>"],
		arguments["<finishByMCommand>"])
	if err != nil {
		panic(err)
	}

	var filter *utils.Regexp
	if arguments["--grep"] != nil {
		query := arguments["--grep"].(string)
		if arguments["--fuzzy"].(bool) {
			regex := ""
			for _, character := range query {
				regex += ".*?" + regexp.QuoteMeta(string(character))
			}
			query = regex + ".*?"
		}
		filter = utils.CreateRegexp(query)
	}

	logger.WithFields(logrus.Fields{
		"slice":  slice,
		"filter": filter,
	}).Info("Arguments of 'show'")

	commands.Show(slice, filter, env)
}

func executeListTrace(arguments map[string]interface{}, env *environments.Environment) {
	cmd := arguments["<numberOfCommandYouWantToCheck>"].(string)

	logger, _ := env.GetLogger()
	logger.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info("Arguments of 'listTrace'")

	commands.ListTrace(cmd, env)
}

func executeBookmark(arguments map[string]interface{}, env *environments.Environment) {
	var commandNumber int
	if number, err := strconv.Atoi(arguments["<commandNumber>"].(string)); err != nil {
		panic(fmt.Sprintf("Cannot understand command number: %d", commandNumber))
	} else {
		commandNumber = number
	}

	bookmarkAs := arguments["<bookmarkAs>"].(string)
	if !validateBookmarkName.Match(bookmarkAs) {
		panic("Incorrect bookmark name!")
	}

	logger, _ := env.GetLogger()
	logger.WithFields(logrus.Fields{
		"commandNumber": commandNumber,
		"bookmarkAs":    bookmarkAs,
	}).Info("Arguments of 'bookmark'")

	commands.Bookmark(commandNumber, bookmarkAs, env)
}

func executeExec(arguments map[string]interface{}, env *environments.Environment) {
	commandNumberOrBookMarkName := arguments["<commandNumberOrBookMarkName>"].(string)

	logger, _ := env.GetLogger()
	logger.WithFields(logrus.Fields{
		"commandNumberOrBookMarkName": commandNumberOrBookMarkName,
	}).Info("Arguments of 'bookmark'")

	if commandNumber, err := strconv.Atoi(commandNumberOrBookMarkName); err == nil {
		logger.Info("Execute command number ", commandNumber)
		commands.ExecuteCommandNumber(commandNumber, env)
	} else if validateBookmarkName.Match(commandNumberOrBookMarkName) {
		logger.Info("Execute bookmark ", commandNumberOrBookMarkName)
		commands.ExecuteBookmark(commandNumberOrBookMarkName, env)
	} else {
		panic("Incorrect bookmark name! It should be started with alphabet letter, and alphabet or digits after!")
	}
}

func executeGC(arguments map[string]interface{}, env *environments.Environment) {
	logger, _ := env.GetLogger()

	var param int
	var gcType commands.GcType

	gcDir := commands.GcTracesDir
	if arguments["gb"].(bool) {
		gcDir = commands.GcBookmarksDir
	}

	if arguments["--keepLatest"].(bool) {
		gcType = commands.GcKeepLatest
		paramString := arguments["<keepLatest>"].(string)
		paramConverted, err := strconv.Atoi(paramString)
		if err != nil {
			panic(err)
		}
		param = paramConverted
	} else if arguments["--olderThan"].(bool) {
		gcType = commands.GcOlderThan
		paramString := arguments["<olderThan>"].(string)
		paramConverted, err := strconv.Atoi(paramString)
		if err != nil {
			panic(err)
		}
		param = paramConverted
	} else if arguments["--all"].(bool) {
		gcType = commands.GcAll
		param = 1
	} else {
		panic("Unknown subcommand command")
	}

	if param <= 0 {
		panic("Parameter of garbage collection has to be > 0")
	}

	logger.WithFields(logrus.Fields{
		"gcType": gcType,
		"param":  param,
	}).Info("Arguments")

	commands.GC(gcType, gcDir, param, env)
}

func executeListBookmarks(_ map[string]interface{}, env *environments.Environment) {
	commands.ListBookmarks(env)
}

func executeRemoveBookmarks(arguments map[string]interface{}, env *environments.Environment) {
	logger, _ := env.GetLogger()

	fmt.Println(arguments["<bookmarkToRemove>"])

	bookmarks, ok := arguments["<bookmarkToRemove>"].([]string)
	if !ok || bookmarks == nil || len(bookmarks) == 0 {
		logger.Info("Nothing to do here")
		return
	}

	for _, bookmark := range bookmarks {
		if !validateBookmarkName.Match(bookmark) {
			logger.WithFields(logrus.Fields{
				"bookmark": bookmark,
			}).Warn("Bookmark name is invalid")
			panic(fmt.Sprintf("Bookmark %s is invalid", bookmark))
		}
	}

	commands.RemoveBookmarks(bookmarks, env)
}
