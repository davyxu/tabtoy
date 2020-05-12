package binpak

import (
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
)

// 写入表的一行
func writeStruct(globals *model.Globals, tab *model.DataTable, row int) (*BinaryWriter, error) {

	structWriter := NewBinaryWriter()

	// 一个结构体
	for _, header := range tab.Headers {

		if header == nil {
			continue
		}

		cell := tab.GetCell(row, header.Cell.Col)

		if cell == nil {
			continue
		}

		goType := model.LanguagePrimitive(header.TypeInfo.FieldType, "go")

		// 写入字段
		if header.TypeInfo.IsArray() {

			for _, elementValue := range strings.Split(cell.Value, header.TypeInfo.ArraySplitter) {

				if err := writePair(globals, structWriter, header.TypeInfo, goType, elementValue, header.Cell.Col); err != nil {
					return nil, err
				}
			}

		} else {

			// 空格不输出
			if cell.Value != "" {

				if err := writePair(globals, structWriter, header.TypeInfo, goType, cell.Value, header.Cell.Col); err != nil {
					return nil, err
				}
			}

		}

	}

	return structWriter, nil
}
