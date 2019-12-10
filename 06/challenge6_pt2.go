package main

import (
  "os"
  "bufio"
  "fmt"
  "strings"
  "container/heap"
  "math"
)

type Body struct {
  name string
  distance int
  index int
  prev *Body
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
  old[n-1] = nil // stop mem leak by removing the reference
  item.index = -1 // to make sure we don't reuse the index
  *pq = old[0:n-1] // Make the array one shorter, popping off the last ele
  return item
}

func (pq *PriorityQueue) update(body *Body, name string, distance int) {
  body.name = name
  body.distance = distance
  heap.Fix(pq, body.index)
}

func main() {
  file, _ := os.Open("input")
  defer file.Close()
  scanner := bufio.NewScanner(file)
  // Maps parent -> map(child -> struct{})
  orbits := map[string]map[string]struct{}{}

  fmt.Print("Creating graph... ")
  for scanner.Scan() {
    split := strings.Split(scanner.Text(), ")")
    parent, child := split[0], split[1]

    _, parentExists := orbits[parent]
    if !parentExists {
      orbits[parent] = map[string]struct{}{}
    }
    _, childExists  := orbits[child]
    if !childExists {
      orbits[child] = map[string]struct{}{}
    }

    orbits[parent][child] = struct{}{}
    orbits[child][parent] = struct{}{}
  }
  fmt.Println("Done")

  pq := make(PriorityQueue, len(orbits))
  bodies := map[string]*Body{}
  i := 0
  var santaParent *Body = nil

  fmt.Print("Populating queue... ")
  for planetName, children := range orbits {
    _, isParentOfYou := children["YOU"]

    distance := math.MaxInt32
    if isParentOfYou {
      distance = 0
    }

    pq[i] = &Body {
      name: planetName,
      distance: distance,
      index: i,
      prev: nil,
    }
    bodies[planetName] = pq[i]

    // We need to know where we're measuring _to_,
    // so check if this is a parent of Santa and store
    // that in the santaParent
    _, isParentOfSanta := children["SAN"]
    if isParentOfSanta {
      fmt.Println("Santa's parent is", planetName)
      santaParent = bodies[planetName]
    }

    i++
  }
  // Sorts the objects into a heap structure
  heap.Init(&pq)
  fmt.Println("Done")

  for pq.Len() != 0 {
    u := heap.Pop(&pq).(*Body)

    for neighbourPlanetName, _ := range orbits[u.name] {
      v := bodies[neighbourPlanetName]
      // No weights in this graph, all paths are len 1
      altPathLen := u.distance + 1
      if altPathLen < v.distance {
        v.distance = altPathLen
        v.prev = u
        heap.Fix(&pq, v.index)
      }
    }
  }

  fmt.Println(santaParent)
}
