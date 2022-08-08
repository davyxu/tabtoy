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

// 这里每加一个字段, 需要在InitBuiltinTypes中定义表的解析方式
type TypeDefine struct {
	Kind          TypeUsage `tb_name:"种类"`
	ObjectType    string    `tb_name:"对象类型"`
	Name          string    `tb_name:"标识名"`
	FieldName     string    `tb_name:"字段名"`
	Note          string    `tb_name:"备注"`
	FieldType     string    `tb_name:"字段类型"`
	Value         string    `tb_name:"值" json:",omitempty"`
	ArraySplitter string    `tb_name:"数组切割" json:",omitempty"`
	MakeIndex     bool      `tb_name:"索引" json:",omitempty"`
	Tags          []string  `tb_name:"标记" json:",omitempty"`
	IsBuiltin     bool      `json:",omitempty"`
}

func (self *TypeDefine) ContainTag(tag string) bool {
	for _, libtag := range self.Tags {
		if tag == libtag {
			return true
		}
	}

	return false
}

func (self *TypeDefine) IsArray() bool {
	return self.ArraySplitter != ""
}

// 内建表的列功能
func InitBuiltinTypes(typeTab *TypeTable) {

	for _, tf := range []*TypeDefine{

		// 类型表类型
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "", FieldName: "None", FieldType: "int", Value: "0"},
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "表头", FieldName: "HeaderStruct", FieldType: "int", Value: "1"},
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "枚举", FieldName: "Enum", FieldType: "int", Value: "2"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "种类", FieldName: "Kind", FieldType: "TypeUsage"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "对象类型", FieldName: "ObjectType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "标识名", FieldName: "Name", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "字段名", FieldName: "FieldName", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "字段类型", FieldName: "FieldType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "值", FieldName: "Value", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "数组切割", FieldName: "ArraySplitter", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "索引", FieldName: "MakeIndex", FieldType: "bool"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "标记", FieldName: "Tags", FieldType: "string", ArraySplitter: "|"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "备注", FieldName: "Note", FieldType: "string"},

		// 索引表类型
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "", FieldName: "None", FieldType: "int", Value: "0"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "类型表", FieldName: "Type", FieldType: "int", Value: "1"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "数据表", FieldName: "Data", FieldType: "int", Value: "2"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "键值表", FieldName: "KeyValue", FieldType: "int", Value: "3"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "模式", FieldName: "TableKind", FieldType: "TableKind"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "表类型", FieldName: "TableType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "表文件名", FieldName: "TableFileName", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "标记", FieldName: "Tags", FieldType: "string", ArraySplitter: "|"},

		// KV表类型
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "字段名", FieldName: "FieldName", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "字段类型", FieldName: "FieldType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "标识名", FieldName: "Name", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "值", FieldName: "Value", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "数组切割", FieldName: "ArraySplitter", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "标记", FieldName: "Tags", FieldType: "string", ArraySplitter: "|"},
	} {
		tf.IsBuiltin = true

		typeTab.AddField(tf, nil, 0)
	}
}
