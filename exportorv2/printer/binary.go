package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

const combineFileVersion = 1

type binaryPrinter struct {
}

func (self *binaryPrinter) Run(g *Globals) *BinaryFile {

	bf := NewBinaryFile()
	bf.WriteString("TABTOY")
	bf.WriteInt32(combineFileVersion)

	for index, tab := range g.Tables {

		if !tab.LocalFD.MatchTag(".bin") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		bf.WriteInt32(model.MakeTag(model.FieldType_Struct, int32(index)))

		if !writeTableBinary(bf, tab) {
			return nil
		}
	}

	return bf
}

func writeTableBinary(self *BinaryFile, tab *model.Table) bool {

	// 表所在的字段

	self.WriteInt32(int32(len(tab.Recs)))

	// 遍历每一行
	for _, r := range tab.Recs {

		// 遍历每一列
		for _, node := range r.Nodes {

			// 写入字段索引
			self.WriteInt32(node.Tag())

			// 写入数量
			if node.IsRepeated {
				self.WriteInt32(int32(len(node.Child)))
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				for _, valueNode := range node.Child {

					self.WriteNodeValue(node.Type, valueNode)
				}

			} else {

				// 遍历repeated的结构体
				for _, structNode := range node.Child {

					// 遍历一个结构体的字段
					for _, fieldNode := range structNode.Child {

						// 写入字段索引
						self.WriteInt32(fieldNode.Tag())

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						self.WriteNodeValue(fieldNode.Type, valueNode)

					}

				}

			}

		}

	}

	return true

}

func init() {

	RegisterPrinter("bin", &binaryPrinter{})

}
