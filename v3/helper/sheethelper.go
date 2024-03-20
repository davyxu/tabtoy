package helper

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

func WriteIndexTableHeader(sheet TableSheet) {
	sheet.WriteRow("模式", "表类型", "表文件名")
}

func WriteTypeTableHeader(sheet TableSheet) {
	sheet.WriteRow("种类", "对象类型", "标识名", "字段名", "字段类型", "数组切割", "值", "索引", "备注")
}

func WriteRowValues(sheet TableSheet, valueList ...string) {
	sheet.WriteRow(valueList...)
}

func ConvertToCSV(inputFile TableFile) (outputFile TableFile) {

	csvFile := NewCSVFile()

	outSheet := csvFile.Sheets()[0]

	inSheet := inputFile.Sheets()[0]

	// 遍历所有数据行
	for row := 0; ; row++ {

		if inSheet.IsRowEmpty(row, -1) {
			break
		}

		rows := ReadSheetRow(inSheet, row)

		outSheet.WriteRow(rows...)
	}

	return csvFile
}
