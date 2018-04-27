package v3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
)

func readOneRow(sheet *xlsx.Sheet, tab *model.DataTable, row int) (eachRow model.DataRow) {

	for col := 0; col < tab.HeaderFieldCount(); col++ {

		value := util.GetSheetValueString(sheet, row, col)

		eachRow = append(eachRow, value)
	}

	return
}

func loadTable(fileName string) (tab *model.DataTable, err error) {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	// TODO 表名默认来自于文件名，当不使用默认规则时，需要准备Pragma表描述对应关系
	ext := filepath.Ext(fileName)
	tableName := strings.TrimSuffix(fileName, ext)

	for _, sheet := range file.Sheets {

		if tab == nil {
			tab = loadheader(sheet, tableName)
		}

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

	return tab, nil
}
