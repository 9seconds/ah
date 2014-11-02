package main

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"

	logrus "github.com/Sirupsen/logrus"
	docopt "github.com/docopt/docopt-go"

	"github.com/9seconds/ah/app/commands"
	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/slices"
	"github.com/9seconds/ah/app/utils"
)

const docoptOptions = `ah - A better history.

Ah is a better way to traverse the history of your shell prompts. You can
store an outputs of your commands, you can search with regular expressions
out of box, bookmark commands etc.

You may want to check detailed readme at https://github.com/9seconds/ah

Just a short reminder on possible subcommands:
    - s  - shows extended output from your HISTFILE
    - b  - bookmarks any command you want to have a faster access.
    - e  - executes a command by its bookmark name or history number.
    - t  - traces an output of the command and stores it safely.
    - l  - lists you an output of the command.
    - lb - lists available bookmarks.
    - rb - removes bookmarks.
    - gt - garbage collecting of the traces. Cleans old outputs.
    - gb - garbage collecting of the bookmarks. Swipes out old ones.

Usage:
    ah [options] s [-z] [-g PATTERN] [<lastNcommands> | <startFromNCommand> <finishByMCommand>]
    ah [options] b <commandNumber> <bookmarkAs>
    ah [options] e [-x] [-y] <commandNumberOrBookMarkName>
    ah [options] t [-x] [-y] [--] <command>...
    ah [options] l <numberOfCommandYouWantToCheck>
    ah [options] lb
    ah [options] rb <bookmarkToRemove>...
    ah [options] (gt | gb) (--keepLatest <keepLatest> | --olderThan <olderThan> | --all)
    ah (-h | --help)
    ah --version

Options:
    -s SHELL, --shell=SHELL
       Shell flavour you are using.
       By default, ah will do some shallow investigations.
    -f HISTFILE, --histfile=HISTFILE
       The path to a history file.
       By default ah will try to use default history file of your shell.
    -t HISTTIMEFORMAT, --histtimeformat=HISTTIMEFORMAT
       A time format for history output. Will use $HISTTIMEFORMAT by default.
    -d APPDIR, --appdir=APPDIR
       A place where ah has to store its data.
    -m TMPDIR, --tmpdir=TMPDIR
       A temporary place where ah stores an output. Set it only if you need it.
    -g PATTERN, --grep PATTERN
       A pattern to filter command lines. It is regular expression if no -f option is set.
    -y, --tty
       Allocates pseudo-tty is necessary.
    -x, --run-in-real-shell
       Runs a command in real interactive shell.
    -z, --fuzzy
       Interpret -g pattern as fuzzy match string.
    -v, --debug
       Shows a debug log of command execution.`

const version = "ah 0.9"

var validateBookmarkName = utils.CreateRegexp(`^[A-Za-z_]\w*$`)

type executor func(map[string]interface{}, *environments.Environment)

