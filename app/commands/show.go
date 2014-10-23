package commands

import (
	"fmt"
	// "regexp"

	"../environments"
	"../history_entries"
	"../slices"
	"../utils"
)

func Show(slice *slices.Slice, filter *utils.Regexp, env *environments.Environment) {
	var commands []history_entries.HistoryEntry

	if slice.Start >= 0 && slice.Finish >= 0 {
		keeper, err := history_entries.GetCommands(history_entries.GET_COMMANDS_RANGE, filter, env)
		if err != nil {
			return
		}
		commands = keeper.Result().([]history_entries.HistoryEntry)
	} else {
		keeper, err := history_entries.GetCommands(history_entries.GET_COMMANDS_ALL, filter, env)
		if err != nil {
			return
		}
		toBeRanged := keeper.Result().([]history_entries.HistoryEntry)
		sliceStart := slices.GetSliceIndex(slice.Start, len(toBeRanged))
		sliceFinish := slices.GetSliceIndex(slice.Finish, len(toBeRanged))
		if sliceStart < 0 || sliceFinish < 0 || sliceFinish <= sliceStart {
			return
		}
		commands = toBeRanged[sliceStart:sliceFinish]
	}

	for idx := 0; idx < len(commands); idx++ {
		fmt.Println(commands[idx].ToString(env))
	}
}
