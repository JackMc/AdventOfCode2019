package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
)

type point struct {
	x int
	y int
}

type game struct {
	ball   point
	paddle point
	score  int
	screen [24][40]int
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (a point) distance(b point) float64 {
	dx := a.x - b.x
	dy := a.x - b.x

	return math.Sqrt(float64(dx*dx + dy*dy))
}

const (
	empty            = 0
	wall             = 1
	block            = 2
	horizontalPaddle = 3
	ball             = 4
)

var (
	tileToChar = map[int]string{
		empty:            " ",
		wall:             "W",
		block:            "B",
		horizontalPaddle: "_",
		ball:             "X",
	}
)

func gameState(interpreter intcode.InterpreterState) game {
	outputBuilder := interpreter.Sink().(*strings.Builder)
	output := outputBuilder.String()

	split := strings.Split(strings.TrimSpace(output), "\n")
	game := game{}
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
				game.score = tileType
			} else {
				if tileType == ball {
					game.ball = point{x: xPos, y: yPos}
				}

				if tileType == horizontalPaddle {
					game.paddle = point{x: xPos, y: yPos}
				}

				game.screen[yPos][xPos] = tileType
			}
		}
	}
	return game
}

func printScreen(game *game) {
	for _, row := range game.screen {
		for _, tileType := range row {
			fmt.Print(tileToChar[tileType])
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("Score:", game.score)
}

func main() {
	program := intcode.ParseProgram("input")
	// reader := bufio.NewReader(os.Stdin)
	// Unlimited quarters mode!
	program[0] = 2
	interpreter, err := intcode.EvaluateProgramWithStringInput(program, "")
	if err != intcode.NeedMoreInput {
		panic("Interpreter quit early")
	}
	var lastGame *game = nil

	for !interpreter.IsTerminated() {
		game := gameState(interpreter)
		// fmt.Println(game, lastGame)
		paddleMove := "0\n"

		if lastGame != nil {
			ballYVelocity := game.ball.y - lastGame.ball.y
			ballXVelocity := game.ball.x - lastGame.ball.x
			paddleBallYDiff := game.paddle.y - game.ball.y

			if ballYVelocity > 0 {
				// The ball is actively moving towards us
				predictedBallPosition := point{}

				predictedBallPosition.x = game.ball.x + (ballXVelocity * paddleBallYDiff)
				predictedBallPosition.y = game.ball.y + (ballYVelocity * paddleBallYDiff)

				// fmt.Println("XVel:", ballXVelocity)
				// fmt.Println("YVel:", ballYVelocity)
				// fmt.Println("Paddle Y diff:", paddleBallYDiff)
				// fmt.Println("Predicted ball position:", predictedBallPosition)

				// input, err := reader.ReadString('\n')

				// if err != nil {
				// 	fmt.Println(err)
				// 	panic("Error while getting input")
				// }

				predictedDeltaX := abs(predictedBallPosition.x - game.paddle.x)
				// fmt.Println("Predicted delta:", predictedDeltaX)
				if predictedDeltaX > 1 {
					// If we stay where we are, we won't hit the ball
					paddleMove = strconv.FormatInt(int64(ballXVelocity), 10) + "\n"
				}
			} else {
				// The ball is moving up, but we should position ourselves to be
				// near it. We're not predicting the balls position, just moving
				// to where it is right now
				deltaX := game.ball.x - game.paddle.x
				if deltaX != 0 {
					paddleMove = strconv.FormatInt(int64(ballXVelocity), 10) + "\n"
				}
			}
		}
		// fmt.Print("Moving:", paddleMove)
		printScreen(&game)
		time.Sleep(66 * time.Millisecond)
		interpreter.AddNewInput(strings.NewReader(paddleMove))
		_, err = interpreter.Run()

		lastGame = &game
	}

	game := gameState(interpreter)
	printScreen(&game)
	// fmt.Println(game)
}
