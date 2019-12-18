package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

const (
	empty            = 0
	wall             = 1
	block            = 2
	horizontalPaddle = 3
	ball             = 4
)

var (
	tileToChar = map[int]string{
		0: " ",
		1: "W",
		2: "B",
		3: "_",
		4: "X",
	}
)

func displayScreen(screenStr string) {
	split := strings.Split(strings.TrimSpace(screenStr), "\n")
	screen := [24][40]int{}
	score := -2
	xPos, yPos, tileType := -2, -2, -2

	for i, line := range split {
		idx := i % 3

		lineInt, err := strconv.Atoi(line)

		if err != nil {
			fmt.Println(err)
			panic("Error while Atoiing output")
		}

		if idx == 0 {
			xPos = lineInt
		} else if idx == 1 {
			yPos = lineInt
		} else if idx == 2 {
			tileType = lineInt
			// Save
			if xPos == -2 || yPos == -2 || tileType == -2 {
				panic("bad input")
			}

			if xPos == -1 && yPos == 0 {
				score = tileType
			} else {
				screen[yPos][xPos] = tileType
				// if xPos == 0 {
				// 	fmt.Println()
				// }

				// fmt.Print(tileToChar[tileType])
			}
		}
	}

	for _, row := range screen {
		for _, tileType := range row {
			fmt.Print(tileToChar[tileType])
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Score:", score)
}

func main() {
	program := intcode.ParseProgram("input")
	reader := bufio.NewReader(os.Stdin)
	// Unlimited quarters mode!
	program[0] = 2
	interpreter, err := intcode.EvaluateProgramWithStringInput(program, "")
	if err != intcode.NeedMoreInput {
		panic("Interpreter quit early")
	}

	for !interpreter.IsTerminated() {
		outputBuilder := interpreter.Sink().(*strings.Builder)
		output := outputBuilder.String()
		displayScreen(output)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			panic("Error while getting input")
		}
		interpreter.AddNewInput(strings.NewReader(input))
		memory, err := interpreter.Run()
		fmt.Println(memory)

		if err != intcode.NeedMoreInput {
			panic("Interpreter quit early")
		}
	}
}
