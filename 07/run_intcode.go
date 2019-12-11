package main

import (
	"os"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

func main() {
	interpreter := intcode.NewInteractiveInterpreter(os.Args[1])
	interpreter.Run()
}
