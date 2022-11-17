package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"strings"
)

func loadDataTable(file util.TableFile, fileName string, types *model.TypeTable) (ret []*model.DataTable) {
	for _, sheet := range file.Sheets() {

		tab := model.NewDataTable()
		tab.HeaderType = sheet.Name()
		tab.FileName = fileName

		ret = append(ret, tab)

		maxCol := loadHeader(sheet, tab)
		checkHeaderTypes(tab, types)

		// 遍历所有数据行
		for row := maxHeaderRow; ; row++ {

			if sheet.IsRowEmpty(row, maxCol+1) {
				break
			}

			// 读取每一行
			readOneRow(sheet, tab, row)
		}
	}

	return
}

func readOneRow(sheet util.TableSheet, tab *model.DataTable, row int) bool {

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 浮点数用库取时，需要特殊处理
		isFloat := util.LanguagePrimitive(header.TypeInfo.FieldType, "go") == "float32"

		value := sheet.GetValue(row, header.Col, &util.ValueOption{ValueAsFloat: isFloat})

		// 首列带#时，本行忽略
		if header.Col == 0 && strings.HasPrefix(value, "#") {
			return false
		}

		cell := tab.MustGetCell(row-maxHeaderRow, header.Col)
		cell.Value = value
	}

	return true
}
