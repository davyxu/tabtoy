package compiler

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/tealeg/xlsx"
	"strings"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int) bool {

	for _, header := range tab.Headers {

		var value string

		// 浮点数单元格按原样输出
		if model.LanguagePrimitive(header.TypeInfo.FieldType, "go") == "float32" {
			value = helper.GetSheetValueAsNumericString(sheet, row, header.Cell.Col)
		} else {
			// 取列头所在列和当前行交叉的单元格
			value = helper.GetSheetValueString(sheet, row, header.Cell.Col)
		}

		// 首列带#时，本行忽略
		if header.Cell.Col == 0 && strings.HasPrefix(value, "#") {
			return false
		}

		cell := tab.MustGetCell(row, header.Cell.Col)
		cell.Value = value
	}

	return true
}

func LoadDataTable(filegetter helper.FileGetter, fileName, headerType, resolveHeaderType string, typeTab *model.TypeTable) (ret []*model.DataTable, err error) {
	file, err := filegetter.GetFile(fileName)
	if err != nil {
		return nil, err
	}

	for _, sheet := range file.Sheets {

		tab := model.NewDataTable()
		tab.HeaderType = headerType
		tab.FileName = fileName
		tab.SheetName = sheet.Name

		ret = append(ret, tab)

		Loadheader(sheet, tab, resolveHeaderType, typeTab)

		// 遍历所有数据行
		for row := 0; ; row++ {

			if helper.IsFullRowEmpty(sheet, row) {
				break
			}

			// 读取每一行
			readOneRow(sheet, tab, row)
		}

	}

	return
}

func CheckRepeat(inputList *model.DataTableList) {

	for _, tab := range inputList.AllTables() {

		// 遍历输入数据的每一列
		for _, header := range tab.Headers {

			// 输入的列头

			if header.TypeInfo == nil {
				continue
			}

			// 这列需要建立索引
			if header.TypeInfo.MakeIndex {

				checker := map[string]*model.Cell{}

				for row := 1; row < len(tab.Rows); row++ {

					inputCell := tab.GetCell(row, header.Cell.Col)

					// 这行被注释，无效行
					if inputCell == nil {
						break
					}

					if _, ok := checker[inputCell.Value]; ok {

						report.ReportError("DuplicateValueInMakingIndex", inputCell.String())

					} else {
						checker[inputCell.Value] = inputCell
					}

				}
			}
		}
	}
}
