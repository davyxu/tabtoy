package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/tealeg/xlsx"
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

func testOutput() {
	fileOut := xlsx.NewFile()
	outSheet := fileOut.AddSheet("sync")

	file, err := xlsx.OpenFile("Actor.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	sheet := file.Sheets[0]

	for y := 0; y <= 2; y++ {

		outRow := outSheet.AddRow()

		row := sheet.Rows[y].Cells

		for x := 0; x < len(row); x++ {
			cell := outRow.AddCell()
			cell.Value = row[x].Value
		}
	}

	fileOut.Save("copy.xlsx")

}

func main() {

}
