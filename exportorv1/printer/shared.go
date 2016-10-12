package printer

import (
	"bytes"
	"os"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/exportorv1/data"
)

type IPrinter interface {
	PrintMessage(msg *data.DynamicMessage) bool
	WriteToFile(filename string) bool
}

type printerfile struct {
	printer bytes.Buffer
}

func (self *printerfile) WriteToFile(filename string) bool {

	// 创建输出文件
	file, err := os.Create(filename)
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	// 写入文件头

	file.WriteString(self.printer.String())

	file.Close()

	return true
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

		if i > 0 && valueWrite > 0 {
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
