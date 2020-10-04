package pbdata

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"strconv"
)

func buildDynamicType(globals *model.Globals) (protoreflect.FileDescriptor, error) {
	var file descriptorpb.FileDescriptorProto
	file.Syntax = proto.String("proto3")
	file.Name = proto.String(globals.CombineStructName)
	file.Package = proto.String(globals.PackageName)

	for _, tab := range globals.Datas.AllTables() {

		var desc descriptorpb.DescriptorProto
		desc.Name = proto.String(tab.OriginalHeaderType)
		for index, field := range globals.Types.AllFieldByName(tab.OriginalHeaderType) {
			var fd descriptorpb.FieldDescriptorProto
			fd.Name = proto.String(field.FieldName)
			fd.Number = proto.Int32(int32(index + 1))
			fd.JsonName = proto.String(field.FieldName)
			tableType2PbType(globals, field, &fd)
			if field.IsArray() {
				fd.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
			} else {
				fd.Label = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
			}

			desc.Field = append(desc.Field, &fd)
		}

		file.MessageType = append(file.MessageType, &desc)
	}

	for _, enumName := range globals.Types.EnumNames() {

		var ed descriptorpb.EnumDescriptorProto
		ed.Name = proto.String(enumName)

		for _, field := range globals.Types.AllFieldByName(enumName) {
			var vd descriptorpb.EnumValueDescriptorProto
			vd.Name = proto.String(field.FieldName)
			v, _ := strconv.Atoi(field.Value)
			vd.Number = proto.Int32(int32(v))
			ed.Value = append(ed.Value, &vd)
		}
		file.EnumType = append(file.EnumType, &ed)
	}

	var combine descriptorpb.DescriptorProto
	combine.Name = proto.String(globals.CombineStructName)
	for index, md := range file.MessageType {
		var fd descriptorpb.FieldDescriptorProto
		fd.Name = proto.String(md.GetName())
		fd.Number = proto.Int32(int32(index + 1))
		fd.JsonName = proto.String(md.GetName())
		fd.Type = descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum()
		fd.TypeName = proto.String(globals.PackageName + "." + md.GetName())
		fd.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
		combine.Field = append(combine.Field, &fd)
	}

	file.MessageType = append(file.MessageType, &combine)

	return protodesc.NewFile(&file, nil)
}

// https://farer.org/2020/04/17/go-protobuf-apiv2-reflect-dynamicpb/
func Generate(globals *model.Globals) (data []byte, err error) {

	pbFile, err := buildDynamicType(globals)
	if err != nil {
		return nil, err
	}

	combineD := pbFile.Messages().ByName(protoreflect.Name(globals.CombineStructName))

	combineRoot := dynamicpb.NewMessage(combineD)

	// 所有的表
	for _, tab := range globals.Datas.AllTables() {

		md := pbFile.Messages().ByName(protoreflect.Name(tab.OriginalHeaderType))

		combineField := combineD.Fields().ByName(protoreflect.Name(tab.OriginalHeaderType))
		list := combineRoot.NewField(combineField).List()

		// 每个表的所有列
		headers := globals.Types.AllFieldByName(tab.OriginalHeaderType)

		// 遍历每一行
		for row := 1; row < len(tab.Rows); row++ {

			msg := dynamicpb.NewMessage(md)

			for col, field := range headers {

				fd := md.Fields().ByName(protoreflect.Name(field.FieldName))

				// 在单元格找到值
				valueCell := tab.GetCell(row, col)
				if valueCell == nil {
					continue
				}

				if field.IsArray() {
					list := msg.NewField(fd).List()
					tableValue2PbValueList(globals, valueCell, field, list)
					msg.Set(fd, protoreflect.ValueOfList(list))
				} else {
					pbValue := tableValue2PbValue(globals, valueCell, field)
					msg.Set(fd, pbValue)
				}

			}

			list.Append(protoreflect.ValueOf(msg))
		}

		combineRoot.Set(combineField, protoreflect.ValueOfList(list))
	}

	return proto.Marshal(combineRoot)
}

