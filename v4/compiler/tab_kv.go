package compiler

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"path/filepath"
	"strings"
)

const (
	kvHeaderCol_Key     = 0
	kvHeaderCol_Type    = 1
	kvHeaderCol_Value   = 2
	kvHeaderCol_Comment = 3
	kvHeaderCol_Meta    = 4
	maxKVHeaderCol      = 5
)

func loadKVHeader(sheet util.TableSheet) (colByHeaderType [maxKVHeaderCol]int, ok bool) {
	for col := 0; col < maxKVHeaderCol; col++ {

		headerValue := sheet.GetValue(0, col, nil)

		var headerType int
		switch headerValue {
		case "Key":
			headerType = kvHeaderCol_Key
		case "Type":
			headerType = kvHeaderCol_Type
		case "Value":
			headerType = kvHeaderCol_Value
		case "Comment":
			headerType = kvHeaderCol_Comment
		case "Meta":
			headerType = kvHeaderCol_Meta
		default:
			return
		}

		colByHeaderType[headerType] = col
	}

	ok = true

	return
}

func loadKVTable(file util.TableFile, fileName string, types *model.TypeManager) (ret []*model.DataTable) {
	for _, sheet := range file.Sheets() {

		colByHeaderType, ok := loadKVHeader(sheet)

		if !ok {
			util.ReportError("InvalidKVHeader", fileName)
			return
		}

		tab := model.NewDataTable()
		if sheet.Name() == "" {
			tab.HeaderType = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		} else {
			tab.HeaderType = sheet.Name()
		}
		tab.FileName = fileName

		ret = append(ret, tab)

		// 添加输出数据行, 只有一行
		tab.AddRow()

		// 遍历所有数据行
		for row := 1; ; row++ {
			if sheet.IsRowEmpty(row, maxKVHeaderCol+1) {
				break
			}

			firstCol := sheet.GetValue(row, 0, nil)
			// 首列带#时，本行忽略
			if strings.HasPrefix(firstCol, "#") {
				continue
			}

			header := model.NewHeaderField(row - 1)
			tab.AddHeader(header)

			field := header.TypeInfo
			field.Usage = model.FieldUsage_Struct
			field.ObjectType = tab.HeaderType
			field.FieldName = sheet.GetValue(row, colByHeaderType[kvHeaderCol_Key], nil)
			field.FieldType = sheet.GetValue(row, colByHeaderType[kvHeaderCol_Type], nil)

			// 原始类型检查
			if !util.PrimitiveExists(field.FieldType) && !types.ObjectExists(field.FieldType) { // 对象检查
				cellLocation := kvCellToString(row, colByHeaderType[kvHeaderCol_Type], field.FieldType, fileName, sheet.Name())
				util.ReportError("UnknownFieldType", field.FieldType, cellLocation)
			}
			field.Comment = sheet.GetValue(row, colByHeaderType[kvHeaderCol_Comment], nil)
			fieldMeta := sheet.GetValue(row, colByHeaderType[kvHeaderCol_Meta], nil)
			if errStr := parseMeta(field, fieldMeta); errStr != "" {
				cellLocation := kvCellToString(row, colByHeaderType[kvHeaderCol_Meta], fieldMeta, fileName, sheet.Name())
				util.ReportError(errStr, fieldMeta, cellLocation)
			}

			if field.FieldName == "" {
				cellLocation := kvCellToString(row, colByHeaderType[kvHeaderCol_Key], field.FieldName, fileName, sheet.Name())
				util.ReportError("UnknownFieldName", cellLocation)
			}

			if types.FieldByName(field.ObjectType, field.FieldName) != nil {
				cellLocation := kvCellToString(row, colByHeaderType[kvHeaderCol_Key], field.FieldName, fileName, sheet.Name())
				util.ReportError("DuplicateKVField", cellLocation)
			}

			types.AddField(field, tab, row)

			// 数据
			isFloat := util.LanguagePrimitive(field.FieldType, "go") == "float32"
			fieldValue := sheet.GetValue(row, colByHeaderType[kvHeaderCol_Value], &util.ValueOption{ValueAsFloat: isFloat})
			cell := tab.AddCell(0)
			cell.Value = fieldValue

			if field.IsArray() {
				for _, element := range strings.Split(fieldValue, field.ArraySplitter) {
					cell.ValueList = append(cell.ValueList, element)
				}
			}
		}
	}

	return
}

func kvCellToString(row, col int, value, file, sheet string) string {
	return fmt.Sprintf("'%s' @%s|%s(%s)", value, file, sheet, util.R1C1ToA1(row+1, col+1))
}
