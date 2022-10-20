package helpers

import (
	"math"
)

func IntToCode(num int, charArr []string) string {
	arrLen := len(charArr)
	var strOut string

	for num > arrLen-1 {
		idx := (num % int(arrLen))
		strOut = charArr[idx] + strOut
		num = int(math.Floor(float64(num / arrLen)))
	}

	return charArr[num%arrLen] + strOut
}

func CodeToInt(code string, charArr []string) int {
	retVal := 0

	for i := 0; i < len(code); i++ {
		x := float64(len(charArr))
		y := float64(len(code) - i - 1)
		retVal += int(math.Pow(x, y)) * indexOf(string(code[i]), charArr)
	}
	return retVal
}

func indexOf(ch string, charArr []string) int {
	for i := 0; i < len(charArr); i++ {
		if ch == charArr[i] {
			return i
		}
	}

	return -1
}
