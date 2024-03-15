package bindata

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
)

// 写入表的一行
func writeStruct(globals *model.Globals, tab *model.DataTable, row int) (*BinaryWriter, error) {

	structWriter := NewBinaryWriter()

	// 一个结构体
	for _, header := range tab.Headers {

		if header == nil {
			continue
		}

		cell := tab.GetCell(row, header.Col)

		if cell == nil {
			continue
		}

		goType := util.LanguagePrimitive(header.TypeInfo.FieldType, "go")

		// 写入字段
		if header.TypeInfo.IsArray() {

			for _, elementValue := range cell.ValueList {

				if err := writePair(globals, structWriter, header.TypeInfo, goType, elementValue, header.Col); err != nil {
					return nil, err
				}
			}

		} else {

			// 空格不输出
			if cell.Value != "" {

				if err := writePair(globals, structWriter, header.TypeInfo, goType, cell.Value, header.Col); err != nil {
					return nil, err
				}
			}

		}

	}

	return structWriter, nil
}
