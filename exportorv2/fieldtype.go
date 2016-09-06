package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/proto/tool"
)

type FieldType int

const (
	FieldType_None FieldType = iota
	FieldType_Int32
	FieldType_Int64
	FieldType_UInt32
	FieldType_UInt64
	FieldType_Float
	FieldType_String
	FieldType_Bool
	FieldType_Enum
	FieldType_Struct
)

// 一列的描述
type FieldDefine struct {
	Name string

	Type FieldType

	BuildInType *BuildInType // 复杂类型: 枚举或者结构体

	IsRepeated bool

	Meta tool.FieldMetaV2 // 扩展字段
}

func (self *FieldDefine) DefaultValue() string {

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

		if self.BuildInType == nil {
			log.Debugln("build type null while get default value", self.Name)
			return ""
		}

		if len(self.BuildInType.Fields) == 0 {
			return ""
		}

		return self.BuildInType.Fields[0].Name

	}

	return ""
}

func (self *FieldDefine) ListSpliter() string {

	return self.Meta.ListSpliter
}

func (self *FieldDefine) RepeatCheck() bool {

	return self.Meta.RepeatCheck
}

var strByFieldType = map[FieldType]string{
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

func FieldTypeToString(t FieldType) string {
	if v, ok := strByFieldType[t]; ok {
		return v
	}

	return "unknown"
}

const repeatedKeyword = "repeated"
const repeatedKeywordLen = len(repeatedKeyword)

func (self *FieldDefine) ParseType(tts *BuildInTypeSet, rawstr string) bool {

	if strings.HasPrefix(rawstr, repeatedKeyword) {

		rawstr = rawstr[repeatedKeywordLen+1:]

		self.IsRepeated = true
	}

	for ft, s := range strByFieldType {
		if rawstr == s {
			self.Type = ft
			return true
		}
	}

	if buildinType, ok := tts.TypeByName[rawstr]; ok {
		self.BuildInType = buildinType

		// 根据内建类型转成字段类型
		switch buildinType.Kind {
		case BuildInTypeKind_Struct:
			self.Type = FieldType_Enum
		case BuildInTypeKind_Enum:
			self.Type = FieldType_Struct
		}

	} else {
		return false
	}

	return true
}
