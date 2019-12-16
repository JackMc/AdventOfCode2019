package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type vec3d struct {
	x, y, z int
}

type planet struct {
	position vec3d
	velocity vec3d
}

func (planet *planet) timeStep() {
	planet.position.x += planet.velocity.x
	planet.position.y += planet.velocity.y
	planet.position.z += planet.velocity.z
}

func (planet planet) potentialEnergy() int {
	return abs(planet.position.x) + abs(planet.position.y) + abs(planet.position.z)
}

func (planet planet) kineticEnergy() int {
	return abs(planet.velocity.x) + abs(planet.velocity.y) + abs(planet.velocity.z)
}

func (planet planet) totalEnergy() int {
	return planet.kineticEnergy() * planet.potentialEnergy()
}

func (planetA *planet) changeVelocity(planetB *planet) {
	if planetA.position.x < planetB.position.x {
		planetA.velocity.x++
		planetB.velocity.x--
	} else if planetA.position.x > planetB.position.x {
		planetA.velocity.x--
		planetB.velocity.x++
	}

	if planetA.position.y < planetB.position.y {
		planetA.velocity.y++
		planetB.velocity.y--
	} else if planetA.position.y > planetB.position.y {
		planetA.velocity.y--
		planetB.velocity.y++
	}

	if planetA.position.z < planetB.position.z {
		planetA.velocity.z++
		planetB.velocity.z--
	} else if planetA.position.z > planetB.position.z {
		planetA.velocity.z--
		planetB.velocity.z++
	}
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

func abs(x int) int {
	if x < 0 {
		x = -x
	}

	return x
}

func (a *planet) distance(b planet) float64 {
	dx := a.position.x - b.position.x
	dy := a.position.x - b.position.x

	return math.Sqrt(float64(dx*dx + dy*dy))
}

func between(a, b, c planet) bool {
	left := min(a.position.x, b.position.x)
	right := max(a.position.x, b.position.x)
	top := min(a.position.y, b.position.y)
	bottom := max(a.position.y, b.position.y)

	if c.position.x >= left && c.position.x <= right && c.position.y >= top && c.position.y <= bottom {
		return true
	}

	return false
}

func main() {
	fileContentsBytes, err := ioutil.ReadFile("input")

	if err != nil {
		panic("Couldn't read input file")
	}

	fileContents := strings.TrimSpace(string(fileContentsBytes))
	splitByNewline := strings.Split(fileContents, "\n")
	planets := make([]planet, len(splitByNewline))
	states := [][]planet{}
	re := regexp.MustCompile("-?[0-9]+")

	for i := 0; i < len(planets); i++ {
		fmt.Println("Matching against", splitByNewline[i])
		match := re.FindAll([]byte(splitByNewline[i]), -1)
		x, err := strconv.Atoi(string(match[0]))
		if err != nil {
			panic("x invalid")
		}
		y, err := strconv.Atoi(string(match[1]))
		if err != nil {
			panic("y invalid")
		}
		z, err := strconv.Atoi(string(match[2]))
		if err != nil {
			panic("z invalid")
		}

		planets[i] = planet{position: vec3d{x: x, y: y, z: z}, velocity: vec3d{}}
	}

	fmt.Println("Step 0:")
	// Simulate velocity
	for _, planet := range planets {
		planet.position.x += planet.velocity.x
		planet.position.y += planet.velocity.y
		planet.position.z += planet.velocity.z
		fmt.Println("pos =", planet.position, "vel =", planet.velocity)
	}

	for timeStep := 1; timeStep <= 1000; timeStep++ {
		// Simulate gravity looping through each unique pair of planets
		for i := 0; i < len(planets); i++ {
			for j := i + 1; j < len(planets); j++ {
				planetA := &planets[i]
				planetB := &planets[j]

				planetA.changeVelocity(planetB)
			}
		}

		if timeStep%10 == 0 {
			fmt.Println("Step", timeStep, ":")
		}
		// Simulate velocity
		for i := range planets {
			planet := &planets[i]

			planet.position.x += planet.velocity.x
			planet.position.y += planet.velocity.y
			planet.position.z += planet.velocity.z
			if timeStep%10 == 0 {
				fmt.Println("pos =", planet.position, "vel =", planet.velocity)
			}
		}
	}

	fmt.Println("----- Final results -----")
	sum := 0
	for _, planet := range planets {
		fmt.Println("pot =", planet.potentialEnergy(), "kin =", planet.kineticEnergy(), "total =", planet.totalEnergy())
		sum += planet.totalEnergy()
	}

	fmt.Println(sum)
}
