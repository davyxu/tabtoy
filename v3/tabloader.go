package v3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int) (eachRow model.DataRow) {

	for col := 0; col < tab.HeaderFieldCount(); col++ {

		value := util.GetSheetValueString(sheet, row, col)

		eachRow = append(eachRow, value)
	}

	return
}

func LoadTableData(fileName string, tab *model.DataTable) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	for _, sheet := range file.Sheets {

		loadheader(sheet, tab)

		// 遍历所有行
		for row := 1; ; row++ {

			if util.IsFullRowEmpty(sheet, row) {
				break
			}

			// 读取每一行
			eachRow := readOneRow(sheet, tab, row)

			tab.AddRow(eachRow)
		}

	}

	return nil
}
