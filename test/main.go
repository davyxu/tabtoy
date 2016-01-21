package main

import (
	"fmt"
	"math"
	"strconv"
)

func str2int(s string) int {
	return int([]byte(s)[0])
}

func mod(a, b int) int {
	return int(math.Mod(float64(a), float64(b)))
}

var asciiA = str2int("A")
var asciiZ = str2int("Z")
var ascii0 = str2int("0")

func num2col(index int) string {

	if index < 1 {
		return "invalid col"
	}

	str := strconv.FormatInt(int64(index), 26)
	fmt.Println(str)

	out := make([]byte, len(str))

	for i, c := range str {

		cIndex, _ := strconv.Atoi(string(c))

		fmt.Println(i, cIndex)

		out[i] = byte(cIndex + asciiA - 1)

	}

	return string(out)
}

func main() {

	fmt.Println("Hello World", num2col(128))
}
