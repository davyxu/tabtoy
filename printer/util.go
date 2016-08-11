package printer

import (
	"fmt"

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
