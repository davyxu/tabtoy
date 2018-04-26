package v3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
)

func loadTable(globals *model.Globals, fileName string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	// TODO 表名
	ext := filepath.Ext(fileName)
	tableName := strings.TrimSuffix(fileName, ext)

	tab := model.NewDataTable(tableName)

	for _, sheet := range file.Sheets {

		// 读取表头
		for col := 0; ; col++ {

			header := util.GetSheetValueString(sheet, 0, col)

			// 空列，终止
			if header == "" {
				break
			}

			t := globals.Symbols.QueryType(tableName, header)

			if t == nil {
				panic("types not found:" + header)
			}

			tab.AddHeader(t)
		}

		for row := 1; ; row++ {

			if util.IsFullRowEmpty(sheet, row) {
				break
			}

			for col := 0; col < tab.MaxColumns(); col++ {

				value := util.GetSheetValueString(sheet, row, col)

				if value == "" {
					break
				}

				tab.AddRow(row-1, col, value)

			}

		}

	}

	globals.AddData(tab)

	return nil
}
