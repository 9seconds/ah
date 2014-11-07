package environments

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	strftime "github.com/jehiah/go-strftime"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/9seconds/ah/app/utils"
)

const (
	defaultAppDir = ".ah"
	tracesDir     = "traces"
	bookmarksDir  = "bookmarks"

	defaultZshHistFile          = ".zsh_history"
	defaultBashHistFile         = ".bash_history"
	defaultAutoCommandsFileName = "autocommands.gob"
)

// ShellType sets a type of the shell
type ShellType string

// Codes for the supported shells
const (
	ShellZsh  ShellType = "zsh"
	ShellBash ShellType = "bash"
	ShellFish ShellType = "fish"
)

var (
	// HomeDir defines a home of current user.
	HomeDir string

	// CreatedAt defines a time when program was executed.
	CreatedAt = time.Now().Unix()

	defaultFishHistFile = filepath.Join(".config", "fish", "fish_history")
)

func init() {
	currentHomeDir, err := homedir.Dir()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Cannot fetch your home directory: %v", err))
		os.Exit(1)
	}
	HomeDir = currentHomeDir
}

// Environment has all required information about environment ah is executed in
type Environment struct {
	appDir         string
	histFile       string
	histTimeFormat string
	shell          ShellType
	tmpDir         string
}

// OK checks if all required information was collected.
func (e *Environment) OK() bool {
	return e.appDir != "" && e.histFile != "" && e.shell != ""
}

// DiscoverTmpDir sets default directory for storing ah temporary data.
func (e *Environment) DiscoverTmpDir() {
	e.SetTmpDir(os.TempDir())
}

// SetTmpDir sets directory where to store temporary files for ah.
func (e *Environment) SetTmpDir(dir string) {
	e.tmpDir = dir
}

// GetTmpDir returns a path to temporary directory
// For the reason why it is required, please see this:
// https://github.com/9seconds/ah/issues/3
func (e *Environment) GetTmpDir() string {
	return e.tmpDir
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
func (e *Environment) DiscoverAppDir() {
	e.SetAppDir(filepath.Join(HomeDir, defaultAppDir))
}

// GetAppDir returns an absolute path of the app storage directory and error if occured.
func (e *Environment) GetAppDir() string {
	return e.appDir
}

// SetAppDir sets an app storage directory. Has to be an absolute path.
func (e *Environment) SetAppDir(path string) {
	e.appDir = path
}

// DiscoverShell discovers shell from the actual environment.
func (e *Environment) DiscoverShell() error {
	return e.SetShell(os.Getenv("SHELL"))
}

// GetShell returns a shell code from the actual envrionment.
func (e *Environment) GetShell() ShellType {
	return e.shell
}

// SetShell explicitly sets shell.
func (e *Environment) SetShell(shell string) (err error) {
	baseShell := ShellType(path.Base(shell))

	switch baseShell {
	case ShellZsh:
		fallthrough
	case ShellBash:
		fallthrough
	case ShellFish:
		e.shell = baseShell
	default:
		err = fmt.Errorf("Shell %s is not supported", shell)
	}

	return
}

// DiscoverHistFile tries to discover history file from the actual environment.
func (e *Environment) DiscoverHistFile() {
	switch e.shell {
	case ShellBash:
		e.histFile = filepath.Join(HomeDir, defaultBashHistFile)
	case ShellZsh:
		e.histFile = filepath.Join(HomeDir, defaultZshHistFile)
	case ShellFish:
		e.histFile = filepath.Join(HomeDir, defaultFishHistFile)
	}
}

// GetHistFile returns an absolute path of the history file from the actual environment.
func (e *Environment) GetHistFile() string {
	return e.histFile
}

// SetHistFile sets a path to the history file.
func (e *Environment) SetHistFile(path string) {
	e.histFile = path
}

// DiscoverHistTimeFormat discovers time format of the history entries from the environment.
func (e *Environment) DiscoverHistTimeFormat() {
	e.histTimeFormat = os.Getenv("HISTTIMEFORMAT")
}

// GetHistTimeFormat returns a history time format.
func (e *Environment) GetHistTimeFormat() string {
	return e.histTimeFormat
}

// SetHistTimeFormat sets a history time format.
func (e *Environment) SetHistTimeFormat(histTimeFormat string) {
	e.histTimeFormat = histTimeFormat
}

// FormatTimeStamp formats a timestamp according to the environment settings.
func (e *Environment) FormatTimeStamp(timestamp int64) (string, error) {
	return e.FormatTime(utils.ConvertTimestamp(timestamp))
}

// FormatTime formats a time structure according to the environment settings.
func (e *Environment) FormatTime(timestamp *time.Time) (string, error) {
	if e.histTimeFormat == "" {
		return "", errors.New("Cannot format time for absent time format")
	}
	return strftime.Format(e.histTimeFormat, *timestamp), nil
}

// GetTraceFilenames returns a list of filenames for traces.
func (e *Environment) GetTraceFilenames() ([]os.FileInfo, error) {
	return e.getFilenames(e.GetTracesDir())
}

// GetBookmarkFilenames returns a list of bookmarks for traces.
func (e *Environment) GetBookmarkFilenames() ([]os.FileInfo, error) {
	return e.getFilenames(e.GetBookmarksDir())
}

func (e *Environment) GetAutoCommandFileName() string {
	return filepath.Join(e.appDir, defaultAutoCommandsFileName)
}

func (e *Environment) getFilenames(directory string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfos = append(fileInfos, file)
	}

	return fileInfos, nil
}
