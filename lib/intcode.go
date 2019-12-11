package intcode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type opFunc func(state *InterpreterState) (err error)

const (
	debug = false
)

var (
	operations = map[int]operation{
		1: operation{
			f:         addInts,
			numParams: 3,
			name:      "ADD",
		},
		2: operation{
			f:         multiplyInts,
			numParams: 3,
			name:      "MUL",
		},
		3: operation{
			f:         receiveInput,
			numParams: 1,
			name:      "INPUT",
		},
		4: operation{
			f:         writeOutput,
			numParams: 1,
			name:      "OUTPUT",
		},
		5: operation{
			f:         jumpIfTrue,
			numParams: 2,
			name:      "JUMP-IF-TRUE",
		},
		6: operation{
			f:         jumpIfFalse,
			numParams: 2,
			name:      "JUMP-IF-FALSE",
		},
		7: operation{
			f:         lessThan,
			numParams: 3,
			name:      "LESS-THAN",
		},
		8: operation{
			f:         equalTo,
			numParams: 3,
			name:      "EQUAL-TO",
		},
		9: operation{
			f:         adjustRelativeBase,
			numParams: 1,
			name:      "ADJUST-RELATIVE-BASE",
		},
	}
	NeedMoreInput = interpreterError{description: "The interpreter requires more input"}
)

type interpreterError struct {
	description string
}

func (e interpreterError) Error() string {
	return e.description
}

type operation struct {
	f         opFunc
	numParams int
	name      string
}

type InterpreterState struct {
	operations       map[int]operation
	memory           []int
	argumentModeMask []int
	cursor           int
	relativeBase     int
	interactive      bool
	terminated       bool
	source           io.ByteReader
	sink             io.Writer
}

