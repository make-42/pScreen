package utils

import (
	"fmt"
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteItem[T any](list []T, index int) []T {
	return append(list[:index], list[index+1:]...)
}

type Coords struct {
	X int
	Y int
}

type CoordsFloat struct {
	X float64
	Y float64
}

// https://stackoverflow.com/questions/64108933/how-to-use-math-pow-with-integers-in-golang
// IntPow calculates n to the mth power. Since the result is an int, it is assumed that m is a positive power
func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func FormatDuration(seconds float64) string {
	hour := int(seconds / 3600)
	minute := int(seconds/60) % 60
	second := int(seconds) % 60
	return fmt.Sprintf("%d:%02d:%02d", hour, minute, second)
}
