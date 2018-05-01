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

		header := helper.GetSheetValueString(sheet, 0, col)

		// 空列，终止
		if header == "" {
			break
		}

		//if headerRow.Exists(header) {
		//	panic("Duplicate header value")
		//}

		headerRow = append(headerRow, header)
	}

	tab.RawHeader = headerRow
}

func ResolveHeaderFields(tab *model.DataTable, tableObjectType string, symbols *model.SymbolTable) {

	for col, value := range tab.RawHeader {

		tf := symbols.FindField(tableObjectType, value)
		if tf == nil {
			helper.ReportError("HeaderFieldNotDefined", value, helper.Location(tab.FileName, 0, col))
		}

		tab.AddHeaderField(tf)
	}

}
