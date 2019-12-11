package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	imageWidth  = 25
	imageHeight = 6
)

func main() {
	bytes, _ := ioutil.ReadFile("input")
	imageText := strings.TrimSpace(string(bytes))
	layerSize := imageWidth * imageHeight
	layerCount := len(imageText) / layerSize
	// Array of maps position -> colour
	image := make([][]byte, imageHeight)

	for i := 0; i < imageHeight; i++ {
		image[i] = make([]byte, imageWidth)
	}

	for layerNumber := layerCount - 1; layerNumber >= 0; layerNumber-- {
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				index := (layerSize * layerNumber) + (y * imageWidth) + x
				value := imageText[index]
				if value != '2' {
					image[y][x] = value
				}
			}
		}
	}

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			value := image[y][x]

			if value == '1' {
				fmt.Print("1")
			} else {
				fmt.Print(" ")
			}
		}

		fmt.Println()
	}
}
