package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
	"strings"
)

func ExportTypes(globals *model.Globals) error {

	for _, oft := range globals.SourceTypes {

		var disableKind string
		if oft.Kind == table.TableKind_None {
			disableKind = "#"
		}

		helper.WriteRowValues(globals.TargetTypesSheet,
			disableKind+oft.Kind.String(),
			oft.ObjectType,
			oft.Name,
			oft.FieldName,
			oft.FieldType,
			oft.ArraySplitter,
			oft.Value)
	}

	return nil
}

func importTypes(globals *model.Globals, sheet *xlsx.Sheet, tabPragma *golexer.KVPair) error {

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

		oft.Value = helper.GetSheetValueString(sheet, row, 3)

		oft.Name = helper.GetSheetValueString(sheet, row, 4)

		// V3无需添加数组前缀

		// 元信息
		meta := helper.GetSheetValueString(sheet, row, 5)

		kvpair := golexer.NewKVPair()
		if err := kvpair.Parse(meta); err != nil {
			continue
		}

		if oft.Value == "" {
			oft.Kind = table.TableKind_None

			globals.AddSourceType(oft)
			continue
		}

		oft.Kind = table.TableKind_Enum

		globals.AddSourceType(oft)
	}

	return nil
}
