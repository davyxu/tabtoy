package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
	"strings"
)

func procDataHeader(globals *model.Globals, sourceSheet, targetSheet *xlsx.Sheet, tableName string) (headerList []model.ObjectFieldType) {

	headerRow := targetSheet.AddRow()

	// 遍历所有行
	for col := 0; ; col++ {

		var oft model.ObjectFieldType
		oft.ObjectType = tableName + "Define"

		oft.FieldName = helper.GetSheetValueString(sourceSheet, 0, col)

		// 空列，终止
		if oft.FieldName == "" {
			break
		}

		oft.FieldType = helper.GetSheetValueString(sourceSheet, 1, col)

		if strings.HasPrefix(oft.FieldType, "[]") {
			oft.FieldType = oft.FieldType[2:]
			oft.IsArray = true
		}

		// 元信息
		meta := helper.GetSheetValueString(sourceSheet, 2, col)

		oft.Meta = golexer.NewKVPair()
		if err := oft.Meta.Parse(meta); err != nil {
			continue
		}

		oft.Comment = helper.GetSheetValueString(sourceSheet, 3, col)

		var disabledForV3 string

		// 添加V3表头
		if globals.TypeIsStruct(oft.FieldType) {
			disabledForV3 = "#"
		}

		headerRow.AddCell().SetValue(disabledForV3 + oft.Comment)

		helper.WriteRowValues(globals.TargetTypesSheet,
			disabledForV3+"表头",
			oft.ObjectType,
			oft.Comment,
			oft.FieldName,
			oft.FieldType,
			oft.Meta.GetString("ListSpliter"),
			"")

		globals.SourceTypes = append(globals.SourceTypes, oft)

		headerList = append(headerList, oft)
	}

	return
}
