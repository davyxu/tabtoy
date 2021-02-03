package pbdata

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
)

// https://farer.org/2020/04/17/go-protobuf-apiv2-reflect-dynamicpb/
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
