package printer

import (
	"encoding/binary"
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

func writeBinary(bf *BinaryFile, ft model.FieldType, value *model.Node) {

	switch ft {
	case model.FieldType_Int32:
		v, _ := strconv.ParseInt(value.Value, 10, 32)
		binary.Write(&bf.buf, binary.LittleEndian, int32(v))
	case model.FieldType_UInt32:
		v, _ := strconv.ParseUint(value.Value, 10, 32)

		binary.Write(&bf.buf, binary.LittleEndian, uint32(v))
	case model.FieldType_Int64:
		v, _ := strconv.ParseInt(value.Value, 10, 64)

		binary.Write(&bf.buf, binary.LittleEndian, int64(v))
	case model.FieldType_UInt64:
		v, _ := strconv.ParseUint(value.Value, 10, 64)

		binary.Write(&bf.buf, binary.LittleEndian, uint64(v))
	case model.FieldType_Float:
		v, _ := strconv.ParseFloat(value.Value, 32)

		binary.Write(&bf.buf, binary.LittleEndian, float32(v))
	case model.FieldType_Bool:
		v, _ := strconv.ParseBool(value.Value)

		binary.Write(&bf.buf, binary.LittleEndian, v)
	case model.FieldType_String:
		bf.WriteString(value.Value)
	case model.FieldType_Enum:
		binary.Write(&bf.buf, binary.LittleEndian, value.EnumValue)
	case model.FieldType_Bytes:
		binary.Write(&bf.buf, binary.LittleEndian, int32(len(value.Raw)))
		binary.Write(&bf.buf, binary.LittleEndian, value.Raw)
	default:
		panic("unsupport type" + model.FieldTypeToString(ft))
	}

}

type CombineBinaryFile struct {
	Datas      []*BinaryFile
	DataByName map[string]*BinaryFile
}

func (self *CombineBinaryFile) Add(f *BinaryFile) bool {

	if _, ok := self.DataByName[f.Name]; ok {
		log.Errorln("duplicate table name in combine binary output:", f.Name)
		return false
	}

	self.Datas = append(self.Datas, f)
	self.DataByName[f.Name] = f

	return true
}

const combineFileVersion = 1

func (self *CombineBinaryFile) Write(outfile string) bool {

	bf := NewBinaryFile("Combine")

	bf.WriteString("TABTOY")
	bf.WriteInt32(combineFileVersion)

	for index, fileBF := range self.Datas {
		bf.WriteInt32(model.MakeTag(model.FieldType_Struct, int32(index)))
		binary.Write(&bf.buf, binary.LittleEndian, fileBF.Data())
	}

	return bf.Write(outfile)
}

func NewCombineBinaryFile() *CombineBinaryFile {
	return &CombineBinaryFile{
		DataByName: make(map[string]*BinaryFile),
	}
}

func PrintBinary(tab *model.Table, rootName string, version string) *BinaryFile {

	bf := NewBinaryFile(rootName)

	bf.WriteInt32(int32(len(tab.Recs)))

	// 遍历每一行
	for _, r := range tab.Recs {

		// 遍历每一列
		for _, node := range r.Nodes {

			// 写入字段索引
			bf.WriteInt32(node.Tag())

			// 写入数量
			if node.IsRepeated {
				bf.WriteInt32(int32(len(node.Child)))
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				for _, valueNode := range node.Child {

					writeBinary(bf, node.Type, valueNode)
				}

			} else {

				// 遍历repeated的结构体
				for _, structNode := range node.Child {

					// 遍历一个结构体的字段
					for _, fieldNode := range structNode.Child {

						// 写入字段索引
						bf.WriteInt32(fieldNode.Tag())

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						writeBinary(bf, fieldNode.Type, valueNode)

					}

				}

			}

		}

	}

	return bf

}
