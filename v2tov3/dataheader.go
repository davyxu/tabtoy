package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
	"strings"
)

func importDataHeader(globals *model.Globals, sourceSheet, targetSheet *xlsx.Sheet, tableName string) (headerList []model.ObjectFieldType) {

	headerRow := targetSheet.AddRow()

	// 遍历所有行
	for col := 0; ; col++ {

		var oft model.ObjectFieldType
		oft.ObjectType = tableName + "Define"
		oft.Kind = table.TableKind_HeaderStruct

		oft.FieldName = helper.GetSheetValueString(sourceSheet, 0, col)

		// 空列，终止
		if oft.FieldName == "" {
			break
		}

		oft.FieldType = helper.GetSheetValueString(sourceSheet, 1, col)

		// 元信息
		meta := helper.GetSheetValueString(sourceSheet, 2, col)

		oft.Meta = golexer.NewKVPair()
		if err := oft.Meta.Parse(meta); err != nil {
			continue
		}

		if strings.HasPrefix(oft.FieldType, "[]") {
			oft.FieldType = oft.FieldType[2:]
			oft.ArraySplitter = oft.Meta.GetString("ListSpliter")
		}

		oft.Name = helper.GetSheetValueString(sourceSheet, 3, col)

		var disabledForV3 string

		// 添加V3表头
		if globals.TypeIsNoneKind(oft.FieldType) {
			disabledForV3 = "#"
		}

		// 新表的表头加列
		headerRow.AddCell().SetValue(disabledForV3 + oft.Name)

		// 拆分字段填充的数组
		if !globals.SourceTypeExists(oft.ObjectType, oft.FieldName) {

			globals.AddSourceType(oft)
		}

		headerList = append(headerList, oft)
	}

	return
}
