package environments

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	logrus "github.com/Sirupsen/logrus"
	strftime "github.com/weidewang/go-strftime"

	"github.com/9seconds/ah/app/utils"
)

const (
	defaultAppDir = ".ah"
	tracesDir     = "traces"
	bookmarksDir  = "bookmarks"

	defaultZshHistFile  = ".zsh_history"
	defaultBashHistFile = ".bash_history"
)

// Codes for the supported shells
const (
	ShellZsh  = "zsh"
	ShellBash = "bash"
)

var currentUser *user.User

func init() {
	fetchedCurrentUser, err := user.Current()
	if err != nil {
		os.Stderr.WriteString("Impossible to detect current user\n")
		os.Exit(1)
	}
	currentUser = fetchedCurrentUser
}

// Environment has all required information about environment ah is executed in
type Environment struct {
	appDir         string
	histFile       string
	histTimeFormat string
	shell          string
	log            *logrus.Logger
}

// OK checks if all required information was collected.
func (e *Environment) OK() bool {
	return e.log != nil && e.appDir != "" && e.histFile != "" && e.shell != ""
}

// GetTracesDir returns an absolute path for the directory where traces should be stored.
func (e *Environment) GetTracesDir() string {
	return filepath.Join(e.appDir, tracesDir)
}

// GetBookmarksDir returns an absolute path for the directory where bookmarks should be stored.
func (e *Environment) GetBookmarksDir() string {
	return filepath.Join(e.appDir, bookmarksDir)
}

// GetTraceFileName returns an absolute file path for the trace with a given name.
func (e *Environment) GetTraceFileName(hash string) string {
	return filepath.Join(e.GetTracesDir(), hash)
}

// GetBookmarkFileName returns an absolute file path for the bookmark with a given name.
func (e *Environment) GetBookmarkFileName(name string) string {
	return filepath.Join(e.GetBookmarksDir(), name)
}

// DiscoverAppDir tries to discover app storage directory from the environment itself.
func (e *Environment) DiscoverAppDir() error {
	return e.SetAppDir(filepath.Join(currentUser.HomeDir, defaultAppDir))
}

// GetAppDir returns an absolute path of the app storage directory and error if occured.
func (e *Environment) GetAppDir() (string, error) {
	if e.appDir == "" {
		return "", errors.New("AppDir is not set yet")
	}
	return e.appDir, nil
}

// SetAppDir sets an app storage directory. Has to be an absolute path.
func (e *Environment) SetAppDir(path string) error {
	e.appDir = path
	return nil
}

// DiscoverShell discovers shell from the actual environment.
func (e *Environment) DiscoverShell() error {
	return e.SetShell(os.Getenv("SHELL"))
}

// GetShell returns a shell code from the actual envrionment.
func (e *Environment) GetShell() (string, error) {
	if e.shell == "" {
		return "", errors.New("Shell is not set yet")
	}
	return e.shell, nil
}

// SetShell explicitly sets shell.
func (e *Environment) SetShell(shell string) error {
	baseShell := path.Base(shell)

	if baseShell == ShellZsh || baseShell == ShellBash {
		e.shell = baseShell
		return nil
	}
	return fmt.Errorf("Shell %s is not supported", shell)
}

// DiscoverHistFile tries to discover history file from the actual environment.
func (e *Environment) DiscoverHistFile() error {
	shell, err := e.GetShell()
	if err != nil {
		return err
	}

	if shell == ShellZsh {
		return e.SetHistFile(filepath.Join(currentUser.HomeDir, defaultZshHistFile))
	}
	return e.SetHistFile(filepath.Join(currentUser.HomeDir, defaultBashHistFile))
}

// GetHistFile returns an absolute path of the history file from the actual environment.
func (e *Environment) GetHistFile() (string, error) {
	if e.histFile == "" {
		return "", errors.New("HistFile is not set yet")
	}
	return e.histFile, nil
}

// SetHistFile sets a path to the history file.
func (e *Environment) SetHistFile(path string) error {
	e.histFile = path
	return nil
}

// DiscoverHistTimeFormat discovers time format of the history entries from the environment.
func (e *Environment) DiscoverHistTimeFormat() error {
	e.histTimeFormat = os.Getenv("HISTTIMEFORMAT")
	return nil
}

// GetHistTimeFormat returns a history time format.
func (e *Environment) GetHistTimeFormat() (string, error) {
	if e.histTimeFormat == "" {
		return "", errors.New("HistTimeFormat is not set yet")
	}
	return e.histTimeFormat, nil
}

// SetHistTimeFormat sets a history time format.
func (e *Environment) SetHistTimeFormat(histTimeFormat string) error {
	e.histTimeFormat = histTimeFormat
	return nil
}

// FormatTimeStamp formats a timestamp according to the environment settings.
func (e *Environment) FormatTimeStamp(timestamp int) (string, error) {
	return e.FormatTime(utils.ConvertTimestamp(timestamp))
}

// FormatTime formats a time structure according to the environment settings.
func (e *Environment) FormatTime(timestamp *time.Time) (string, error) {
	format, err := e.GetHistTimeFormat()
	if err != nil {
		return "", err
	}
	return strftime.Strftime(timestamp, format), nil
}

// GetLogger returns a preconfigured logger.
func (e *Environment) GetLogger() (*logrus.Logger, error) {
	if e.log == nil {
		return nil, errors.New("Logger is not set yet")
	}
	return e.log, nil
}

// EnableDebugLog makes logging verbose.
func (e *Environment) EnableDebugLog() {
	e.log = logrus.New()
	e.log.Out = os.Stderr
	e.log.Level = logrus.DebugLevel
}

// DisableDebugLog makes logging silent.
func (e *Environment) DisableDebugLog() {
	e.log = logrus.New()
	e.log.Out = ioutil.Discard
	e.log.Level = logrus.ErrorLevel
}

// GetTraceFilenames returns a list of filenames for traces.
func (e *Environment) GetTraceFilenames() ([]os.FileInfo, error) {
	fileInfos := make([]os.FileInfo, 0, 16)

	files, err := ioutil.ReadDir(e.GetTracesDir())
	if err != nil {
		return fileInfos, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfos = append(fileInfos, file)
	}

	return fileInfos, nil
}
