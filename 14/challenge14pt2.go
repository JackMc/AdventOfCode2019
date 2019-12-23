package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

type formula struct {
	result    chemical
	reactants []chemical
}

func newFormula(formulaString string) formula {
	arrowSplit := strings.Split(formulaString, "=>")
	if len(arrowSplit) != 2 {
		panic("Length of arrow split is not 2")
	}

	reactantsStr, resultStr := arrowSplit[0], arrowSplit[1]

	reactantStrings := strings.Split(reactantsStr, ", ")
	reactants := make([]chemical, len(reactantStrings))

	for index, reactant := range reactantStrings {
		reactants[index] = newChemical(reactant)
	}

	result := newChemical(resultStr)

	return formula{result, reactants}
}

type chemical struct {
	name   string
	amount int64
}

func newChemical(chemicalString string) chemical {
	chemicalSplit := strings.Split(strings.TrimSpace(chemicalString), " ")
	if len(chemicalSplit) != 2 {
		panic("Incorrect chemical split")
	}

	amountStr, name := chemicalSplit[0], chemicalSplit[1]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		//fmt.Println(err)
		panic("Error while atoiing chemical name")
	}

	return chemical{name, int64(amount)}
}

func readInput() map[string]formula {
	// maps results to reactants
	formulas := map[string]formula{}

	inputBytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		//fmt.Println(err)
		panic("Error reading input")
	}

	inputLines := strings.Split(string(inputBytes), "\n")

	for _, line := range inputLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		formula := newFormula(line)
		resultName := formula.result.name
		formulas[resultName] = formula
	}
	//fmt.Println(formulas)
	return formulas
}

func solveRecursiveImpl(need chemical, formulas map[string]formula, leftover map[string]int64) (map[string]int64, int64) {
	if need.name == "ORE" {
		return leftover, need.amount
	}

	formula := formulas[need.name]

	sumOre := int64(0)
	ore := int64(0)

	times := int64(math.Ceil(float64(need.amount) / float64(formula.result.amount)))

	//fmt.Println("Satisfying requirement", need, "with formula", formula)

	//fmt.Println(leftover)
	// if the amount needed is less than the amount left over, just take the amount needed
	leftoverReactant := leftover[need.name]
	// The required amount of reactant to require at least one less execution
	goalReactant := formula.result.amount * (times - 1)
	needIfAllLeftoverUsed := max(0, need.amount-leftoverReactant)
	if needIfAllLeftoverUsed <= goalReactant {
		// originalNeed := need
		usedLeftover := need.amount - goalReactant
		adjustedNeedAmount := need.amount - usedLeftover
		times = int64(math.Ceil(float64(adjustedNeedAmount) / float64(formula.result.amount)))
		leftover[formula.result.name] -= usedLeftover
		// We need less of it now
		need.amount = adjustedNeedAmount
		//fmt.Println("Used", usedLeftover, "leftover", formula.result.name, "on", originalNeed)
	} else if leftoverReactant != 0 {
		//fmt.Println("Not using", leftoverReactant, "for", formula.result.name, "because it wouldn't decrease waste")
	}

	for _, reactant := range formula.reactants {
		amountNeededForReaction := reactant.amount * times

		leftover, ore = solveRecursiveImpl(chemical{reactant.name, amountNeededForReaction}, formulas, leftover)
		sumOre += ore
	}

	amountLeftover := (formula.result.amount * times) - need.amount
	leftover[need.name] += amountLeftover

	//fmt.Println("Resolved", need, "with", amountLeftover, "left over")
	//fmt.Println(leftover)
	return leftover, sumOre
}

func solveRecursive(need chemical, formulas map[string]formula) int64 {
	leftover := map[string]int64{}
	sum := int64(0)

	leftover, sum = solveRecursiveImpl(need, formulas, leftover)

	return sum
}

func main() {
	formulas := readInput()

	desiredResult := newChemical("1 FUEL")
	var input int64 = 1
	var output int64 = 0
	const increment int64 = 1000

	fmt.Println("Finding close match to 1 trillion in increments of ", increment, "...")
	for output < 1000000000000 {
		desiredResult.amount = input
		output = solveRecursive(desiredResult, formulas)
		input += increment
	}

	fmt.Println("Found", desiredResult.amount, "=", output, ". Working backwards from there to find answer")

	// We overshot it, so let's go backwards until we find the right
	// answer

	for output > 1000000000000 {
		desiredResult.amount = input
		output = solveRecursive(desiredResult, formulas)
		input--
	}

	fmt.Println("FUEL =", desiredResult.amount, "ORE =", output)
}
