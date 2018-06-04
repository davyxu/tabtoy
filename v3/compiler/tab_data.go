package compiler

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"strings"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int) bool {

	for _, header := range tab.Headers {

		// 取列头所在列和当前行交叉的单元格
		value := helper.GetSheetValueString(sheet, row, header.Cell.Col)

		// 首列带#时，本行忽略
		if header.Cell.Col == 0 && strings.HasPrefix(value, "#") {
			return false
		}

		cell := tab.MustGetCell(row, header.Cell.Col)
		cell.Value = value
	}

	return true
}

func LoadDataTable(filegetter helper.FileGetter, fileName, headerType string) (ret []*model.DataTable, err error) {
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

		loadheader(sheet, tab)

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
