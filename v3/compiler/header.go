package compiler

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
	"strings"
)

func CheckHeaderTypes(tab *model.DataTable, types *model.TypeTable) {

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 原始类型检查
		if !table.PrimitiveExists(header.TypeInfo.FieldType) &&
			!types.ObjectExists(header.TypeInfo.FieldType) { // 对象检查

			report.ReportError("UnknownFieldType", header.Cell.String())
		}
	}

}

func loadheader(sheet *xlsx.Sheet, tab *model.DataTable) {
	// 读取表头

	for col := 0; ; col++ {

		headerValue := helper.GetSheetValueString(sheet, 0, col)

		// 空列，终止
		if headerValue == "" {
			break
		}
		// 列头带#时，本列忽略
		if strings.HasPrefix(headerValue, "#") {
			continue
		}

		header := tab.MustGetHeader(col)
		header.Cell.CopyFrom(&model.Cell{
			Value: headerValue,
			Col:   col,
			Row:   0,
			Table: tab,
		})

	}

}

func ResolveHeaderFields(tab *model.DataTable, tableObjectType string, symbols *model.TypeTable) {

	tab.OriginalHeaderType = tableObjectType
	for _, header := range tab.Headers {

		if header.Cell.Value == "" {
			continue
		}

		tf := symbols.FieldByName(tableObjectType, header.Cell.Value)
		if tf == nil {
			report.ReportError("HeaderFieldNotDefined", header.Cell.String())
		}

		// 解析好的类型
		header.TypeInfo = tf
	}

}
