package main

import (
  "fmt"
)

func main() {
  start := 240920
  end := 789857
  count := 0

  for passwordCandidate := start; passwordCandidate <= end; passwordCandidate++ {
    digitNum := 5
    digitDivider := passwordCandidate
    // 10 > any digit
    previousDigit := 10
    repeatCount := 0
    satisfiesRepeatProperty := false
    alwaysIncreases := true

    for ; digitDivider != 0; digitNum-- {
      digit := digitDivider % 10
      digitDivider /= 10

      if previousDigit == digit {
        repeatCount += 1
      } else {
        if repeatCount == 1 {
          satisfiesRepeatProperty = true
        }
        repeatCount = 0
      }

      if digit > previousDigit {
        alwaysIncreases = false
        break
      }
      previousDigit = digit
    }

    if repeatCount == 1 {
      satisfiesRepeatProperty = true
    }

    if satisfiesRepeatProperty && alwaysIncreases {
      fmt.Println(passwordCandidate)
      count += 1
    }
  }
  fmt.Println("Count:", count)
}
