package main

import (
  "os"
  "bufio"
  "fmt"
  "strings"
  "strconv"
)

type OpFunc func(state *InterpreterState)

const(
  debug = false
)

type Operation struct {
  f OpFunc
  numParams int
  name string
}

type InterpreterState struct {
  operations map[int]Operation
  memory []int
  argumentModeMask []int
  cursor int
  reader *bufio.Reader
}

func newInterpreterState(memory []int) InterpreterState {
  operations := map[int]Operation {
    1: Operation {
      f: addInts,
      numParams: 3,
      name: "ADD",
    },
    2: Operation {
      f: multiplyInts,
      numParams: 3,
      name: "MUL",
    },
    3: Operation {
      f: receiveInput,
      numParams: 1,
      name: "INPUT",
    },
    4: Operation {
      f: writeOutput,
      numParams: 1,
      name: "OUTPUT",
    },
    5: Operation {
      f: jumpIfTrue,
      numParams: 2,
      name: "JUMP-IF-TRUE",
    },
    6: Operation {
      f: jumpIfFalse,
      numParams: 2,
      name: "JUMP-IF-FALSE",
    },
    7: Operation {
      f: lessThan,
      numParams: 3,
      name: "LESS-THAN",
    },
    8: Operation {
      f: equalTo,
      numParams: 3,
      name: "EQUAL-TO",
    },
  }

  interpreterState := InterpreterState{
    operations: operations,
    memory: memory,
    cursor: 0,
    reader: bufio.NewReader(os.Stdin),
  }

  return interpreterState
}

func argMode(argumentIndex int, state *InterpreterState) int {
  if argumentIndex >= len(state.argumentModeMask) {
    return 0
  }

  return state.argumentModeMask[argumentIndex]
}

func readValueAtArgument(argumentIndex int, state *InterpreterState) int {
  firstArgumentIndex := state.cursor + 1
  address := firstArgumentIndex + argumentIndex
  if debug {
    fmt.Println("Reading address", address, ", got value", state.memory[address])
  }
  return state.memory[address]
}

func readArgument(argumentIndex int, state *InterpreterState) int {
  // Hardcoded position mode
  argValue := readValueAtArgument(argumentIndex, state)

  if argMode(argumentIndex, state) == 0 {
    if debug {
      fmt.Println("Reading address", argValue, ", got value", state.memory[argValue])
    }

    return state.memory[argValue]
  } else {
    return argValue
  }
}

func writeArgument(argumentIndex int, value int, state *InterpreterState) {
  // Always position mode
  address := readValueAtArgument(argumentIndex, state)
  if debug {
    fmt.Println("Writing address", address, ", with value", value)
  }

  if argMode(argumentIndex, state) == 0 {
    state.memory[address] = value
  } else {
    panic("Can't write memory in immediate mode")
  }

}

func addInts(state *InterpreterState) {
  lhs, rhs := readArgument(0, state), readArgument(1, state)
  result := lhs + rhs

  fmt.Println(lhs, "+", rhs, "=", result)

  writeArgument(2, result, state)
}

func multiplyInts(state *InterpreterState) {
  lhs, rhs := readArgument(0, state), readArgument(1, state)
  result := lhs * rhs

  fmt.Println(lhs, "*", rhs, "=", result)

  writeArgument(2, result, state)
}

func receiveInput(state *InterpreterState) {
  fmt.Print("Program is requesting input: ")
  text, _ := state.reader.ReadString('\n')
  strippedText := strings.TrimSpace(text)
  toWrite, _ := strconv.Atoi(strippedText)

  writeArgument(0, toWrite, state)
}

func writeOutput(state *InterpreterState) {
  valueToPrint := readArgument(0, state)

  fmt.Println("Program is outputting:", valueToPrint)
}

func jumpIfTrue(state *InterpreterState) {
  test, jumpToAddress := readArgument(0, state), readArgument(1, state)

  if test != 0 {
    state.cursor = jumpToAddress
  }
}

func jumpIfFalse(state *InterpreterState) {
  test, jumpToAddress := readArgument(0, state), readArgument(1, state)

  if test == 0 {
    state.cursor = jumpToAddress
  }
}

func lessThan(state *InterpreterState) {
  lhs, rhs := readArgument(0, state), readArgument(1, state)

  valueToWrite := 0
  if lhs < rhs {
    valueToWrite = 1
  }

  fmt.Println(lhs, "<", rhs, "=>", valueToWrite)

  writeArgument(2, valueToWrite, state)
}

func equalTo(state *InterpreterState) {
  lhs, rhs := readArgument(0, state), readArgument(1, state)

  valueToWrite := 0
  if lhs == rhs {
    valueToWrite = 1
  }

  fmt.Println(lhs, "=", rhs, "=>", valueToWrite)

  writeArgument(2, valueToWrite, state)
}

func parseOpcode(rawOpcode int) (opcode int, mask []int) {
  opcode = rawOpcode % 100

  maskInt := rawOpcode / 100
  mask = []int{}

  for ; maskInt != 0; maskInt /= 10 {
    mask = append(mask, maskInt % 10)
  }

  return
}

func interpretIntCode(memory []int) []int {
  // NOTE TO SELF: COPY THE MEMORY IF YOU WANT TO USE THIS MORE
  // THAN ONCE PER RUN
  state := newInterpreterState(memory)

  if debug {
    fmt.Println("state:", state.memory)
  }

  for state.cursor < len(memory) {
    rawOpcode := memory[state.cursor]

    if rawOpcode == 99 {
      fmt.Println("99: terminating")
      break
    }

    opcode, mask := parseOpcode(rawOpcode)

    state.argumentModeMask = mask

    operation, exists := state.operations[opcode]

    if !exists {
      panic("Invalid opcode: " + strconv.FormatInt(int64(rawOpcode), 10))
    }

    fmt.Println("Running operation", operation.name, "(", rawOpcode, "):", operation)

    oldcursor := state.cursor

    operation.f(&state)

    // If there's a jump, we don't want to modify the state
    if state.cursor == oldcursor {
      state.cursor += operation.numParams + 1
    }

    if debug {
      fmt.Println("state:", state.memory)
    }
  }

  return state.memory
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
  initialProgramCodes := parseProgram("input")

  outputMemory := interpretIntCode(initialProgramCodes)

  if debug {
    fmt.Println(outputMemory)
  }
}
