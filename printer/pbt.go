package printer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
)

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

func makeMsg(msg *data.DynamicMessage, indent int) string {
	msgStr := strings.TrimSpace(writeString(msg, indent))

	if msgStr == "" {

		// 当一个消息有默认值, 且不是最顶层消息时, 输出一个空, 产生空消息使用默认值
		if indent > 1 && getDefaultValueCount(msg) > 0 {
			return "{}"
		}

		return ""
	}
	return fmt.Sprintf("{%s}", msgStr)
}

func makeValue(fd *pbmeta.FieldDescriptor, value string) string {
	return fmt.Sprintf("%s: %s", fd.Name(), valueWrapper(fd, value))
}

func valueWrapper(fd *pbmeta.FieldDescriptor, value string) string {
	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_STRING:
		return fmt.Sprintf("\"%s\"", value)
	}

	return value
}

func writeString(msg *data.DynamicMessage, indent int) string {
	var out bytes.Buffer

	for i := 0; i < msg.Desc.FieldCount(); i++ {
		fd := msg.Desc.Field(i)

		if i > 0 {
			out.WriteString(" ")
		}

		if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {

			if fd.IsRepeated() {

				// 多结构
				if msgList := msg.GetRepeatedMessage(fd); msgList != nil {

					for _, msg := range msgList {

						msgStr := makeMsg(msg, indent+1)

						if msgStr != "" {

							out.WriteString(fmt.Sprintf("%s %s", fd.Name(), msgStr))

							if indent == 0 {
								out.WriteString("\n")
							} else {
								out.WriteString(" ")
							}
						}

					}

					continue
				}

			} else {
				// 单结构
				if msg := msg.GetMessage(fd); msg != nil {
					msgStr := makeMsg(msg, indent+1)
					if msgStr != "" {
						out.WriteString(fmt.Sprintf("%s %s", fd.Name(), msgStr))
					}

					continue
				}
			}

		} else {

			// 多值
			if fd.IsRepeated() {

				if valuelist := msg.GetRepeatedValue(fd); valuelist != nil {
					for vIndex, value := range valuelist {

						if vIndex > 0 {
							out.WriteString(" ")
						}

						out.WriteString(makeValue(fd, value))

					}

					continue
				}

			} else {

				// 单值
				if value, ok := msg.GetValue(fd); ok {
					out.WriteString(makeValue(fd, value))
					continue
				}

			}

		}

	}

	return out.String()

}

func WriteProtoBufferText(msg *data.DynamicMessage) string {
	return writeString(msg, 0)
}
