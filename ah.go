package main

import (
	"fmt"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	fmt.Println("Hello world")
}
