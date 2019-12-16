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

func between(a point, b point, c point) bool {
	left := min(a.x, b.x)
	right := max(a.x, b.x)
	top := min(a.y, b.y)
	bottom := max(a.y, b.y)

	if c.x >= left && c.x <= right && c.y >= top && c.y <= bottom {
		return true
	}

	return false
}

func main() {
	fileContentsBytes, err := ioutil.ReadFile("smallsample")

	if err != nil {
		panic("Couldn't read input file")
	}

	fileContents := string(fileContentsBytes)
	splitByNewline := strings.Split(fileContents, "\n")
	grid := make([][]string, len(splitByNewline))
	asteroids := map[point]struct{}{}

	for yCoord, rowString := range splitByNewline {
		grid[yCoord] = strings.Split(rowString, "")

		for xCoord, itemInPosition := range grid[yCoord] {
			if itemInPosition == "#" {
				asteroids[point{xCoord, yCoord}] = struct{}{}
			}
		}
	}

	reachableAsteroids := map[point]map[point]bool{}

	for asteroid := range asteroids {
		reachableAsteroids[asteroid] = map[point]bool{}
	}

	for originAsteroid := range asteroids {
		for destinationAsteroid := range asteroids {
			if originAsteroid == destinationAsteroid {
				continue
			}

			oppositeVisible, exists := reachableAsteroids[destinationAsteroid][originAsteroid]

			if exists {
				reachableAsteroids[originAsteroid][destinationAsteroid] = oppositeVisible
				continue
			}

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
				for potentiallyObstructingAsteroid := range asteroids {
					if potentiallyObstructingAsteroid == originAsteroid || potentiallyObstructingAsteroid == destinationAsteroid {
						continue
					}
					// y = mx + b
					yForX := float64(potentiallyObstructingAsteroid.x)*slope + yIntercept
					howCloseToLine := yForX - float64(potentiallyObstructingAsteroid.y)

					if math.Abs(howCloseToLine) < 0.0001 && between(originAsteroid, destinationAsteroid, potentiallyObstructingAsteroid) {
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

	asteroidWithStation := point{-1, -1}
	maxCount := 0

	for origin, visibleAsteroidMap := range reachableAsteroids {
		count := 0
		for _, isVisible := range visibleAsteroidMap {
			if isVisible {
				count++
			}
		}

		if count > maxCount {
			asteroidWithStation = origin
			maxCount = count
		}
	}

	fmt.Println(asteroidWithStation, maxCount)

	// PART 2
	// Angle, 0 to 360.
	angle := 0

	for len(asteroids) > 0 {
		if angle == 0 || angle == 180 {

		} else {

		}
		angle = (angle + 1) % 360
	}
}
