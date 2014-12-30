package commands

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

const teeDelta = 1

// Tee implements t (trace, tee) command.
func Tee(input string, interactive bool, pseudoTTY bool, env *environments.Environment) {
	output, err := ioutil.TempFile(env.TmpDir, "ah")
	if err != nil {
		utils.Logger.Panic("Cannot create temporary file")
	}

	bufferedOutput := bufio.NewWriter(output)
	gzippedWrapper := utils.NewSynchronizedWriter(gzip.NewWriter(bufferedOutput))
	combinedStdout := io.MultiWriter(os.Stdout, gzippedWrapper)
	combinedStderr := io.MultiWriter(os.Stderr, gzippedWrapper)

	var commandError *exec.ExitError
	defer func() {
		// defer here because command may cause a panic but we do not want to lose any output
		gzippedWrapper.Close()
		bufferedOutput.Flush()
		output.Close()

		if hash, err := getPreciseHash(input, env); err == nil {
			err = os.Rename(output.Name(), env.GetTraceFileName(hash))
			if err != nil {
				utils.Logger.Errorf("Cannot save trace: %v. Get it here: %s", err, output.Name())
			} else {
				os.Remove(output.Name())
			}
		} else {
			utils.Logger.Errorf("Error occured on fetching command number: %v", err)
		}

		if commandError != nil {
			os.Exit(utils.GetStatusCode(commandError))
		}
	}()

	commandError = utils.Exec(input,
		string(env.Shell), interactive, pseudoTTY,
		os.Stdin, combinedStdout, combinedStderr)
}

func getPreciseHash(cmd string, env *environments.Environment) (hash string, err error) {
	commands, err := historyentries.GetCommands(historyentries.GetCommandsAll, nil, env)
	if err != nil {
		err = fmt.Errorf("Cannot fetch commands list: %v", err)
		return
	}
	commandList := commands.Result().([]historyentries.HistoryEntry)

	if len(commandList) == 0 {
		err = errors.New("Command list is empty")
		return
	}

	// Why do I need such complicated logic here? The reason is trivial:
	// one may use substitutions in a command like
	// ah t -- `which python` script $(docker images -a -q)
	// basically it manages these situations as precise as possible.

	// Candidates here just means that several command may be executed
	// simultaneously (e.g. with tmuxinator) so timestamp is not precise
	// identifier here.

	candidates, err := getCandidates(commandList)
	if err != nil {
		err = fmt.Errorf("Cannot detect proper command: %v", err)
		return
	}

	if len(candidates) == 1 {
		hash = candidates[0].GetTraceName()
		return
	}

	suitable := getSuitable(candidates, cmd)
	hash = (*suitable).GetTraceName()

	return
}

func getCandidates(commands []historyentries.HistoryEntry) (candidates []historyentries.HistoryEntry, err error) {
	start := -1
	finish := -1

	for idx := len(commands) - 1; idx >= 0; idx-- {
		timestamp := commands[idx].GetTimestamp()
		switch {
		case timestamp < environments.CreatedAt-teeDelta:
			break
		case timestamp > environments.CreatedAt+teeDelta:
			continue
		default:
			if finish == -1 {
				finish = idx + 1
			}
			start = idx
		}
	}

	if start == -1 || finish == -1 {
		utils.Logger.Warn("Cannot find anything on the time range " +
			"so assume that deduplication works. Take the latest.")
		candidates = commands[len(commands)-1:]
	} else {
		candidates = commands[start:finish]
	}

	return
}

func getSuitable(candidates []historyentries.HistoryEntry, cmd string) *historyentries.HistoryEntry {
	chunks := strings.Split(cmd, " ")
	maxHaveElements := -1
	maxIdx := 0

	for idx := 0; idx < len(candidates); idx++ {
		haveElements := countElements(candidates[idx].GetCommand(), chunks)
		if haveElements == len(chunks) {
			maxIdx = idx
			break
		}
		if haveElements > maxHaveElements {
			maxHaveElements = haveElements
			maxIdx = idx
		}
	}

	return &candidates[maxIdx]
}

func countElements(cmd string, chunks []string) (count int) {
	for idx := 0; idx < len(chunks); idx++ {
		if strings.Contains(cmd, chunks[idx]) {
			count++
		}
	}

	return
}
