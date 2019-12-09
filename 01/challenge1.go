package main

import (
  "os"
  "bufio"
  "strconv"
  "fmt"
  "math"
)

func calculateFuel(mass int) int {
  return int(math.Floor(float64(mass) / 3) - 2)
}

func main() {
  file, _ := os.Open("input")
  defer file.Close()
  scanner := bufio.NewScanner(file)
  sum := 0

  for scanner.Scan() {
    totalFuelRequired := 0
    massThisCycle, _ := strconv.Atoi(scanner.Text())
    // Placeholder
    fuelRequiredThisCycle := calculateFuel(massThisCycle)
    massThisCycle = fuelRequiredThisCycle
    i := 0

    for fuelRequiredThisCycle > 0 {
      fuelRequiredThisCycle = calculateFuel(massThisCycle)
      massThisCycle = fuelRequiredThisCycle
      totalFuelRequired += fuelRequiredThisCycle
    }

    sum += totalFuelRequired
  }

  fmt.Println(sum)
}
