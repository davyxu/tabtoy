package data

import (
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
)

func GetDefaultValue(fd *pbmeta.FieldDescriptor) string {
	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_FLOAT,
		pbprotos.FieldDescriptorProto_TYPE_INT64,
		pbprotos.FieldDescriptorProto_TYPE_UINT64,
		pbprotos.FieldDescriptorProto_TYPE_INT32,
		pbprotos.FieldDescriptorProto_TYPE_UINT32:
		if fd.DefaultValue() != "" && fd.DefaultValue() != "0" {
			break
		}

		return "0"

	case pbprotos.FieldDescriptorProto_TYPE_BOOL:

		if fd.DefaultValue() != "" && fd.DefaultValue() != "false" {
			break
		}

		return "false"

	case pbprotos.FieldDescriptorProto_TYPE_STRING:

		if fd.DefaultValue() != "" {
			break
		}

		return ""

	case pbprotos.FieldDescriptorProto_TYPE_ENUM:

		ed := fd.EnumDesc()

		// 好奇葩
		if ed.ValueCount() == 0 {
			return ""
		}

		evd := ed.Value(0)

		if fd.DefaultValue() != "" && fd.DefaultValue() != evd.Name() {
			break
		}

		return evd.Name()

	}

	return fd.DefaultValue()
}
