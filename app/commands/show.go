package commands

import (
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/slices"
	"github.com/9seconds/ah/app/utils"
)

// Show implements s (show) command.
func Show(slice *slices.Slice, filter *utils.Regexp, env *environments.Environment) {
	var commands []historyentries.HistoryEntry

	if slice.Start >= 0 && slice.Finish >= 0 {
		keeper, err := historyentries.GetCommands(historyentries.GetCommandsRange,
			filter, env, slice.Start, slice.Finish)
		if err != nil {
			return
		}
		commands = keeper.Result().([]historyentries.HistoryEntry)
	} else {
		keeper, err := historyentries.GetCommands(historyentries.GetCommandsAll, filter, env)
		if err != nil {
			return
		}
		toBeRanged := keeper.Result().([]historyentries.HistoryEntry)
		sliceStart := slices.GetSliceIndex(slice.Start, len(toBeRanged))
		sliceFinish := slices.GetSliceIndex(slice.Finish, len(toBeRanged))
		if sliceStart < 0 || sliceFinish < 0 || sliceFinish <= sliceStart {
			return
		}
		commands = toBeRanged[sliceStart:sliceFinish]
	}

	for idx := 0; idx < len(commands); idx++ {
		os.Stdout.WriteString(commands[idx].ToString(env))
		os.Stdout.WriteString("\n")
	}
}
