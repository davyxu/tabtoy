package compiler

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"strings"
)

func LoadHeader(sheet helper.TableSheet, tab *model.DataTable, resolveTableType string, typeTab *model.TypeTable) {
	// 读取表头

	for col := 0; ; col++ {

		headerValue := sheet.GetValue(0, col, nil)

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

	resolveHeaderFields(tab, resolveTableType, typeTab)

	checkHeaderTypes(tab, typeTab)
}

func checkHeaderTypes(tab *model.DataTable, typeTab *model.TypeTable) {

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 原始类型检查
		if !model.PrimitiveExists(header.TypeInfo.FieldType) &&
			!typeTab.ObjectExists(header.TypeInfo.FieldType) { // 对象检查

			report.ReportError("UnknownFieldType", header.TypeInfo.FieldType, header.Cell.String())
		}
	}

}

func headerValueExists(offset int, name string, headers []*model.HeaderField) bool {

	for i := offset; i < len(headers); i++ {
		if headers[i].Cell.Value == name {
			return true
		}
	}

	return false
}

func resolveHeaderFields(tab *model.DataTable, tableObjectType string, typeTab *model.TypeTable) {

	tab.OriginalHeaderType = tableObjectType
	for index, header := range tab.Headers {

		if header.Cell.Value == "" {
			continue
		}

		tf := typeTab.FieldByName(tableObjectType, header.Cell.Value)
		if tf == nil {
			report.ReportError("HeaderFieldNotDefined", header.Cell.String())
		}

		if headerValueExists(index+1, header.Cell.Value, tab.Headers) && !tf.IsArray() {
			report.ReportError("DuplicateHeaderField", header.Cell.String())
		}

		// 解析好的类型
		header.TypeInfo = tf
	}

}
