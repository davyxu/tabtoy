package util

import (
	"github.com/tealeg/xlsx"
	"strings"
)

func GetSheetValueString(sheet *xlsx.Sheet, row, col int) string {
	c := sheet.Cell(row, col)

	return strings.TrimSpace(c.Value)
}

// 整行都是空的
func IsFullRowEmpty(sheet *xlsx.Sheet, row int) bool {

	for col := 0; col < sheet.MaxCol; col++ {

		data := GetSheetValueString(sheet, row, col)

		if data != "" {
			return false
		}
	}

	return true
}
