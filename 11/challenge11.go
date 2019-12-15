package main

import (
	"fmt"
	"strconv"
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
	currentPosition := [2]int{0, 0}
	// 0 = left, 1 = up, 2 = right, 3 = down, all done mod 4
	currentDirection := 1
	// Position -> colour (0 = black, 1 = white), default 0 (black)
	paintedPanels := map[[2]int]int{}
	program := intcode.ParseProgram("input")
	interpreter, err := intcode.EvaluateProgramWithStringInput(program, "")
	if err != intcode.NeedMoreInput {
		panic("Interpreter quit early")
	}

	// outputs := [][2]int{[2]int{1, 0}, [2]int{0, 0}, [2]int{1, 0}, [2]int{1, 0}, [2]int{0, 1}, [2]int{1, 0}, [2]int{1, 0}}
	// Sink stays the same, so we can just get it now and use it later
	builder := interpreter.Sink().(*strings.Builder)

	for !interpreter.IsTerminated() {
		currentPanelColour, _ := paintedPanels[currentPosition]
		fmt.Println("Current panel colour", currentPanelColour)
		strPanelColour := strconv.FormatInt(int64(currentPanelColour), 10)

		// Robot takes the current colour and outputs the new colour
		// + which direction to move next
		reader := strings.NewReader(strPanelColour + "\n")

		interpreter.AddNewInput(reader)
		_, err = interpreter.Run()
		output := builder.String()
		builder.Reset()
		split := strings.Split(output, "\n")
		newColour, direction := split[0], split[1]
		newColourInt, err := strconv.Atoi(newColour)

		if err != nil {
			fmt.Println(err)
			panic("Error while atoi colour")
		}

		directionInt, err := strconv.Atoi(direction)
		if err != nil {
			fmt.Println(err)
			panic("Error while atoiing direction")
		}

		paintedPanels[currentPosition] = newColourInt
		fmt.Println(directionInt)
		if directionInt == 0 {
			currentDirection = (currentDirection - 1 + 4) % 4
		} else {
			currentDirection = (currentDirection + 1) % 4
		}

		currentPosition = moveRobot(currentPosition, currentDirection)
	}

	fmt.Println(len(paintedPanels))
}
