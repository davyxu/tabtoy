package printer

import (
	"bytes"
	"fmt"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
)

type IWriter interface {
	PrintMessage(msg *data.DynamicMessage)
}

type dataWriter interface {
	RepeatedMessageBegin(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int)

	WriteMessage(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, indent int)

	RepeatedMessageEnd(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int)

	RepeatedValueBegin(fd *pbmeta.FieldDescriptor)

	WriteValue(fd *pbmeta.FieldDescriptor, value string, indent int)

	RepeatedValueEnd(fd *pbmeta.FieldDescriptor)

	WriteFieldSpliter()
}

func rawWriteMessage(printer *bytes.Buffer, writer dataWriter, msg *data.DynamicMessage, indent int) (valueWrite int) {

	for i := 0; i < msg.Desc.FieldCount(); i++ {
		fd := msg.Desc.Field(i)

		var needSpliter bool

		pos := printer.Len()

		if i > 0 {
			writer.WriteFieldSpliter()
		}

		if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {

			if fd.IsRepeated() {

				// 多结构
				if msgList := msg.GetRepeatedMessage(fd); msgList != nil {

					writer.RepeatedMessageBegin(fd, msg, len(msgList), indent+1)

					for msgIndex, msg := range msgList {

						beginPos := printer.Len()
						writer.WriteMessage(fd, msg, indent+1)
						endPos := printer.Len()

						// 是否有数据写入
						if endPos > beginPos {
							valueWrite++

							// 最后一个不要加
							if msgIndex < len(msgList)-1 {
								writer.WriteFieldSpliter()
							}
						}

						if indent == 0 {
							printer.WriteString("\n")
						}

					}

					writer.RepeatedMessageEnd(fd, msg, len(msgList), indent+1)

					needSpliter = true

					continue
				}

			} else {
				// 单结构
				if msg := msg.GetMessage(fd); msg != nil {

					beginPos := printer.Len()
					writer.WriteMessage(fd, msg, indent+1)
					endPos := printer.Len()

					// 是否有数据写入
					if endPos > beginPos {
						valueWrite++

						needSpliter = true
					}

				}
			}

		} else {

			// 多值
			if fd.IsRepeated() {

				if valuelist := msg.GetRepeatedValue(fd); valuelist != nil {

					writer.RepeatedValueBegin(fd)

					for vIndex, value := range valuelist {

						writer.WriteValue(fd, value, indent)
						valueWrite++

						if vIndex < len(valuelist)-1 {
							writer.WriteFieldSpliter()

						}
					}

					writer.RepeatedValueEnd(fd)

					needSpliter = true
				}

			} else {

				// 单值
				if value, ok := msg.GetValue(fd); ok {
					writer.WriteValue(fd, value, indent)
					valueWrite++

					needSpliter = true
				}

			}

		}

		//没有值输出, 把打分割之前的位置恢复
		if !needSpliter {
			printer.Truncate(pos)
		}

	}

	return

}

func getDefaultValueCount(msg *data.DynamicMessage) int {

	var valueCount int

	for i := 0; i < msg.Desc.FieldCount(); i++ {
		fd := msg.Desc.Field(i)
		if fd.DefaultValue() != "" {
			valueCount++
		}
	}

	return valueCount
}

func valueWrapper(fd *pbmeta.FieldDescriptor, value string) string {
	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_STRING:
		return strEscape(value)
	}

	return value
}

func strEscape(s string) string {

	b := make([]byte, 0)

	var index int

	// 表中直接使用换行会干扰最终合并文件格式, 所以转成\n,由pbt文本解析层转回去
	for index < len(s) {
		c := s[index]

		switch c {
		case '"':
			b = append(b, '\\')
			b = append(b, '"')
		case '\x0A':
			b = append(b, '\\')
			b = append(b, 'n')
		case '\x0D':
			b = append(b, '\\')
			b = append(b, 'r')
		default:
			b = append(b, c)
		}

		index++

	}

	return fmt.Sprintf("\"%s\"", string(b))

}
