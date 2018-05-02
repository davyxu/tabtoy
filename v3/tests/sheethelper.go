package tests

import "github.com/tealeg/xlsx"

func WriteRowValues(sheet *xlsx.Sheet, valueList ...interface{}) {

	row := sheet.AddRow()

	if sheet.MaxCol != 0 && sheet.MaxCol != len(valueList) {
		panic("diff col count")
	}

	for _, value := range valueList {

		cell := row.AddCell()
		cell.SetValue(value)
	}

}

func createSheet() *xlsx.Sheet {
	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("Sheet1")

	return sheet
}

func WriteTypeTableHeader(sheet *xlsx.Sheet) {

	WriteRowValues(sheet, "种类", "对象类型", "标识名", "字段名", "字段类型", "数组切割", "值")
}

func WriteIndexTableHeader(sheet *xlsx.Sheet) {

	WriteRowValues(sheet, "模式", "表类型", "表文件名")
}
