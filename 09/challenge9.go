package main

import (
	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

func main() {
	interpreter := intcode.NewInteractiveInterpreter("input")
	interpreter.Run()
}
