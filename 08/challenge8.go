package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

const (
	imageWidth  = 25
	imageHeight = 6
)

func main() {
	bytes, _ := ioutil.ReadFile("input")
	image := strings.TrimSpace(string(bytes))
	layerSize := imageWidth * imageHeight
	layerCount := len(image) / layerSize
	digitCounts := make([]map[rune]int, layerCount)

	for i := 0; i < layerCount; i++ {
		digitCounts[i] = map[rune]int{}
	}

	fmt.Println(len(image))

	for index, character := range image {
		layerNumber := index / layerSize
		digitCounts[layerNumber][character]++
	}

	minZeroes := math.MaxInt32
	minZeroesLayer := -1

	for layerNumber, charMap := range digitCounts {
		if charMap['0'] < minZeroes {
			minZeroesLayer = layerNumber
			minZeroes = charMap['0']
		}
	}

	layer := digitCounts[minZeroesLayer]
	fmt.Println(layer['1'] * layer['2'])
}
