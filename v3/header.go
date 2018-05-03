package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
)

func loadheader(sheet *xlsx.Sheet, tab *model.DataTable) {
	// 读取表头
	var headerRow model.DataRow
	for col := 0; ; col++ {

		headerValue := helper.GetSheetValueString(sheet, 0, col)

		// 空列，终止
		if headerValue == "" {
			break
		}

		headerRow = append(headerRow, model.Cell{
			Value: headerValue,
			Col:   col,
			Row:   0,
			File:  tab.FileName,
			Sheet: sheet.Name,
		})
	}

	tab.RawHeader = headerRow
}

func ResolveHeaderFields(tab *model.DataTable, tableObjectType string, symbols *model.TypeTable) {

	tab.OriginalHeaderType = tableObjectType
	for _, cell := range tab.RawHeader {

		tf := symbols.FieldByName(tableObjectType, cell.Value)
		if tf == nil {
			helper.ReportError("HeaderFieldNotDefined", cell.String())
		}

		tab.AddHeaderField(tf)
	}

}
