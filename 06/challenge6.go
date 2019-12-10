package main

import (
  "os"
  "bufio"
  "fmt"
  "strings"
)

func main() {
  file, _ := os.Open("input")
  defer file.Close()
  scanner := bufio.NewScanner(file)
  orbits := map[string]map[string]struct{}{}

  for scanner.Scan() {
    split := strings.Split(scanner.Text(), ")")
    parent, child := split[0], split[1]

    _, parentExists := orbits[parent]
    if !parentExists {
      orbits[parent] = map[string]struct{}{}
    }
    orbits[parent][child] = struct{}{}

    _, childExists  := orbits[child]
    if !childExists {
      orbits[child] = map[string]struct{}{}
    }
  }
  count := 0.0
  numOrbits := float64(len(orbits) * len(orbits) * len(orbits))

  for i := range orbits {
    for v := range orbits {
      for w := range orbits {
        count += 1
        if int(count) % 10000000 == 0 {
          percentage := (count / numOrbits)*100
          fmt.Printf("%.2f%%\n", percentage)
        }
        _, viExists := orbits[v][i]
        _, iwExists := orbits[i][w]

        if viExists && iwExists {
          orbits[v][w] = struct{}{}
        }
      }
    }
  }

  count = 0

  for _, childMap := range orbits {
    for _, _ = range childMap {
      count++
    }
  }

  fmt.Println(count)
}
