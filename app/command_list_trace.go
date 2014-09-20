package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func CommandListTrace(argument string, env *Environment) {
	number, err := strconv.Atoi(argument)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert argument to a number: %s", argument))
	}

	filename := env.GetTraceFileName(number)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("Output for %s is not exist", argument))
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("File %s is not readable", filename))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("\n\n-----------------n")
		fmt.Println("An error occured", err)
	}
}
