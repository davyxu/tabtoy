package data

import (
	"bytes"
	"fmt"
	"math"
)

func mod(a, b int) int {
	return int(math.Mod(float64(a), float64(b)))
}

func str2int(s string) int {
	return int([]byte(s)[0])
}

var asciiA = str2int("A")

const unit = 26

// 按excel格式 R1C1格式转A1
func index2Alphabet(number int) string {

	if number < 1 {
		return ""
	}

	n := number

	nl := make([]int, 0)

	for {

		quo := n / unit

		var reminder int
		x := mod(n, unit)

		// 余数为0时, 要跳过这个0, 重新计算除数(影响进位)
		if x == 0 {
			reminder = unit
			n--
			quo = n / unit
		} else {
			reminder = x
		}

		nl = append(nl, reminder)

		if quo == 0 {
			break
		}

		n = quo
	}

	var out bytes.Buffer

	for i := len(nl) - 1; i >= 0; i-- {

		v := nl[i]

		out.WriteString(string(v + asciiA - 1))
	}

	return out.String()
}

func ConvR1C1toA1(r, c int) string {
	return fmt.Sprintf("%s%d", index2Alphabet(c), r)
}
