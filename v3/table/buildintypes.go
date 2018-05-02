package table

type TableKind int32

const (
	TableKind_None         = 0 //
	TableKind_HeaderStruct = 1 // 表头
	TableKind_Enum         = 2 // 枚举
)

type TableField struct {
	Kind          TableKind `tb_name:"种类"`
	ObjectType    string    `tb_name:"对象类型"`
	Name          string    `tb_name:"标识名"`
	FieldName     string    `tb_name:"字段名"`
	FieldType     string    `tb_name:"字段类型"`
	Value         string    `tb_name:"值"`
	ArraySplitter string    `tb_name:"数组切割"`
	IsBuiltin     bool
}

func (self *TableField) IsArray() bool {
	return self.ArraySplitter != ""
}

var CoreSymbols = []*TableField{
	{Kind: TableKind_Enum, ObjectType: "TableKind", Name: "", FieldName: "None", FieldType: "int", Value: "0"},
	{Kind: TableKind_Enum, ObjectType: "TableKind", Name: "表头", FieldName: "HeaderStruct", FieldType: "int", Value: "1"},
	{Kind: TableKind_Enum, ObjectType: "TableKind", Name: "枚举", FieldName: "Enum", FieldType: "int", Value: "2"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "种类", FieldName: "Kind", FieldType: "TableKind"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "对象类型", FieldName: "ObjectType", FieldType: "string"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "标识名", FieldName: "Name", FieldType: "string"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "字段名", FieldName: "FieldName", FieldType: "string"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "字段类型", FieldName: "FieldType", FieldType: "string"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "值", FieldName: "Value", FieldType: "string"},
	{Kind: TableKind_HeaderStruct, ObjectType: "TableField", Name: "数组切割", FieldName: "ArraySplitter", FieldType: "string"},
}

func init() {
	for _, symbol := range CoreSymbols {
		symbol.IsBuiltin = true
	}
}
