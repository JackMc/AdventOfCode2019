package main

import (
  "os"
  "bufio"
  "fmt"
  "strings"
  "strconv"
)

func intCode(inputProgramCodes []int) []int {
  memory := make([]int, len(inputProgramCodes))
  copy(memory, inputProgramCodes)

  cursor := 0

  for ; cursor < len(memory); cursor += 4 {
    opcode := memory[cursor]
    if opcode == 99 {
      fmt.Println("99: terminating")
      break
    }
    lhsIdx, rhsIdx := memory[cursor + 1], memory[cursor + 2]
    fmt.Println("Reading address", cursor + 1, ", got value", lhsIdx)
    fmt.Println("Reading address", cursor + 2, ", got value", rhsIdx)
    fmt.Println("Reading address", lhsIdx, ", got value", memory[lhsIdx])
    fmt.Println("Reading address", rhsIdx, ", got value", memory[rhsIdx])
    resultIdx := memory[cursor + 3]
    val1, val2 := memory[lhsIdx], memory[rhsIdx]

    if opcode == 1 {
      fmt.Println(val1, "+", val2, "=", val1 + val2)
      memory[resultIdx] = val1 + val2
    }
    if opcode == 2 {
      fmt.Println(val1, "+", val2, "=", val1 * val2)
      memory[resultIdx] = val1 * val2
    }

    fmt.Println("state:", memory)
  }

  return memory
}

func parseProgram(filename string) []int {
  file, _ := os.Open(filename)
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  programText := scanner.Text()
  programStrings := strings.Split(programText, ",")
  initialProgramCodes := []int{}

  for _, convertibleInt := range programStrings {
    convertedInt, _ := strconv.Atoi(convertibleInt)
    initialProgramCodes = append(initialProgramCodes, convertedInt)
  }

  return initialProgramCodes
}

func main() {
  // Pt 1

  initialProgramCodes := parseProgram("pt1_input")

  outputMemory := intCode(initialProgramCodes)

  fmt.Println(outputMemory)

  // Pt 2

  // initialProgramCodes = parseProgram("pt2_input")
  //
  // for noun := 0; noun < 100; noun ++ {
  //   initialProgramCodes[1] = noun
  //   for verb := 0; verb < 100; verb ++ {
  //     initialProgramCodes[2] = verb
  //
  //     outputMemory := intCode(initialProgramCodes)
  //
  //     if outputMemory[0] == 19690720 {
  //       fmt.Println("Noun:", noun, "Verb:", verb, "Composite:", (noun * 100) + verb)
  //       break
  //     }
  //   }
  // }

  fmt.Println("Done!")
}
