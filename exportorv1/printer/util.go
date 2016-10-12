package printer

import (
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/exportorv1/data"
	"github.com/davyxu/tabtoy/util"
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
		return util.StringEscape(value)
	}

	return value
}
