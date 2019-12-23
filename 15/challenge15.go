package main

import (
	"container/heap"
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

type Body struct {
	position vec2d
	tileType string
	distance int
	index    int
	prev     *Body
}

type PriorityQueue []*Body

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// The least distance should come first
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Body)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil     // stop mem leak by removing the reference
	item.index = -1    // to make sure we don't reuse the index
	*pq = old[0 : n-1] // Make the array one shorter, popping off the last ele
	return item
}

func (pq *PriorityQueue) update(body *Body, position vec2d, distance int) {
	body.position = position
	body.distance = distance
	heap.Fix(pq, body.index)
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

		switch sensor {
		case 0:
			room[newPositionCandidate] = "#"
		case 1:
			room[newPositionCandidate] = "."
			position = newPositionCandidate
		case 2:
			room[newPositionCandidate] = "X"
			position = newPositionCandidate
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
				fmt.Print("O")
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

	pq := PriorityQueue{}
	pqLookup := map[vec2d]*Body{}
	originPosition := vec2d{0, 0}
	var oxygenSystem *Body = nil
	i := 0

	fmt.Print("Populating queue... ")
	for position, tileType := range room {
		distance := math.MaxInt32
		if position == originPosition {
			distance = 0
		}

		if tileType == "#" {
			continue
		}

		pq = append(pq, &Body{
			position: position,
			tileType: tileType,
			distance: distance,
			index:    i,
			prev:     nil,
		})
		pqLookup[position] = pq[len(pq)-1]

		// We need to know where we're measuring _to_,
		// so check if this is a parent of Santa and store
		// that in the santaParent
		if tileType == "X" {
			fmt.Println("Oxygen's position is", position)
			oxygenSystem = pq[i]
		}

		i++
	}
	// Sorts the objects into a heap structure
	heap.Init(&pq)
	fmt.Println("Done")

	for pq.Len() != 0 {
		u := heap.Pop(&pq).(*Body)
		// fmt.Println(u)

		for direction := int64(1); direction <= int64(4); direction++ {
			vPosition := moveRobot(u.position, direction)
			v, exists := pqLookup[vPosition]

			if room[vPosition] == "#" {
				continue
			}

			if !exists {
				// fmt.Println(vPosition)
				panic("Position does not exist!!")
			}

			// No weights in this graph, all paths are len 1
			altPathLen := u.distance + 1
			if altPathLen < v.distance {
				v.distance = altPathLen
				v.prev = u
				heap.Fix(&pq, v.index)
			}
		}
	}

	fmt.Println(oxygenSystem)
}
