package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"github.com/docopt/docopt-go"

	"./app"
)

const DOCOPT = `ah - A better history.

Ah is a better way to traverse the history of your shell prompts. Right now it supports only 3 additional possibilities you are probably have dreamt about:
    1. Good searching;
    2. Better memoizing of entries;
    3. Storing an output of the commands.

Searching is done using 's' command. You may filter output or fetch a history slice you are wondering about. Filtering uses regular expressions or fuzzy matching out of box.

Memoizing means that you have an ability to bookmark some favourite commands and use ah to store some sort of short snippets or ad-hoc shell scripts.

And ah gives you a possibility to have a persistent storage of an output of any command you are executing. And you can return back to it any time you want.

Usage:
    ah [options] s [-z] [-g PATTERN] [<lastNcommands> | <startFromNCommand> <finishByMCommand>]
    ah [options] b <commandNumber> <bookmarkAs>
    ah [options] e <commandNumberOrBookMarkName>
    ah [options] t [--] <command>...
    ah [options] l <numberOfCommandYouWantToCheck>
    ah (-h | --help)
    ah --version

Options:
    -s SHELL, --shell=SHELL                               Shell flavour you are using. By default, ah will do some shallow investigations.
    -f HISTFILE, --histfile=HISTFILE                      The path to a history file. By default ah will try to use default history file of your shell
    -t HISTTIMEFORMAT, --histtimeformat=HISTTIMEFORMAT    A time format for history output. Will use $HISTTIMEFORMAT by default.
    -d APPDIR, --appdir=APPDIR                            A place where ah has to store its data.
    -g PATTERN, --grep PATTERN                            A pattern to filter command lines. It is regular expression if no -f option is set.
    -z, --fuzzy                                           Interpret -g pattern as fuzzy match string.`

const (
	DEFAULT_SHELL = "bash"
)

var (
	DEFAULT_BASH_HISTFILE = ".bash_history"
	DEFAULT_ZSH_HISTFILE  = ".zsh_history"
	DEFAULT_APPDIR        = ".ah"

	CURRENT_USER *user.User
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Impossible to detect current user\n"))
		os.Exit(1)
	}
	CURRENT_USER = currentUser

}

func main() {
	defer func() {
		if exc := recover(); exc != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v\n", exc))
			os.Exit(1)
		}
	}()

	DEFAULT_BASH_HISTFILE = filepath.Join(CURRENT_USER.HomeDir, DEFAULT_BASH_HISTFILE)
	DEFAULT_ZSH_HISTFILE = filepath.Join(CURRENT_USER.HomeDir, DEFAULT_ZSH_HISTFILE)
	DEFAULT_APPDIR = filepath.Join(CURRENT_USER.HomeDir, DEFAULT_APPDIR)

	arguments, err := docopt.Parse(DOCOPT, nil, true, "ah 0.1", false)
	if err != nil {
		panic(err)
	}
	env := app.Environment{}

	var argShell interface{} = arguments["--shell"]
	if argShell == nil {
		env.Shell = os.Getenv("SHELL")
	} else {
		env.Shell = argShell.(string)
	}
	env.Shell = path.Base(env.Shell)
	if env.Shell != "zsh" && env.Shell != "bash" {
		panic("Sorry, ah supports only bash and zsh")
	}

	var argHistFile interface{} = arguments["--histfile"]
	if argHistFile == nil {
		env.HistFile = os.Getenv("HISTFILE")
	} else {
		env.HistFile = argHistFile.(string)
	}
	if env.HistFile == "" {
		if env.Shell == "bash" {
			env.HistFile = DEFAULT_BASH_HISTFILE
		} else {
			env.HistFile = DEFAULT_ZSH_HISTFILE
		}
	}

	var argHistTimeFormat interface{} = arguments["--histtimeformat"]
	if argHistTimeFormat == nil {
		env.HistTimeFormat = os.Getenv("HISTTIMEFORMAT")
	} else {
		env.HistTimeFormat = argHistTimeFormat.(string)
	}

	var argAppDir interface{} = arguments["--appdir"]
	if argAppDir == nil {
		env.AppDir = DEFAULT_APPDIR
	} else {
		env.AppDir = argAppDir.(string)
	}

	os.MkdirAll(env.GetTracesDir(), 0777)
	os.MkdirAll(env.GetBookmarksDir(), 0777)

	if arguments["s"].(bool) {
		slice, err := app.ExtractSlice(
			arguments["<lastNcommands>"],
			arguments["<startFromNCommand>"],
			arguments["<finishByMCommand>"])
		if err != nil {
			panic(err)
		}

		var filter *regexp.Regexp
		if arguments["--grep"] != nil {
			query := arguments["--grep"].(string)
			if arguments["--fuzzy"].(bool) {
				regex := ""
				for _, character := range query {
					regex += ".*?" + regexp.QuoteMeta(string(character))
				}
				query = regex + ".*?"
			}
			filter = regexp.MustCompile(query)
		}

		app.CommandShow(slice, filter, &env)
	} else if arguments["t"].(bool) {
		commands := arguments["<command>"].([]string)

		app.CommandTee(commands, &env)
	} else {
		panic("Unknown command. Please be more precise.")
	}
}
