package bindata

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func writeHeader(writer *BinaryWriter) error {
	if err := writer.WriteString("TABTOY"); err != nil {
		return err
	}
	if err := writer.WriteUInt32(4); err != nil {
		return err
	}

	return nil
}

func Generate(globals *model.Globals) (data []byte, err error) {

	totalWriter := NewBinaryWriter()

	if err := writeHeader(totalWriter); err != nil {
		return nil, err
	}

	for _, tab := range globals.Datas.AllTables() {

		// 结构体的标记头, 方便跨过不同类型
		structTag := MakeTagStructArray()
		if err := totalWriter.WriteUInt32(structTag); err != nil {
			return nil, err
		}

		totalWriter.WriteString(tab.HeaderType)

		totalDataRow := len(tab.Rows) - 1
		totalWriter.WriteUInt32(uint32(totalDataRow))

		// 表的每一个行
		for row := 1; row < len(tab.Rows); row++ {

			if swriter, err := writeStruct(globals, tab, row); err != nil {
				return nil, err
			} else {
				structData := swriter.Bytes()
				// 结构体二进制边界
				totalWriter.WriteUInt32(uint32(len(structData)))
				totalWriter.Write(structData)
			}
		}
	}

	data = totalWriter.Bytes()

	return
}
