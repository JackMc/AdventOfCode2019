package main

import (
  "os"
  "bufio"
  "fmt"
  "strings"
  "strconv"
)

func abs(x int) int {
  if x < 0 {
    return -x
  }
  return x
}

func manhattan(point1 [2]int, point2 [2]int) int {
  return abs(point1[0] - point2[0]) + abs(point1[1] - point2[1])
}

func main() {
  file, _ := os.Open("input")
  defer file.Close()
  scanner := bufio.NewScanner(file)
  // Maps (x,y) -> (wire1length, wire2length, lastSeenWire)
  wirePoints := map[[2]int][2]int{}
  collisions := [][2]int{}
  wireNum := 0

  for scanner.Scan() {
    wireStr := scanner.Text()
    wireSplit := strings.Split(wireStr, ",")
    // Start at 0,0 for drawing
    cursor := [2]int{0,0}
    totalDist := 0

    for _, instruction := range wireSplit {
      directionChr := instruction[0:1]
      magnitude, _ := strconv.Atoi(instruction[1:])
      xDirection, yDirection := 0, 0

      switch directionChr {
      case "R":
        xDirection = 1
      case "L":
        xDirection = -1
      case "U":
        yDirection = 1
      case "D":
        yDirection = -1
      default:
        panic("Invalid direction char")
      }

      for i := 0; i < magnitude; i++ {
        cursor[0] += xDirection
        cursor[1] += yDirection
        totalDist += 1

        value, exists := wirePoints[cursor]

        // If it doesn't exist, we're the first to see this space
        if !exists {
          wirePoints[cursor] = [2]int{-1,-1}
        } else if value[wireNum] == -1 {
          // The other wire has been here, but we have not. Collision!
          collisions = append(collisions, cursor)
        }

        // We may have modified the wirePoint so we refetch it
        value = wirePoints[cursor]
        value[wireNum] = totalDist
        wirePoints[cursor] = value
      }
    }

    wireNum++
  }

  // hardcoded intmax...
  min := 2147483647
  fmt.Println(collisions)
  for _, collision := range collisions {
    dist := wirePoints[collision][0] + wirePoints[collision][1]
    fmt.Println("Collision:", collision)
    fmt.Println("Distance wire A:", wirePoints[collision][0])
    fmt.Println("Distance wire B:", wirePoints[collision][1])
    fmt.Println("Total distance:", dist)

    if dist < min {
      min = dist
    }
  }

  fmt.Println(min)
}
