package main

import (
	"os"
	"strconv"
	"runtime/debug"
	"github.com/9seconds/ah/utils"
	// "runtime/pprof"
	"github.com/weidewang/go-strftime"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	// prof, _ := os.Create("ah.prof")
	// pprof.StartCPUProfile(prof)
	// defer pprof.StopCPUProfile()
    file, _ := os.Open(utils.HistoryFilePath)
	defer file.Close()

	scanner := utils.ShellEnv.GetHistoryScanner(file)
	commands, _ := scanner.GetCommands()

	content := make([][]string, len(commands))
	for idx, value := range commands {
		content[idx] = []string{
			strconv.FormatInt(int64(idx), 10),
			strftime.Strftime(&value.Timestamp, "%H:%M:%S"),
			value.Command,
		}
	}
	shrinkStatus := []bool{false, false, true}

	utils.PrintFormatted(content, shrinkStatus, 4)
}

