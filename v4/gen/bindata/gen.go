package bindata

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
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

func exportTable(globals *model.Globals, writer *BinaryWriter, tab *model.DataTable) error {
	// 结构体的标记头, 方便跨过不同类型
	if err := writer.WriteUInt32(MakeTagStructArray()); err != nil {
		return err
	}

	writer.WriteString(tab.HeaderType)

	totalDataRow := len(tab.Rows)
	writer.WriteUInt32(uint32(totalDataRow))

	// 表的每一个行
	for row := 0; row < totalDataRow; row++ {

		if swriter, err := writeStruct(globals, tab, row); err != nil {
			return err
		} else {
			structData := swriter.Bytes()
			// 结构体二进制边界
			writer.WriteUInt32(uint32(len(structData)))
			writer.Write(structData)
		}
	}

	return nil
}

func OutputFile(globals *model.Globals, outFile string) (err error) {

	totalWriter := NewBinaryWriter()

	if err = writeHeader(totalWriter); err != nil {
		return err
	}

	for _, tab := range globals.Datas.AllTables() {

		err = exportTable(globals, totalWriter, tab)
		if err != nil {
			return err
		}
	}

	data := totalWriter.Bytes()

	err = util.WriteFile(outFile, data)

	return
}

func OutputDir(globals *model.Globals, outputDir string) (err error) {

	for _, tab := range globals.Datas.AllTables() {

		writer := NewBinaryWriter()
		if err = writeHeader(writer); err != nil {
			return err
		}

		err = exportTable(globals, writer, tab)
		if err != nil {
			return err
		}

		err = util.WriteFile(fmt.Sprintf("%s/%s.bin", outputDir, tab.HeaderType), writer.Bytes())

		if err != nil {
			return err
		}
	}

	return nil
}
