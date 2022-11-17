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

func loadHeader(sheet util.TableSheet, tab *model.DataTable) (maxCol int) {
	for col := 0; ; col++ {

		header := tab.MustGetHeader(col)
		header.Col = col

		for row := 0; row < maxHeaderRow; row++ {
			headerValue := sheet.GetValue(row, col, nil)
			// 空列，终止
			if headerValue == "" {
				return
			}

			maxCol = col

			tinfo := header.TypeInfo
			tinfo.Usage = model.FieldUsage_Struct
			tinfo.ObjectType = tab.HeaderType

			switch row {
			case headerRow_FieldName:
				tinfo.FieldName = headerValue
			case headerRow_FieldType:
				tinfo.FieldType = headerValue
			case headerRow_Meta:
				if !parseMeta(header.TypeInfo, headerValue) {
					util.ReportError("InvalidMetaFormat", headerValue, header.Location())
				}
			case headerRow_Comment:
				tinfo.Comment = headerValue
			}
		}

	}
}

func parseMeta(field *model.DataField, meta string) bool {
	if meta == "" {
		return true
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
			field.MakeIndex = true
		case "Spliter":
			field.ArraySplitter = value
		default:
			return false
		}
	}

	return true
}

func checkHeaderTypes(tab *model.DataTable, types *model.TypeTable) {

	for _, header := range tab.Headers {

		// 原始类型检查
		if !util.PrimitiveExists(header.TypeInfo.FieldType) &&
			!types.ObjectExists(header.TypeInfo.FieldType) { // 对象检查

			util.ReportError("UnknownFieldType", header.TypeInfo.FieldType, header.Location())
		}
	}

}
