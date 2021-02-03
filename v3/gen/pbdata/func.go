package pbdata

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
)

func tableValue2PbValueList(globals *model.Globals, cell *model.Cell, valueType *model.TypeDefine, list protoreflect.List) {

	pbType := model.LanguagePrimitive(valueType.FieldType, "pb")

	if globals.Types.IsEnumKind(pbType) {

		for _, str := range cell.ValueList {
			enumValue := globals.Types.ResolveEnumValue(pbType, str)
			v, _ := strconv.ParseInt(enumValue, 10, 32)
			list.Append(protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)))
		}

	} else {
		for _, str := range cell.ValueList {

			v := tableValue2PbValue(globals, str, valueType)
			list.Append(v)
		}
	}

}

func tableValue2PbValue(globals *model.Globals, cellValue string, valueType *model.TypeDefine) protoreflect.Value {

	pbType := model.LanguagePrimitive(valueType.FieldType, "pb")

	switch pbType {
	case "int32":
		if cellValue == "" {
			return protoreflect.ValueOfInt32(0)
		}
		v, _ := strconv.ParseInt(cellValue, 10, 32)
		return protoreflect.ValueOfInt32(int32(v))
	case "uint32":
		if cellValue == "" {
			return protoreflect.ValueOfUint32(0)
		}
		v, _ := strconv.ParseUint(cellValue, 10, 32)
		return protoreflect.ValueOfUint32(uint32(v))
	case "int64":
		if cellValue == "" {
			return protoreflect.ValueOfInt64(0)
		}
		v, _ := strconv.ParseInt(cellValue, 10, 64)
		return protoreflect.ValueOfInt64(v)
	case "uint64":
		if cellValue == "" {
			return protoreflect.ValueOfUint64(0)
		}
		v, _ := strconv.ParseUint(cellValue, 10, 64)
		return protoreflect.ValueOfUint64(v)
	case "float":
		if cellValue == "" {
			return protoreflect.ValueOfFloat32(0)
		}
		v, _ := strconv.ParseFloat(cellValue, 32)
		return protoreflect.ValueOfFloat32(float32(v))
	case "double":
		if cellValue == "" {
			return protoreflect.ValueOfFloat64(0)
		}
		v, _ := strconv.ParseFloat(cellValue, 64)
		return protoreflect.ValueOfFloat64(v)
	case "bool":
		if cellValue == "" {
			return protoreflect.ValueOfBool(false)
		}

		v, _ := model.ParseBool(cellValue)
		return protoreflect.ValueOfBool(v)
	case "string":
		return protoreflect.ValueOfString(cellValue)
	default:

		if globals.Types.IsEnumKind(pbType) {

			if cellValue == "" {
				return protoreflect.ValueOfEnum(protoreflect.EnumNumber(0))
			}
			enumValue := globals.Types.ResolveEnumValue(pbType, cellValue)

			v, _ := strconv.ParseInt(enumValue, 10, 32)

			return protoreflect.ValueOfEnum(protoreflect.EnumNumber(v))
		} else {
			panic("unknown pb type: " + pbType)
		}
	}
}

func tableType2PbType(globals *model.Globals, def *model.TypeDefine, pbDesc *descriptorpb.FieldDescriptorProto) {
	pbType := model.LanguagePrimitive(def.FieldType, "pb")

	switch pbType {
	case "int32":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum()
	case "uint32":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_UINT32.Enum()
	case "int64":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum()
	case "uint64":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum()
	case "float":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_FLOAT.Enum()
	case "double":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum()
	case "bool":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum()
	case "string":
		pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()
	default:

		if globals.Types.IsEnumKind(pbType) {
			pbDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum()
			pbDesc.TypeName = proto.String(globals.PackageName + "." + pbType)
		} else {
			panic("unknown pb type: " + pbType)
		}
	}
}
