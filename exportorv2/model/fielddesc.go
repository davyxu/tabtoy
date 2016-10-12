package model

import (
	"fmt"
	"strings"

	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/golang/protobuf/proto"
)

type FieldType int

const (
	FieldType_None   FieldType = 0
	FieldType_Int32  FieldType = 1
	FieldType_Int64  FieldType = 2
	FieldType_UInt32 FieldType = 3
	FieldType_UInt64 FieldType = 4
	FieldType_Float  FieldType = 5
	FieldType_String FieldType = 6
	FieldType_Bool   FieldType = 7
	FieldType_Enum   FieldType = 8
	FieldType_Struct FieldType = 9
)

// 一列的描述
type FieldDescriptor struct {
	Name string

	Type FieldType

	Complex *Descriptor // 复杂类型: 枚举或者结构体

	Order int32 // 在Descriptor中的顺序

	Meta tool.FieldMetaV2 // 扩展字段

	IsRepeated bool

	EnumValue int32 // 枚举值

	Comment string // 注释

	Parent *Descriptor
}

func (self *FieldDescriptor) Tag() int32 {
	return MakeTag(self.Type, self.Order)
}

func MakeTag(t FieldType, order int32) int32 {
	return int32(t)<<16 | order
}

func (self *FieldDescriptor) MetaString() string {

	return proto.MarshalTextString(&self.Meta)
}

func (self *FieldDescriptor) String() string {

	var typestr string
	if self.Complex != nil {
		typestr = fmt.Sprintf("%s|%s", self.Complex.Name, FieldTypeToString(self.Type))
	} else {
		typestr = FieldTypeToString(self.Type)
	}

	var repString string
	if self.IsRepeated {
		repString = "repeated "
	}

	return fmt.Sprintf("name: '%s' %stype: '%s'", self.Name, repString, typestr)
}

func (self *FieldDescriptor) DefaultValue() string {

	if self.Meta.Default != "" {
		return self.Meta.Default
	}

	switch self.Type {
	case FieldType_Int32,
		FieldType_UInt32,
		FieldType_Int64,
		FieldType_UInt64,
		FieldType_Float:
		return "0"
	case FieldType_Bool:
		return "false"
	case FieldType_Enum:

		if self.Complex == nil {
			log.Debugln("build type null while get default value", self.Name)
			return ""
		}

		if len(self.Complex.Fields) == 0 {
			return ""
		}

		return self.Complex.Fields[0].Name

	}

	return ""
}

func (self *FieldDescriptor) ListSpliter() string {

	return self.Meta.ListSpliter
}

func (self *FieldDescriptor) RepeatCheck() bool {

	return self.Meta.RepeatCheck
}

var strByFieldDescriptor = map[FieldType]string{
	FieldType_None:   "none",
	FieldType_Int32:  "int32",
	FieldType_Int64:  "int64",
	FieldType_UInt32: "uint32",
	FieldType_UInt64: "uint64",

	FieldType_Float:  "float",
	FieldType_String: "string",
	FieldType_Bool:   "bool",
	FieldType_Enum:   "enum",
	FieldType_Struct: "struct",
}

var fieldTypeByString = make(map[string]FieldType)

func FieldTypeToString(t FieldType) string {
	if v, ok := strByFieldDescriptor[t]; ok {
		return v
	}

	return "unknown"
}

func ParseFieldType(str string) (t FieldType, ok bool) {
	v, ok := fieldTypeByString[str]
	return v, ok
}

const repeatedKeyword = "repeated"
const repeatedKeywordLen = len(repeatedKeyword)

func (self *FieldDescriptor) ParseType(fileD *FileDescriptor, rawstr string) bool {

	var puretype string

	if strings.HasPrefix(rawstr, repeatedKeyword) {

		puretype = rawstr[repeatedKeywordLen+1:]

		self.IsRepeated = true
	} else {
		puretype = rawstr
	}

	if ft, ok := ParseFieldType(puretype); ok {
		self.Type = ft
		return true
	}

	if desc, ok := fileD.DescriptorByName[puretype]; ok {
		self.Complex = desc

		// 根据内建类型转成字段类型
		switch desc.Kind {
		case DescriptorKind_Struct:
			self.Type = FieldType_Struct
		case DescriptorKind_Enum:
			self.Type = FieldType_Enum
		}

	} else {
		return false
	}

	return true
}

func init() {

	for k, v := range strByFieldDescriptor {
		fieldTypeByString[v] = k
	}

}
