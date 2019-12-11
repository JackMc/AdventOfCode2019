package main

import (
	"fmt"
	"strconv"
	"strings"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
	prmt "github.com/gitchander/permutation"
)

func main() {
	initialMemory := intcode.ParseProgram("input")
	phaseSettings := []string{"9", "7", "8", "5", "6"}

	permutations := prmt.New(prmt.StringSlice(phaseSettings))
	max := 0

	for permutations.Next() {
		amplifierInterpreters := [5]*intcode.InterpreterState{}

		// Start up the amplifiers, give them their phase settings, and
		// check that they all end up waiting for input
		for amplifierNumber, phaseSetting := range phaseSettings {
			interpreter, err := intcode.EvaluateProgramWithStringInput(initialMemory, phaseSetting+"\n")

			if err != intcode.NeedMoreInput {
				panic("Interpreter has not exited with expected error")
			}

			amplifierInterpreters[amplifierNumber] = &interpreter
		}
		// We kick off the process by giving the 0th interpreter a 0
		input := "0\n"
		amplifierNumber := 0
		for !amplifierInterpreters[len(amplifierInterpreters)-1].IsTerminated() {
			interpreter := amplifierInterpreters[amplifierNumber]
			reader := strings.NewReader(input)
			interpreter.AddNewInput(reader)
			interpreter.Run()
			stringBuilder := interpreter.Sink().(*strings.Builder)
			output := stringBuilder.String()
			input = output
			stringBuilder.Reset()

			input = output

			amplifierNumber = (amplifierNumber + 1) % len(amplifierInterpreters)
		}

		thrusterValue, _ := strconv.Atoi(strings.TrimSpace(input))

		if thrusterValue > max {
			fmt.Println("New max: a =", phaseSettings[0], "b =", phaseSettings[1], "c =", phaseSettings[2], "d=", phaseSettings[3], "e=", phaseSettings[4], "result =", thrusterValue)
			max = thrusterValue
		}
	}
}
