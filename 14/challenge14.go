package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func min(a, b int) int {
	if a < b {
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
	amount int
}

func newChemical(chemicalString string) chemical {
	chemicalSplit := strings.Split(strings.TrimSpace(chemicalString), " ")
	if len(chemicalSplit) != 2 {
		panic("Incorrect chemical split")
	}

	amountStr, name := chemicalSplit[0], chemicalSplit[1]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println(err)
		panic("Error while atoiing chemical name")
	}

	return chemical{name, amount}
}

func readInput() map[string]formula {
	// maps results to reactants
	formulas := map[string]formula{}

	inputBytes, err := ioutil.ReadFile("bigsample1")

	if err != nil {
		fmt.Println(err)
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
	fmt.Println(formulas)
	return formulas
}

func solveRecursiveImpl(need chemical, formulas map[string]formula, leftover map[string]int) (map[string]int, int) {
	if need.name == "ORE" {
		return leftover, need.amount
	}

	formula := formulas[need.name]

	sumOre := 0
	ore := 0
	times := int(math.Ceil(float64(need.amount) / float64(formula.result.amount)))

	fmt.Println("Satisfying requirement", need, "with formula", formula, "(", times, "times)")

	for _, reactant := range formula.reactants {
		amountNeededForReaction := reactant.amount * times

		// if the amount needed is less than the amount left over, just take the amount needed
		leftoverReactant := leftover[reactant.name]
		if leftoverReactant >= reactant.amount {
			usedReactant := reactant.amount * (leftoverReactant / reactant.amount)
			amountNeededForReaction -= usedReactant
			leftover[reactant.name] -= usedReactant
			fmt.Println("[", formula, "] Used", usedReactant, "leftover", reactant.name, "on", reactant)
		} else if leftoverReactant != 0 {
			fmt.Println("[", formula, "] Not using", leftoverReactant, "for", reactant, "because it wouldn't decrease waste")
		}

		if amountNeededForReaction > 0 {
			leftover, ore = solveRecursiveImpl(chemical{reactant.name, amountNeededForReaction}, formulas, leftover)
			sumOre += ore
		}
	}

	amountLeftover := (formula.result.amount * times) - need.amount
	leftover[need.name] += amountLeftover

	fmt.Println("Resolved", need, "with", amountLeftover, "left over")
	fmt.Println(leftover)
	return leftover, sumOre
}

func solveRecursive(need chemical, formulas map[string]formula) int {
	leftover := map[string]int{}
	sum := 0

	leftover, sum = solveRecursiveImpl(need, formulas, leftover)

	return sum
}

func solve(need chemical, formulas map[string]formula) int {
	requirements := map[string]int{need.name: need.amount}
	numResolved := 1
	// leftover := map[string]int{}

	for len(requirements) != 1 || requirements["ORE"] == 0 {
		for i := 0; numResolved != 0; i++ {
			numResolved = 0
			// Requirements to be added this time
			newRequirements := map[string]int{}

			// Resolve perfect requirements
			for name, amountRequired := range requirements {
				if name == "ORE" {
					newRequirements["ORE"] += amountRequired
					continue
				}

				formula := formulas[name]
				fmt.Println("Attempting to resolve requirement", name, amountRequired, "with formula", formula)
				times := int(math.Ceil(float64(amountRequired) / float64(formula.result.amount)))
				if amountRequired != times*formula.result.amount {
					fmt.Println("Not resolving imperfect formula")
					newRequirements[name] += amountRequired
				} else {
					fmt.Println("Perfect match!")
					for _, reactant := range formula.reactants {
						newRequirements[reactant.name] += times * reactant.amount
						fmt.Println("Added new requirement", reactant.name, newRequirements[reactant.name])
					}
					numResolved++
				}
			}

			// Resolve imperfect requirements
			for name, amountRequired := range requirements {
				if name == "ORE" {
					newRequirements["ORE"] += amountRequired
					continue
				}

				formula := formulas[name]
				fmt.Println("Attempting to resolve requirement", name, amountRequired, "with formula", formula)
				times := int(math.Ceil(float64(amountRequired) / float64(formula.result.amount)))
				if amountRequired != times*formula.result.amount {
					fmt.Println("Not resolving imperfect formula")
					newRequirements[name] += amountRequired
				} else {
					fmt.Println("Not resolving perfect formula")
				}
			}

			requirements = newRequirements
			fmt.Println("Resolved", numResolved)
			fmt.Println(requirements)
		}
	}

	fmt.Println(requirements)
	return requirements["ORE"]
}

func main() {
	formulas := readInput()

	desiredResult := newChemical("1 FUEL")

	fmt.Println(solveRecursive(desiredResult, formulas))

	fmt.Println()
}
