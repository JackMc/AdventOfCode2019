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
  // Maps (x,y) pairs to which wire they're in (0-indexed)
  wirePoints := map[[2]int]int{}
  collisions := [][2]int{}
  wireNum := 0

  for ; scanner.Scan(); wireNum++ {
    wireStr := scanner.Text()
    wireSplit := strings.Split(wireStr, ",")
    // Start at 0,0 for drawing
    cursor := [2]int{0,0}

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
        value, ok := wirePoints[cursor]

        if ok && value != wireNum {
          collisions = append(collisions, cursor)
        }

        // Avoid counting a collision with ourselves
        wirePoints[cursor] = wireNum
      }
    }
  }

  origin := [2]int{0,0}
  // hardcoded intmax...
  min := 2147483647
  fmt.Println(collisions)
  for _, collision := range collisions {
    dist := manhattan(collision, origin)
    if dist < min {
      min = dist
    }
    fmt.Println("Manhattan for", collision, ":", manhattan(collision, origin))
  }

  fmt.Println(min)
  fmt.Println(len(wirePoints))
}
