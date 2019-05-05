package model

/*
	添加字段和枚举, 需要在model.InitBuiltinTypes函数中添加入口
*/

type TypeUsage int32

const (
	TypeUsage_None         TypeUsage = iota //
	TypeUsage_HeaderStruct                  // 表头
	TypeUsage_Enum                          // 枚举
)

func (self TypeUsage) String() string {

	switch self {
	case TypeUsage_HeaderStruct:
		return "表头"
	case TypeUsage_Enum:
		return "枚举"
	default:
		return "未知"
	}
}

type TypeDefine struct {
	Kind          TypeUsage `tb_name:"种类"`
	ObjectType    string    `tb_name:"对象类型"`
	Name          string    `tb_name:"标识名"`
	FieldName     string    `tb_name:"字段名"`
	FieldType     string    `tb_name:"字段类型"`
	Value         string    `tb_name:"值" json:",omitempty"`
	ArraySplitter string    `tb_name:"数组切割" json:",omitempty"`
	MakeIndex     bool      `tb_name:"索引" json:",omitempty"`
	IsBuiltin     bool      `json:",omitempty"`
}

func (self *TypeDefine) IsArray() bool {
	return self.ArraySplitter != ""
}
