package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func onesDigit(n int) int {
	return abs(n % 10)
}

func main() {
	basePattern := [4]int{0, 1, 0, -1}

	contents, err := ioutil.ReadFile("input")
	if err != nil {
		fmt.Println(err)
		panic("Can't read file")
	}
	trimContents := strings.TrimSpace(string(contents))
	split := strings.Split(trimContents, "")
	input := make([]int, len(split))

	for i, numStr := range split {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println(err)
			panic("Error while atoiing input")
		}
		input[i] = num
	}

	fmt.Println("Input:", input)

	output := make([]int, len(input))

	for phaseNumber := 1; phaseNumber <= 100; phaseNumber++ {
		for outputPosition := 0; outputPosition < len(output); outputPosition++ {
			val := 0

			indexIntoPattern := 0
			// Skip the first 0
			currentPatternIndexUsedCount := 1

			for _, digit := range input {
				// Account for outputPosition being 0-indexed
				// In the first spot, each pattern index is used once, in
				// the second 2 times, 3rd 3 times, etc.
				if currentPatternIndexUsedCount == outputPosition+1 {
					indexIntoPattern = (indexIntoPattern + 1) % 4
					currentPatternIndexUsedCount = 0
				}

				// fmt.Print(digit, "*", basePattern[indexIntoPattern], " + ")
				val += digit * basePattern[indexIntoPattern]
				currentPatternIndexUsedCount++
			}

			digit := onesDigit(val)
			// fmt.Println(" =", digit, "(", val, ")")
			output[outputPosition] = digit
		}

		fmt.Println("After phase", phaseNumber, ": ", output)
		input = output
	}

	// Print the last 8 digits
	fmt.Println("After 100 phases, first 8 digits are:", input[0:8])
}
