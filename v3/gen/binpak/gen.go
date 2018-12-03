package binpak

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func writeHeader(writer *BinaryWriter) error {
	if err := writer.WriteString("TABTOY"); err != nil {
		return err
	}
	if err := writer.WriteUInt32(3); err != nil {
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

		totalDataRow := len(tab.Rows) - 1
		totalWriter.WriteUInt32(uint32(totalDataRow))

		// 结构体数组
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
