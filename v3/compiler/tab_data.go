package compiler

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"strings"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int, eachRow *model.DataRow) bool {

	for _, headerCell := range tab.RawHeader {

		// 取列头所在列和当前行交叉的单元格
		value := helper.GetSheetValueString(sheet, row, headerCell.Col)

		// 首列带#时，本行忽略
		if headerCell.Col == 0 && strings.HasPrefix(value, "#") {
			return false
		}

		*eachRow = append(*eachRow, model.Cell{
			Value: value,
			Row:   row,
			Col:   headerCell.Col,
			File:  tab.FileName,
			Sheet: tab.SheetName,
		})
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

		// 遍历所有行
		for row := 1; ; row++ {

			if helper.IsFullRowEmpty(sheet, row) {
				break
			}

			// 读取每一行
			var eachRow model.DataRow
			if readOneRow(sheet, tab, row, &eachRow) {
				tab.AddRow(eachRow)
			}
		}

	}

	return
}