func main() {
	defer func() {
		if exc := recover(); exc != nil {
			utils.Logger.Fatal(exc)
		}
	}()

	arguments, err := docopt.Parse(docoptOptions, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	if arguments["--debug"].(bool) {
		utils.EnableLogging()
	} else {
		utils.DisableLogging()
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

	argTmpDir := arguments["--tmpdir"]
	if argTmpDir == nil {
		env.DiscoverTmpDir()
	} else {
		env.SetTmpDir(argTmpDir.(string))
	}

	utils.Logger.WithFields(logrus.Fields{
		"arguments":   arguments,
		"environment": env,
	}).Debug("Ready to start")

	utils.Logger.WithFields(logrus.Fields{
		"error": os.MkdirAll(env.GetTracesDir(), 0777),
	}).Info("Create traces dir")
	utils.Logger.WithFields(logrus.Fields{
		"error": os.MkdirAll(env.GetBookmarksDir(), 0777),
	}).Info("Create bookmarks dir")
	utils.Logger.WithFields(logrus.Fields{
		"error": os.MkdirAll(env.GetTmpDir(), 0777),
	}).Info("Create create temporary dir")

	var exec executor
	switch {
	case arguments["t"].(bool):
		utils.Logger.Info("Execute command 'tee'")
		exec = executeTee
	case arguments["s"].(bool):
		utils.Logger.Info("Execute command 'show'")
		exec = executeShow
	case arguments["l"].(bool):
		utils.Logger.Info("Execute command 'listTrace'")
		exec = executeListTrace
	case arguments["b"].(bool):
		utils.Logger.Info("Execute command 'listTrace'")
		exec = executeListTrace
	case arguments["e"].(bool):
		utils.Logger.Info("Execute command 'execute'")
		exec = executeExec
	case arguments["lb"].(bool):
		utils.Logger.Info("Execute command 'listBookmarks'")
		exec = executeListBookmarks
	case arguments["rb"].(bool):
		utils.Logger.Info("Execute command 'removeBookmarks'")
		exec = executeRemoveBookmarks
	case arguments["gt"].(bool) || arguments["gb"].(bool):
		utils.Logger.Info("Execute command 'gc'")
		exec = executeGC
	default:
		utils.Logger.Panic("Unknown command. Please be more precise")
		return
	}
	exec(arguments, env)
}

func executeTee(arguments map[string]interface{}, env *environments.Environment) {
	cmds := arguments["<command>"].([]string)
	cmd := strings.Join(cmds, " ")
	tty := arguments["--tty"].(bool)
	interactive := arguments["--run-in-real-shell"].(bool)

	utils.Logger.WithFields(logrus.Fields{
		"command":     cmd,
		"pseudo-tty":  tty,
		"interactive": interactive,
	}).Info("Arguments of 'tee'")

	commands.Tee(cmd, interactive, tty, env)
}

func executeShow(arguments map[string]interface{}, env *environments.Environment) {
	slice, err := slices.ExtractSlice(
		arguments["<lastNcommands>"],
		arguments["<startFromNCommand>"],
		arguments["<finishByMCommand>"])
	if err != nil {
		utils.Logger.Panic(err)
	}

	var filter *utils.Regexp
	if arguments["--grep"] != nil {
		query := arguments["--grep"].(string)
		if arguments["--fuzzy"].(bool) {
			regex := new(bytes.Buffer)
			for _, character := range query {
				regex.WriteString(".*?")
				regex.WriteString(regexp.QuoteMeta(string(character)))
			}
			regex.WriteString(".*?")
			query = regex.String()
		}
		filter = utils.CreateRegexp(query)
	}

	utils.Logger.WithFields(logrus.Fields{
		"slice":  slice,
		"filter": filter,
	}).Info("Arguments of 'show'")

	commands.Show(slice, filter, env)
}

func executeListTrace(arguments map[string]interface{}, env *environments.Environment) {
	cmd := arguments["<numberOfCommandYouWantToCheck>"].(string)

	utils.Logger.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info("Arguments of 'listTrace'")

	commands.ListTrace(cmd, env)
}

func executeBookmark(arguments map[string]interface{}, env *environments.Environment) {
	number, err := strconv.Atoi(arguments["<commandNumber>"].(string))
	if err != nil {
		utils.Logger.Panicf("Cannot understand command number: %d", number)
	}

	bookmarkAs := arguments["<bookmarkAs>"].(string)
	if !validateBookmarkName.Match(bookmarkAs) {
		utils.Logger.Panic("Incorrect bookmark name!")
	}

	utils.Logger.WithFields(logrus.Fields{
		"commandNumber": number,
		"bookmarkAs":    bookmarkAs,
	}).Info("Arguments of 'bookmark'")

	commands.Bookmark(number, bookmarkAs, env)
}

func executeExec(arguments map[string]interface{}, env *environments.Environment) {
	commandNumberOrBookMarkName := arguments["<commandNumberOrBookMarkName>"].(string)
	tty := arguments["--tty"].(bool)
	interactive := arguments["--run-in-real-shell"].(bool)

	utils.Logger.WithFields(logrus.Fields{
		"commandNumberOrBookMarkName": commandNumberOrBookMarkName,
		"tty":         tty,
		"interactive": interactive,
	}).Info("Arguments of 'bookmark'")

	commandNumber, err := strconv.Atoi(commandNumberOrBookMarkName)
	switch {
	case err == nil:
		utils.Logger.Info("Execute command number ", commandNumber)
		commands.ExecuteCommandNumber(commandNumber, interactive, tty, env)
	case validateBookmarkName.Match(commandNumberOrBookMarkName):
		utils.Logger.Info("Execute bookmark ", commandNumberOrBookMarkName)
		commands.ExecuteBookmark(commandNumberOrBookMarkName, interactive, tty, env)
	default:
		utils.Logger.Panic("Incorrect bookmark name! It should be started with alphabet letter, and alphabet or digits after!")
	}
}

func executeGC(arguments map[string]interface{}, env *environments.Environment) {
	gcDir := commands.GcTracesDir
	if arguments["gb"].(bool) {
		gcDir = commands.GcBookmarksDir
	}

	var gcType commands.GcType
	stringParam := "1"
	switch {
	case arguments["--keepLatest"].(bool):
		gcType = commands.GcKeepLatest
		stringParam = arguments["<keepLatest>"].(string)
	case arguments["--olderThan"].(bool):
		gcType = commands.GcOlderThan
		stringParam = arguments["<keepLatest>"].(string)
	case arguments["--all"].(bool):
		gcType = commands.GcAll
	default:
		utils.Logger.Panic("Unknown subcommand command")
	}

	param, err := strconv.Atoi(stringParam)
	if err != nil {
		utils.Logger.Panic(err)
	} else if param <= 0 {
		utils.Logger.Panic("Parameter of garbage collection has to be > 0")
	}

	utils.Logger.WithFields(logrus.Fields{
		"gcType": gcType,
		"param":  param,
	}).Info("Arguments")

	commands.GC(gcType, gcDir, param, env)
}

func executeListBookmarks(_ map[string]interface{}, env *environments.Environment) {
	commands.ListBookmarks(env)
}

func executeRemoveBookmarks(arguments map[string]interface{}, env *environments.Environment) {
	bookmarks, ok := arguments["<bookmarkToRemove>"].([]string)
	if !ok || bookmarks == nil || len(bookmarks) == 0 {
		utils.Logger.Info("Nothing to do here")
		return
	}

	for _, bookmark := range bookmarks {
		if !validateBookmarkName.Match(bookmark) {
			utils.Logger.WithFields(logrus.Fields{
				"bookmark": bookmark,
			}).Panicf("Bookmark name %s is invalid", bookmark)
		}
	}

	commands.RemoveBookmarks(bookmarks, env)
}
