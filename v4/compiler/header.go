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

func parseHeader(sheet util.TableSheet, tab *model.DataTable, types *model.TypeManager) (maxCol int) {
	for col := 0; ; col++ {

		fieldName := sheet.GetValue(headerRow_FieldName, col, nil)
		fieldType := sheet.GetValue(headerRow_FieldType, col, nil)

		if fieldName == "" && fieldType == "" {
			break
		}

		if strings.HasPrefix(fieldName, "#") {
			continue
		}

		maxCol = col

		header := model.NewHeaderField(col)
		tab.AddHeader(header)
		types.AddField(header.TypeInfo, tab, 0)
		tinfo := header.TypeInfo
		tinfo.Usage = model.FieldUsage_Struct
		tinfo.ObjectType = tab.HeaderType

		if fieldName == "" {
			util.ReportError("UnknownFieldName", header.String())
		}

		tinfo.FieldName = fieldName
		tinfo.FieldType = fieldType
		fieldMeta := sheet.GetValue(headerRow_Meta, col, nil)
		if errStr := parseMeta(tinfo, fieldMeta); errStr != "" {
			util.ReportError(errStr, fieldMeta, header.String())
		}

		tinfo.Comment = sheet.GetValue(headerRow_Comment, col, nil)
	}

	return
}

func parseMeta(field *model.DataField, meta string) string {
	if meta == "" {
		return ""
	}

	features := strings.Split(meta, " ")

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
		case "ArraySpliter":
			if value == "" {
				return "EmptyArraySpliter"
			}
			field.ArraySplitter = value
		default:
			return "UnknownMetaKey"
		}
	}

	return ""
}

func headerValueExists(offset int, name string, headers []*model.HeaderField) bool {

	for i := offset; i < len(headers); i++ {
		if headers[i].TypeInfo.FieldName == name {
			return true
		}
	}

	return false
}

func checkHeaderTypes(tab *model.DataTable, types *model.TypeManager) {

	for index, header := range tab.Headers {

		// 原始类型检查
		if !util.PrimitiveExists(header.TypeInfo.FieldType) &&
			!types.ObjectExists(header.TypeInfo.FieldType) { // 对象检查

			util.ReportError("UnknownFieldType", header.TypeInfo.FieldType, header.String())
		}

		if headerValueExists(index+1, header.TypeInfo.FieldName, tab.Headers) && !header.TypeInfo.IsArray() {
			util.ReportError("DuplicateHeaderField", header.String())
		}

	}

}
