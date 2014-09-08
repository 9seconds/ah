package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/alecthomas/kingpin.v1"

	"./app"
)

const (
	APP_DESCRIPTION          = "ah - A Better history."
	SHELL_FLAVOUR_DESRIPTION = "A shell flavour you are using." +
		" Vaild options are \"bash\" and \"zsh\"." +
		" By default ah tries to sniff what shell you are using and if you" +
		" have a nonstandard setup, bash would be used."
	HISTFILE_DESCRIPTION = "The path to your shell history file. " +
		"No worries if you are not sure or just do not want to set it, " +
		"ah will try to sniff. But it costs 1 interpreter start."
	HISTTIMEFORMAT_DESCRIPTION = "If you want to set different time format " +
		"of the history file, the best way to do it here"
	APP_DIR_DESCRIPTION = "An ah's directory for own storage"

	SHOW_CMD_DESCRIPTION       = "Shows an enhanced history of your commands."
	SHOW_CMD_GREP_DESCRIPTION  = "Filter output by given regular expression."
	SHOW_CMD_FUZZY_DESCRIPTION = "Interpret grep expression as fuzzy search " +
		"string instead of regular expression."
	SHOW_CMD_SLICE_DESCRIPTION = "Basically it could be a single argument " +
		"or a slice. Let's say it is 20. Then ah will show you 20 latest " +
		"records. So it is an equialent of :-20. What does negative number " +
		"means? It is rather simple: n commands from the end so it is a " +
		"shortcut of \"len(history)-20\". By default whole history would " +
		"be shown."
)

const (
	APP_DIR = "~/.ah"
)

var (
	HISTFILE_BASH = ".bash_history"
	HISTFILE_ZSH  = ".zsh_history"
)

var (
	application = kingpin.New("ah", APP_DESCRIPTION)
	histfile    = application.Flag("histfile", HISTFILE_DESCRIPTION).
			Short('f').
			String()
	histtimeformat = application.Flag("histtimeformat", HISTTIMEFORMAT_DESCRIPTION).
			Short('t').
			String()
	shell_flavour = application.Flag("shell", SHELL_FLAVOUR_DESRIPTION).
			Short('s').
			String()
	app_path = application.Flag("dir", APP_DIR_DESCRIPTION).
			Default(APP_DIR).
			String()

	show      = application.Command("s", SHOW_CMD_DESCRIPTION)
	show_grep = show.Flag("grep", SHOW_CMD_GREP_DESCRIPTION).
			Short('g').
			String()
	show_fuzzy = show.Flag("fuzzy", SHOW_CMD_FUZZY_DESCRIPTION).
			Short('f').Bool()
	show_arg = show.Arg("slice", SHOW_CMD_SLICE_DESCRIPTION).
			String()
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Impossible to detect current user\n"))
		os.Exit(1)
	}

	HISTFILE_BASH = filepath.Join(currentUser.HomeDir, HISTFILE_BASH)
	HISTFILE_ZSH = filepath.Join(currentUser.HomeDir, HISTFILE_ZSH)
}

func main() {
	defer func() {
		if exc := recover(); exc != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v\n", exc))
			os.Exit(1)
		}
	}()

	command := kingpin.MustParse(application.Parse(os.Args[1:]))
	env := app.Environment{}

	env.Shell = *shell_flavour
	if env.Shell == "" {
		env.Shell = os.Getenv("SHELL")
	}
	env.Shell = path.Base(env.Shell)
	if env.Shell != "zsh" && env.Shell != "bash" {
		panic("Sorry, ah supports only bash and zsh")
	}

	env.HistFile = *histfile
	if env.HistFile == "" {
		env.HistFile = os.Getenv("HISTFILE")
	}
	if env.HistFile == "" {
		if env.Shell == "bash" {
			env.HistFile = HISTFILE_BASH
		} else {
			env.HistFile = HISTFILE_ZSH
		}
	}

	env.HistTimeFormat = *histtimeformat
	if env.HistTimeFormat == "" {
		env.HistTimeFormat = os.Getenv("HISTTIMEFORMAT")
	}

	env.AppDir = *app_path

	switch command {
	case "s":
		slice, err := app.ExtractSlice(*show_arg)
		if err != nil {
			panic(err)
		}

		var filter *regexp.Regexp = nil
		if *show_grep != "" {
			filter = regexp.MustCompile(*show_grep)
		}

		app.CommandShow(slice, filter, &env)
	default:
		panic("Unknown command. Please specify at least one.")
	}
}
