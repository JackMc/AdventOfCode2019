package main

import (
	"fmt"
	"strings"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func moveRobot(position [2]int, direction int) [2]int {
	moveX := 0
	moveY := 0

	fmt.Println("trying to move robot in direction", direction)

	switch direction {
	case 0:
		fmt.Println("Moving left")
		moveX = -1
	case 1:
		fmt.Println("Moving up")
		moveY = 1
	case 2:
		fmt.Println("Moving right")
		moveX = 1
	case 3:
		fmt.Println("Moving down")
		moveY = -1
	default:
		panic("Unknown direction")
	}

	position[0] += moveX
	position[1] += moveY
	fmt.Println("new position:", position)
	return position
}

func main() {
	program := intcode.ParseProgram("input")
	interpreter, err := intcode.EvaluateProgramWithStringInput(program, "")
	if err != nil {
		fmt.Println(err)
		panic("Error while running interpreter")
	}

	outputBuilder := interpreter.Sink().(*strings.Builder)
	output := outputBuilder.String()
	split := strings.Split(output, "\n")
	blockTiles := 0
	for i, line := range split {
		idx := i % 3
		if idx == 0 {
			// x pos
		} else if idx == 1 {
			// y pos
		} else if idx == 2 {
			// type of tile
			if line == "2" {
				blockTiles++
			}
		}
	}

	fmt.Println(blockTiles)
}
