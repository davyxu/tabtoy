package helper

import (
	"github.com/tealeg/xlsx"
	"strings"
)

func GetSheetValueString(sheet *xlsx.Sheet, row, col int) string {
	c := sheet.Cell(row, col)

	return strings.TrimSpace(c.Value)
}

func GetSheetValueAsNumericString(sheet *xlsx.Sheet, row, col int) string {
	c := sheet.Cell(row, col)

	str, _ := c.GeneralNumeric()

	return strings.TrimSpace(str)
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

func WriteRowValues(sheet *xlsx.Sheet, valueList ...interface{}) {

	row := sheet.AddRow()

	//if sheet.MaxCol != 0 && sheet.MaxCol != len(valueList) {
	//	panic("diff col count")
	//}

	for _, value := range valueList {

		cell := row.AddCell()
		cell.SetValue(value)
	}

}

func WriteTypeTableHeader(sheet *xlsx.Sheet) {

	WriteRowValues(sheet, "种类", "对象类型", "标识名", "字段名", "字段类型", "数组切割", "值", "索引")
}

func WriteIndexTableHeader(sheet *xlsx.Sheet) {

	WriteRowValues(sheet, "模式", "表类型", "表文件名")
}
