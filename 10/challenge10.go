package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type point struct {
	x int
	y int
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (a *point) distance(b point) float64 {
	dx := a.x - b.x
	dy := a.x - b.x

	return math.Sqrt(float64(dx*dx + dy*dy))
}

func main() {
	fileContentsBytes, err := ioutil.ReadFile("sample")

	if err != nil {
		panic("Couldn't read input file")
	}

	fileContents := string(fileContentsBytes)
	splitByNewline := strings.Split(fileContents, "\n")
	grid := make([][]string, len(splitByNewline))
	asteroids := []point{}

	for yCoord, rowString := range splitByNewline {
		grid[yCoord] = strings.Split(rowString, "")

		for xCoord, itemInPosition := range grid[yCoord] {
			if itemInPosition == "#" {
				asteroids = append(asteroids, point{xCoord, yCoord})
			}
		}
	}

	reachableAsteroids := map[point]map[point]bool{}

	for _, asteroid := range asteroids {
		reachableAsteroids[asteroid] = map[point]bool{}
	}

	for _, originAsteroid := range asteroids {
		for _, destinationAsteroid := range asteroids {
			if originAsteroid == destinationAsteroid {
				continue
			}

			// oppositeVisible, exists := reachableAsteroids[destinationAsteroid][originAsteroid]

			// if exists {
			// 	reachableAsteroids[originAsteroid][destinationAsteroid] = oppositeVisible
			// 	continue
			// }

			visible := true

			originX, originY := originAsteroid.x, originAsteroid.y
			destX, destY := destinationAsteroid.x, destinationAsteroid.y
			// Compute the linear equation between the two coordinates
			xDelta := originX - destX
			yDelta := originY - destY
			// m = dy / dx
			if xDelta != 0 {
				slope := float64(yDelta) / float64(xDelta)
				// b = y - mx
				yIntercept := float64(originY) - slope*float64(originX)
				for _, potentiallyObstructingAsteroid := range asteroids {
					if potentiallyObstructingAsteroid == originAsteroid || potentiallyObstructingAsteroid == destinationAsteroid {
						continue
					}
					// y = mx + b
					yForX := float64(potentiallyObstructingAsteroid.x)*slope + yIntercept
					howCloseToLine := yForX - float64(potentiallyObstructingAsteroid.y)

					if howCloseToLine == 0 {
						// We know that it's on the same line segment. Is it
						// between the two points?
						visible = false
						break
					}
				}
			} else {
				// Linear equations don't work if it's straight up and down (xDelta == 0),
				// so we manually have to check all the things on the same
				// ycoord
				beginAt := min(originY, destY)
				endAt := max(originY, destY)

				for checkY := beginAt + 1; checkY < endAt; checkY++ {
					if grid[checkY][originX] == "#" {
						visible = false
						break
					}
				}
			}

			reachableAsteroids[originAsteroid][destinationAsteroid] = visible
		}
	}

	for origin, visibleAsteroidMap := range reachableAsteroids {
		count := 0
		for _, isVisible := range visibleAsteroidMap {
			if isVisible {
				count++
			}
		}

		fmt.Println(origin, "can see", count, "asteroids")
	}

	fmt.Println(reachableAsteroids[point{4, 0}])
}