func tableValue2PbValueList(globals *model.Globals, cell *model.Cell, valueType *model.TypeDefine, list protoreflect.List) {

	pbType := model.LanguagePrimitive(valueType.FieldType, "pb")

	switch pbType {
	case "int32":
		for _, str := range cell.ValueList {
			v, _ := strconv.ParseInt(str, 10, 32)
			list.Append(protoreflect.ValueOfInt32(int32(v)))
		}

	case "int64":
		for _, str := range cell.ValueList {
			v, _ := strconv.ParseInt(str, 10, 64)
			list.Append(protoreflect.ValueOfInt64(v))
		}
	case "float":

		for _, str := range cell.ValueList {
			v, _ := strconv.ParseFloat(str, 32)
			list.Append(protoreflect.ValueOfFloat32(float32(v)))
		}
	case "double":
		for _, str := range cell.ValueList {
			v, _ := strconv.ParseFloat(str, 32)
			list.Append(protoreflect.ValueOfFloat32(float32(v)))
		}
	case "bool":

		for _, str := range cell.ValueList {

			v, _ := model.ParseBool(str)
			list.Append(protoreflect.ValueOfBool(v))
		}

	case "string":
		for _, str := range cell.ValueList {
			list.Append(protoreflect.ValueOfString(str))
		}
	default:

		if globals.Types.IsEnumKind(pbType) {
			for _, str := range cell.ValueList {

				enumValue := globals.Types.ResolveEnumValue(pbType, str)

				v, _ := strconv.ParseInt(enumValue, 10, 32)

				list.Append(protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)))
			}

		} else {
			panic("unknown pb type: " + pbType)
		}
	}
}

func tableValue2PbValue(globals *model.Globals, cell *model.Cell, valueType *model.TypeDefine) protoreflect.Value {

	pbType := model.LanguagePrimitive(valueType.FieldType, "pb")

	switch pbType {
	case "int32":
		if cell.Value == "" {
			return protoreflect.ValueOfInt32(0)
		}
		v, _ := strconv.ParseInt(cell.Value, 10, 32)
		return protoreflect.ValueOfInt32(int32(v))
	case "uint32":
		if cell.Value == "" {
			return protoreflect.ValueOfUint32(0)
		}
		v, _ := strconv.ParseUint(cell.Value, 10, 32)
		return protoreflect.ValueOfUint32(uint32(v))
	case "int64":
		if cell.Value == "" {
			return protoreflect.ValueOfInt64(0)
		}
		v, _ := strconv.ParseInt(cell.Value, 10, 64)
		return protoreflect.ValueOfInt64(v)
	case "uint64":
		if cell.Value == "" {
			return protoreflect.ValueOfUint64(0)
		}
		v, _ := strconv.ParseUint(cell.Value, 10, 64)
		return protoreflect.ValueOfUint64(v)
	case "float":
		if cell.Value == "" {
			return protoreflect.ValueOfFloat32(0)
		}
		v, _ := strconv.ParseFloat(cell.Value, 32)
		return protoreflect.ValueOfFloat32(float32(v))
	case "double":
		if cell.Value == "" {
			return protoreflect.ValueOfFloat64(0)
		}
		v, _ := strconv.ParseFloat(cell.Value, 64)
		return protoreflect.ValueOfFloat64(v)
	case "bool":
		if cell.Value == "" {
			return protoreflect.ValueOfBool(false)
		}

		v, _ := model.ParseBool(cell.Value)
		return protoreflect.ValueOfBool(v)
	case "string":
		return protoreflect.ValueOfString(cell.Value)
	default:

		if globals.Types.IsEnumKind(pbType) {

			if cell.Value == "" {
				return protoreflect.ValueOfInt32(0)
			}
			enumValue := globals.Types.ResolveEnumValue(pbType, cell.Value)

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