func NewInteractiveInterpreter(filename string) InterpreterState {
	state := NewInterpreterFromFile(filename, bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	state.interactive = true

	return state
}

func NewInterpreterFromFile(filename string, source io.ByteReader, sink io.Writer) InterpreterState {
	memory := ParseProgram(filename)

	return NewInterpreter(memory, source, sink)
}

func NewInterpreter(initialMemory []int, source io.ByteReader, sink io.Writer) InterpreterState {
	memory := make([]int, len(initialMemory))
	copy(memory, initialMemory)

	interpreterState := InterpreterState{
		operations:   operations,
		memory:       memory,
		cursor:       0,
		relativeBase: 0,
		source:       source,
		sink:         sink,
		interactive:  false,
	}

	return interpreterState
}

func SimpleEvaluateProgramAsString(initialMemory []int, input string) (output string) {
	sink := &strings.Builder{}
	interpreter := NewInterpreter(initialMemory, strings.NewReader(input), sink)
	interpreter.Run()

	return sink.String()
}

func EvaluateProgramWithStringInput(initialMemory []int, input string) (interpreter InterpreterState, err error) {
	sink := &strings.Builder{}
	interpreter = NewInterpreter(initialMemory, strings.NewReader(input), sink)
	_, err = interpreter.Run()

	return interpreter, err
}

func (state *InterpreterState) checkResize(absoluteAddress int) {
	if absoluteAddress > len(state.memory)-1 {
		oldMemory := state.memory
		newBufferLength := absoluteAddress + 1
		state.memory = make([]int, newBufferLength)
		copy(state.memory, oldMemory)
	}
}

func (state *InterpreterState) readAddress(absoluteAddress int) int {
	state.checkResize(absoluteAddress)

	return state.memory[absoluteAddress]
}

func (state *InterpreterState) writeAddress(absoluteAddress int, value int) {
	state.checkResize(absoluteAddress)

	state.memory[absoluteAddress] = value
}

func (state *InterpreterState) readValueAtArgument(argumentIndex int) int {
	firstArgumentIndex := state.cursor + 1
	address := firstArgumentIndex + argumentIndex
	value := state.readAddress(address)
	if debug {
		if state.interactive {
			fmt.Println("Reading address", address, ", got value", value)
		}
	}

	return value
}

func (state *InterpreterState) readArgument(argumentIndex int) int {
	// Hardcoded position mode
	argValue := state.readValueAtArgument(argumentIndex)
	argMode := state.argMode(argumentIndex)

	if argMode == 0 {
		value := state.readAddress(argValue)
		if debug {
			if state.interactive {
				fmt.Println("Reading address", argValue, ", got value", value)
			}
		}

		return value
	} else if argMode == 1 {
		return argValue
	} else if argMode == 2 {
		return state.readAddress(state.relativeBase + argValue)
	}

	return argValue
}

func (state *InterpreterState) writeArgument(argumentIndex int, value int) {
	argValue := state.readValueAtArgument(argumentIndex)
	argMode := state.argMode(argumentIndex)
	address := -1

	if argMode == 0 {
		address = argValue
	} else if argMode == 1 {
		panic("Can't write memory in immediate mode")
	} else if argMode == 2 {
		address = state.relativeBase + argValue
	}

	state.writeAddress(address, value)

	if debug {
		if state.interactive {
			fmt.Println("Writing address", address, ", with value", value)
		}
	}
}

func (state *InterpreterState) argMode(argumentIndex int) int {
	if argumentIndex >= len(state.argumentModeMask) {
		return 0
	}

	return state.argumentModeMask[argumentIndex]
}

func addInts(state *InterpreterState) error {
	lhs, rhs := state.readArgument(0), state.readArgument(1)
	result := lhs + rhs

	if state.interactive {
		fmt.Println(lhs, "+", rhs, "=", result)
	}

	state.writeArgument(2, result)

	return nil
}

func multiplyInts(state *InterpreterState) error {
	lhs, rhs := state.readArgument(0), state.readArgument(1)
	result := lhs * rhs

	if state.interactive {
		fmt.Println(lhs, "*", rhs, "=", result)
	}

	state.writeArgument(2, result)

	return nil
}

func readStringUntilNewline(reader io.ByteReader) (string, error) {
	// 16 is the highest expected size
	bytes := []byte{}
	bytesRead := 0

	for byteWeRead, err := reader.ReadByte(); byteWeRead != byte('\n'); bytesRead++ {
		if err == io.EOF {
			return "", NeedMoreInput
		}

		bytes = append(bytes, byteWeRead)

		byteWeRead, err = reader.ReadByte()
	}

	return string(bytes[:bytesRead]), nil
}

func receiveInput(state *InterpreterState) error {
	if state.interactive {
		fmt.Print("Program is requesting input: ")
	}
	// Truncate only to bytes read
	input, err := readStringUntilNewline(state.source)

	if err != nil {
		return err
	}

	strippedText := strings.TrimSpace(input)
	toWrite, _ := strconv.Atoi(strippedText)

	state.writeArgument(0, toWrite)

	return nil
}

func writeOutput(state *InterpreterState) error {
	valueToWrite := state.readArgument(0)

	decimalString := strconv.FormatInt(int64(valueToWrite), 10) + "\n"

	state.sink.Write([]byte(decimalString))

	if state.interactive {
		fmt.Println("Program is outputting:", valueToWrite)
	}

	return nil
}

func jumpIfTrue(state *InterpreterState) error {
	test, jumpToAddress := state.readArgument(0), state.readArgument(1)

	if test != 0 {
		state.cursor = jumpToAddress
	}

	return nil
}

func jumpIfFalse(state *InterpreterState) error {
	test, jumpToAddress := state.readArgument(0), state.readArgument(1)

	if test == 0 {
		state.cursor = jumpToAddress
	}

	return nil
}

func lessThan(state *InterpreterState) error {
	lhs, rhs := state.readArgument(0), state.readArgument(1)

	valueToWrite := 0
	if lhs < rhs {
		valueToWrite = 1
	}

	if state.interactive {
		fmt.Println(lhs, "<", rhs, "=>", valueToWrite)
	}

	state.writeArgument(2, valueToWrite)

	return nil
}

func equalTo(state *InterpreterState) error {
	lhs, rhs := state.readArgument(0), state.readArgument(1)

	valueToWrite := 0
	if lhs == rhs {
		valueToWrite = 1
	}

	if state.interactive {
		fmt.Println(lhs, "=", rhs, "=>", valueToWrite)
	}

	state.writeArgument(2, valueToWrite)

	return nil
}

func adjustRelativeBase(state *InterpreterState) error {
	adjustment := state.readArgument(0)
	state.relativeBase += adjustment

	return nil
}

func parseOpcode(rawOpcode int) (opcode int, mask []int) {
	opcode = rawOpcode % 100

	maskInt := rawOpcode / 100
	mask = []int{}

	for ; maskInt != 0; maskInt /= 10 {
		mask = append(mask, maskInt%10)
	}

	return
}

func (state *InterpreterState) AddNewInput(source io.ByteReader) {
	state.source = source
}

func (state *InterpreterState) IsTerminated() bool {
	return state.terminated
}

func (state *InterpreterState) Sink() io.Writer {
	return state.sink
}

func (state *InterpreterState) Run() (finalMemory []int, err error) {
	if state.terminated {
		panic("Trying to run a terminated interpreter")
	}

	if debug {
		if state.interactive {
			fmt.Println("state:", state.memory)
		}
	}

	for state.cursor < len(state.memory) {
		rawOpcode := state.readAddress(state.cursor)
		if rawOpcode == 99 {
			if state.interactive {
				fmt.Println("99: terminating")
			}
			state.terminated = true
			break
		}

		opcode, mask := parseOpcode(rawOpcode)

		state.argumentModeMask = mask

		operation, exists := state.operations[opcode]

		if !exists {
			panic("Invalid opcode: " + strconv.FormatInt(int64(rawOpcode), 10))
		}

		if state.interactive {
			fmt.Println("Running operation", operation.name, "(", rawOpcode, "):", operation)
		}

		oldcursor := state.cursor

		err := operation.f(state)

		if err != nil {
			return state.memory, err
		}

		// If there's a jump, we don't want to modify the state
		if state.cursor == oldcursor {
			state.cursor += operation.numParams + 1
		}

		if debug {
			if state.interactive {
				fmt.Println("state:", state.memory)
			}
		}
	}

	return state.memory, nil
}

func ParseProgram(filename string) []int {
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
