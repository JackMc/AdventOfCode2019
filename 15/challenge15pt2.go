package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

type vec2d struct {
	x, y int
}

func moveRobot(position vec2d, direction int64) vec2d {
	switch direction {
	case 1:
		position.y--
	case 2:
		position.y++
	case 3:
		position.x--
	case 4:
		position.x++
	}

	return position
}

func main() {
	program := intcode.ParseProgram("input")
	position := vec2d{0, 0}
	room := map[vec2d]string{}
	minX, maxX, minY, maxY := math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32

	interpreter, err := intcode.EvaluateProgramWithStringInput(program, "")
	builder := interpreter.Sink().(*strings.Builder)
	if err != intcode.NeedMoreInput {
		panic("Interpreter quit early")
	}

	numDots := 0
	initialOxygenPosition := vec2d{-1, -1}

	for i := int64(0); i < 1000000; i++ {
		direction := rand.Int63n(4) + 1
		newPositionCandidate := moveRobot(position, direction)
		if room[newPositionCandidate] == "#" {
			// Skip anything we know will be a wall
			continue
		}

		input := strconv.FormatInt(direction, 10) + "\n"

		// fmt.Println("Input:", direction)
		reader := strings.NewReader(input)
		interpreter.AddNewInput(reader)
		interpreter.Run()

		output := strings.TrimSpace(builder.String())
		sensor, err := strconv.Atoi(output)
		if err != nil {
			fmt.Println(err)
			panic("Cannot read robot output")
		}

		_, exists := room[newPositionCandidate]

		switch sensor {
		case 0:
			room[newPositionCandidate] = "#"
		case 1:
			if !exists {
				numDots++
			}
			room[newPositionCandidate] = "."
			position = newPositionCandidate
		case 2:
			room[newPositionCandidate] = "O"
			position = newPositionCandidate
			initialOxygenPosition = newPositionCandidate
		}

		if newPositionCandidate.x < minX {
			minX = newPositionCandidate.x
		}
		if newPositionCandidate.x > maxX {
			maxX = newPositionCandidate.x
		}
		if newPositionCandidate.y < minY {
			minY = newPositionCandidate.y
		}
		if newPositionCandidate.y > maxY {
			maxY = newPositionCandidate.y
		}

		builder.Reset()
	}

	vec := vec2d{}
	for vec.y = minY; vec.y <= maxY; vec.y++ {
		for vec.x = minX; vec.x <= maxX; vec.x++ {
			chr, exists := room[vec]
			if vec.x == 0 && vec.y == 0 {
				fmt.Print("X")
				// } else if vec == position {
				// 	fmt.Println("?")
			} else if exists {
				fmt.Print(chr)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println(len(room))

	steps := 0
	oxygenFilledPositions := []vec2d{initialOxygenPosition}
	fmt.Println(numDots, "spaces to fill!")

	for numDots != 0 {
		for _, position := range oxygenFilledPositions {
			for direction := int64(1); direction <= int64(4); direction++ {
				fillToPosition := moveRobot(position, direction)
				if room[fillToPosition] == "." {
					fmt.Println("Filling", fillToPosition)
					room[fillToPosition] = "O"
					numDots--
					oxygenFilledPositions = append(oxygenFilledPositions, fillToPosition)
				}
			}
		}
		steps++
	}

	vec = vec2d{}
	for vec.y = minY; vec.y <= maxY; vec.y++ {
		for vec.x = minX; vec.x <= maxX; vec.x++ {
			chr, exists := room[vec]
			if vec.x == 0 && vec.y == 0 {
				fmt.Print("X")
				// } else if vec == position {
				// 	fmt.Println("?")
			} else if exists {
				fmt.Print(chr)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println(len(room))
	fmt.Println("Took", steps, "steps to fill with oxygen")
}
