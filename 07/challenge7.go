package main

import (
	"fmt"
	"strconv"
	"strings"

	intcode "github.com/JackMc/AdventOfCode2019/lib"
	prmt "github.com/gitchander/permutation"
)

func main() {
	initialMemory := intcode.ParseProgram("input")
	phaseSettings := []string{"0", "1", "2", "3", "4"}
	permutations := prmt.New(prmt.StringSlice(phaseSettings))
	max := 0

	for permutations.Next() {
		output := intcode.SimpleEvaluateProgramAsString(initialMemory, phaseSettings[0]+"\n0\n")
		output = intcode.SimpleEvaluateProgramAsString(initialMemory, phaseSettings[1]+"\n"+output)
		output = intcode.SimpleEvaluateProgramAsString(initialMemory, phaseSettings[2]+"\n"+output)
		output = intcode.SimpleEvaluateProgramAsString(initialMemory, phaseSettings[3]+"\n"+output)
		output = intcode.SimpleEvaluateProgramAsString(initialMemory, phaseSettings[4]+"\n"+output)
		thrusterValue, _ := strconv.Atoi(strings.TrimSpace(output))

		if thrusterValue > max {
			fmt.Println("New max: a =", phaseSettings[0], "b =", phaseSettings[1], "c =", phaseSettings[2], "d=", phaseSettings[3], "e=", phaseSettings[4], "result =", thrusterValue)
			max = thrusterValue
		}
	}
}
