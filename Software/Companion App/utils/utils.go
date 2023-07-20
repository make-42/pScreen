package utils

import (
	"fmt"
	"log"
	"math/rand"
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

func WrapString(str string, charactersOnLine int) string {
	newString := ""
	sinceNewLine := 0
	for _, s := range str {
		newString += string(s)
		if string(s) == "\n" {
			sinceNewLine = 0
		}
		if sinceNewLine >= charactersOnLine {
			newString += "\n"
			sinceNewLine = 0
		}
		if string(s) != "\n" && string(s) != "\r" {
			sinceNewLine++
		}
	}
	return newString
}

func RandFloat64Around0() float64 {
	return (rand.Float64() - 0.5) * 2
}

func IsPointInRectangle(rectPos CoordsFloat, rectW, rectH float64, pointPos CoordsFloat, pointRadius float64) bool {
	if pointPos.X > rectPos.X-pointRadius && pointPos.X < rectPos.X+rectW+pointRadius && pointPos.Y > rectPos.Y-pointRadius && pointPos.Y < rectPos.Y+rectH+pointRadius {
		return true
	}
	return false
}
