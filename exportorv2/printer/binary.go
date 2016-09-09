package printer

import (
	"encoding/binary"
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

func writeData(fp *FilePrinter, ft model.FieldType, value *model.Node) {

	switch ft {
	case model.FieldType_Int32:
		v, _ := strconv.ParseInt(value.Value, 10, 32)

		binary.Write(&fp.buf, binary.LittleEndian, int32(v))
	case model.FieldType_UInt32:
		v, _ := strconv.ParseUint(value.Value, 10, 32)

		binary.Write(&fp.buf, binary.LittleEndian, uint32(v))
	case model.FieldType_Int64:
		v, _ := strconv.ParseInt(value.Value, 10, 64)

		binary.Write(&fp.buf, binary.LittleEndian, int64(v))
	case model.FieldType_UInt64:
		v, _ := strconv.ParseUint(value.Value, 10, 64)

		binary.Write(&fp.buf, binary.LittleEndian, uint64(v))
	case model.FieldType_Float:
		v, _ := strconv.ParseFloat(value.Value, 32)

		binary.Write(&fp.buf, binary.LittleEndian, float32(v))
	case model.FieldType_Bool:
		v, _ := strconv.ParseBool(value.Value)

		binary.Write(&fp.buf, binary.LittleEndian, v)
	case model.FieldType_String:

		rawStr := []byte(value.Value)

		binary.Write(&fp.buf, binary.LittleEndian, int32(len(rawStr)))
		binary.Write(&fp.buf, binary.LittleEndian, rawStr)
	case model.FieldType_Enum:
		binary.Write(&fp.buf, binary.LittleEndian, value.EnumValue)
	default:
		panic("unsupport type" + model.FieldTypeToString(ft))
	}

}

func PrintBinary(tab *model.Table, rootName string, version string) *FilePrinter {

	var fp FilePrinter

	// 遍历每一行
	for _, r := range tab.Recs {

		// 遍历每一列
		for _, node := range r.Nodes {

			// 普通值
			if node.Type != model.FieldType_Struct {

				if node.IsRepeated {

					// repeated 值序列
					for _, valueNode := range node.Child {

						writeData(&fp, node.Type, valueNode)
					}
				} else {
					// 单值
					valueNode := node.Child[0]

					writeData(&fp, node.Type, valueNode)
				}

			} else {

				// 遍历repeated的结构体
				for _, structNode := range node.Child {

					// 遍历一个结构体的字段
					for _, fieldNode := range structNode.Child {

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						writeData(&fp, fieldNode.Type, valueNode)

					}

				}

			}

		}

	}

	return &fp

}
