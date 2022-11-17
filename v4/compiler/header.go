package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"strings"
)

const (
	headerRow_FieldName = 0
	headerRow_FieldType = 1
	headerRow_Meta      = 2
	headerRow_Comment   = 3
	maxHeaderRow        = 4
)

func loadHeaderToCell(sheet util.TableSheet, tab *model.DataTable) (maxCol int) {
	for col := 0; ; col++ {

		for row := 0; row < maxHeaderRow; row++ {
			headerValue := sheet.GetValue(row, col, nil)
			// 空列，终止
			if headerValue == "" {
				return
			}

			maxCol = col

			header := tab.MustGetHeader(col)
			header.Cell.CopyFrom(&model.Cell{
				Value: headerValue,
				Col:   col,
				Row:   row,
				Table: tab,
			})

			tinfo := header.TypeInfo
			tinfo.Usage = model.FieldUsage_Struct
			tinfo.ObjectType = tab.HeaderType

			switch row {
			case headerRow_FieldName:
				tinfo.FieldName = headerValue
			case headerRow_FieldType:
				tinfo.FieldType = headerValue
			case headerRow_Meta:
				parseMeta(header, headerValue)
			case headerRow_Comment:
				tinfo.Comment = headerValue
			}
		}

	}

	return
}

func parseMeta(header *model.HeaderField, meta string) {
	if meta == "" {
		return
	}

	features := strings.Split(meta, ";")

	for _, kvStr := range features {
		if kvStr == "" {
			continue
		}

		kvPair := strings.Split(kvStr, "=")
		var (
			key   string
			value string
		)
		if len(kvPair) == 2 {
			key = kvPair[0]
			value = kvPair[1]
		} else {
			key = kvStr
		}

		switch key {
		case "MakeIndex":
			header.TypeInfo.MakeIndex = true
		case "Spliter":
			header.TypeInfo.ArraySplitter = value
		default:
			util.ReportError("InvalidMetaFormat", key, header.Cell.String())
		}
	}
}

func loadHeader(sheet util.TableSheet, tab *model.DataTable, typeTab *model.TypeTable) (maxCol int) {
	maxCol = loadHeaderToCell(sheet, tab)

	checkHeaderTypes(tab, typeTab)

	return
}

func checkHeaderTypes(tab *model.DataTable, symbols *model.TypeTable) {

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 原始类型检查
		if !util.PrimitiveExists(header.TypeInfo.FieldType) &&
			!symbols.ObjectExists(header.TypeInfo.FieldType) { // 对象检查

			util.ReportError("UnknownFieldType", header.TypeInfo.FieldType, header.Cell.String())
		}
	}

}
