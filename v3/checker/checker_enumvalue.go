package checker

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
)

// 枚举值的解析是放在输出端处理的, 例如json中, 所以在这里提前检查
func checkEnumValue(globals *model.Globals) {

	for _, tab := range globals.Datas.AllTables() {

		// 遍历输入数据的每一列
		for _, header := range tab.Headers {

			// 输入的列头，为空表示改列被注释
			if header.TypeInfo == nil {
				continue
			}

			if !globals.Types.IsEnumKind(header.TypeInfo.FieldType) {
				continue
			}

			for row := 1; row < len(tab.Rows); row++ {

				inputCell := tab.GetCell(row, header.Cell.Col)

				// 这行被注释，无效行
				if inputCell == nil {
					continue
				}

				if header.TypeInfo.IsArray() {

					for _, v := range inputCell.ValueList {
						checkEnumFieldValue(globals, header, v, inputCell)
					}

				} else {
					checkEnumFieldValue(globals, header, inputCell.Value, inputCell)
				}

			}
		}
	}
}

// 检查枚举值是否存在有效
func checkEnumFieldValue(globals *model.Globals, header *model.HeaderField, value string, inputCell *model.Cell) {

	if inputCell.Value == "" {
		return
	}

	enumValue := globals.Types.GetEnumValue(header.TypeInfo.FieldType, value)
	if enumValue == nil {
		util.ReportError("UnknownEnumValue", header.TypeInfo.FieldType, inputCell.String())
	}

}
