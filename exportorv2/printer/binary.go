package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

const combineFileVersion = 2

type binaryPrinter struct {
}

func (self *binaryPrinter) Run(g *Globals) *Stream {

	fileStresam := NewStream()
	fileStresam.WriteString("TABTOY")
	fileStresam.WriteInt32(combineFileVersion)

	for index, tab := range g.Tables {

		if !tab.LocalFD.MatchTag(".bin") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		if !writeTableBinary(fileStresam, tab, int32(index)) {
			return nil
		}

	}

	return fileStresam
}

func writeTableBinary(tabStream *Stream, tab *model.Table, index int32) bool {

	// 遍历每一行
	for _, r := range tab.Recs {

		rowStream := NewStream()

		// 遍历每一列
		for _, node := range r.Nodes {

			// 子节点数量
			if node.IsRepeated {
				rowStream.WriteInt32(int32(len(node.Child)))
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				for _, valueNode := range node.Child {

					// 写入字段索引
					rowStream.WriteInt32(node.Tag())
					rowStream.WriteNodeValue(node.Type, valueNode)
				}

			} else {

				// 遍历repeated的结构体
				for _, structNode := range node.Child {

					structStream := NewStream()

					// 遍历一个结构体的字段
					for _, fieldNode := range structNode.Child {

						// 写入字段索引
						structStream.WriteInt32(fieldNode.Tag())

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						structStream.WriteNodeValue(fieldNode.Type, valueNode)

					}

					// 真正写到文件中
					rowStream.WriteInt32(node.Tag())
					rowStream.WriteInt32(int32(structStream.Len()))
					rowStream.WriteBytes(structStream.Buffer().Bytes())

				}

			}

		}

		tabStream.WriteInt32(model.MakeTag(int32(model.FieldType_Table), index))
		tabStream.WriteInt32(int32(rowStream.Len()))
		tabStream.WriteBytes(rowStream.Buffer().Bytes())
	}

	return true

}

func init() {

	RegisterPrinter("bin", &binaryPrinter{})

}
