package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
	"strings"
)

func procTypes(globals *model.Globals, sheet *xlsx.Sheet, tabPragma *golexer.KVPair) error {

	pragma := helper.GetSheetValueString(sheet, 0, 0)

	if err := tabPragma.Parse(pragma); err != nil {
		return err
	}

	// 遍历所有行
	for row := 4; ; row++ {

		var oft model.ObjectFieldType

		oft.ObjectType = helper.GetSheetValueString(sheet, row, 0)

		// 空列，终止
		if oft.ObjectType == "" {
			break
		}

		oft.FieldName = helper.GetSheetValueString(sheet, row, 1)

		oft.FieldType = helper.GetSheetValueString(sheet, row, 2)
		if strings.HasPrefix(oft.FieldType, "[]") {
			oft.FieldType = oft.FieldType[2:]
		}

		fieldValue := helper.GetSheetValueString(sheet, row, 3)

		oft.Comment = helper.GetSheetValueString(sheet, row, 4)

		// 默认值
		//defaultValue := helper.GetSheetValueString(sheet, row, 5)

		// V3无需添加数组前缀

		// 元信息
		meta := helper.GetSheetValueString(sheet, row, 6)

		kvpair := golexer.NewKVPair()
		if err := kvpair.Parse(meta); err != nil {
			continue
		}

		if fieldValue == "" {
			oft.IsStruct = true

			globals.SourceTypes = append(globals.SourceTypes, oft)

			//log.Warnf("v3不再支持结构体，忽略结构 %s", oft.ObjectType)
			continue
		}

		oft.IsArray = true

		globals.SourceTypes = append(globals.SourceTypes, oft)

		helper.WriteRowValues(globals.TargetTypesSheet,
			"枚举",
			oft.ObjectType,
			oft.Comment,
			oft.FieldName,
			oft.FieldType, "", "")
	}

	return nil
}
