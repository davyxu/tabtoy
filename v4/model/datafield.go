package model

/*
	添加字段和枚举, 需要在model.InitBuiltinTypes函数中添加入口
*/

type FieldUsage int32

const (
	FieldUsage_None   FieldUsage = iota //
	FieldUsage_Struct                   // 结构体
	FieldUsage_Enum                     // 枚举
)

func (self FieldUsage) String() string {

	switch self {
	case FieldUsage_Struct:
		return "表头"
	case FieldUsage_Enum:
		return "枚举"
	default:
		return "未知"
	}
}

type DataField struct {
	Usage         FieldUsage
	ObjectType    string // 枚举类型, 结构体类型
	FieldName     string
	FieldType     string
	Value         string `json:",omitempty"` // 枚举值
	Comment       string `json:",omitempty"` // 枚举值
	ArraySplitter string `json:",omitempty"` // 数组分隔符
	MakeIndex     bool   `json:",omitempty"`
}

func (self *DataField) IsArray() bool {
	return self.ArraySplitter != ""
}
