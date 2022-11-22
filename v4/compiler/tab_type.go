package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"strconv"
	"strings"
)

const (
	typeHeaderCol_ObjectType = 0
	typeHeaderCol_FieldName  = 1
	typeHeaderCol_Value      = 2
	typeHeaderCol_Comment    = 3
	maxTypeHeaderCol         = 4
)

func loadTypeHeader(sheet util.TableSheet) (colByHeaderType [maxTypeHeaderCol]int, ok bool) {
	for col := 0; col < maxTypeHeaderCol; col++ {

		headerValue := sheet.GetValue(0, col, nil)

		var headerType int
		switch headerValue {
		case "ObjectType":
			headerType = typeHeaderCol_ObjectType
		case "FieldName":
			headerType = typeHeaderCol_FieldName
		case "Value":
			headerType = typeHeaderCol_Value
		case "Comment":
			headerType = typeHeaderCol_Comment
		default:
			return
		}

		colByHeaderType[headerType] = col
	}

	ok = true

	return
}

func loadTypeTable(file util.TableFile, meta *model.FileMeta, types *model.TypeManager) {
	for _, sheet := range file.Sheets() {

		colByHeaderType, ok := loadTypeHeader(sheet)

		if !ok {
			util.ReportError("InvalidTypeHeader", meta.FileName)
			return
		}

		// 遍历所有数据行
		for row := 1; ; row++ {
			if sheet.IsRowEmpty(row, maxTypeHeaderCol+1) {
				break
			}

			firstCol := sheet.GetValue(row, 0, nil)
			// 首列带#时，本行忽略
			if strings.HasPrefix(firstCol, "#") {
				continue
			}

			var field model.DataField
			field.Usage = model.FieldUsage_Enum
			field.ObjectType = sheet.GetValue(row, colByHeaderType[typeHeaderCol_ObjectType], nil)
			field.FieldType = "int32"
			field.FieldName = sheet.GetValue(row, colByHeaderType[typeHeaderCol_FieldName], nil)
			field.Value = sheet.GetValue(row, colByHeaderType[typeHeaderCol_Value], nil)
			field.Comment = sheet.GetValue(row, colByHeaderType[typeHeaderCol_Comment], nil)

			if field.FieldName == "" {
				cellLocation := cellToString(row, colByHeaderType[typeHeaderCol_FieldName], field.FieldName, meta.FileName, sheet.Name())
				util.ReportError("UnknownFieldName", field.FieldName, cellLocation)
			}

			if types.FieldByName(field.ObjectType, field.FieldName) != nil {
				cellLocation := cellToString(row, colByHeaderType[typeHeaderCol_FieldName], field.FieldName, meta.FileName, sheet.Name())
				util.ReportError("DuplicateTypeField", field.FieldName, cellLocation)
			}

			if field.Value == "" {
				cellLocation := cellToString(row, colByHeaderType[typeHeaderCol_Value], field.FieldName, meta.FileName, sheet.Name())
				util.ReportError("EnumValueEmpty", field.Value, cellLocation)
			}

			_, err := strconv.ParseInt(field.Value, 10, 32)
			if err != nil {
				cellLocation := cellToString(row, colByHeaderType[typeHeaderCol_Value], field.FieldName, meta.FileName, sheet.Name())
				util.ReportError("RequireInteger", field.Value, cellLocation)
			}

			if pre := types.GetEnumByValue(field.ObjectType, field.Value); pre != nil {
				cellLocation := cellToString(row, colByHeaderType[typeHeaderCol_Value], field.FieldName, meta.FileName, sheet.Name())
				util.ReportError("DuplicateEnumValue", field.Value, cellLocation)
			}

			types.AddField(&field, nil, row)
		}

		//只支持导出第一个sheet
		break
	}

	return
}
