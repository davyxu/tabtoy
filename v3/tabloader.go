package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int) (eachRow model.DataRow) {

	for col := 0; col < tab.HeaderFieldCount(); col++ {

		value := helper.GetSheetValueString(sheet, row, col)

		eachRow = append(eachRow, value)
	}

	return
}

func LoadTableData(fileName, headerType string) (ret []*model.DataTable, err error) {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	for _, sheet := range file.Sheets {

		tab := model.NewDataTable()
		tab.HeaderType = headerType
		tab.FileName = fileName

		ret = append(ret, tab)

		loadheader(sheet, tab)

		// 遍历所有行
		for row := 1; ; row++ {

			if helper.IsFullRowEmpty(sheet, row) {
				break
			}

			// 读取每一行
			eachRow := readOneRow(sheet, tab, row)

			tab.AddRow(eachRow)
		}

	}

	return
}
