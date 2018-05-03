package table

import (
	"encoding/json"
)

var (
	CoreConfig Config // 内嵌数据

	CoreSymbols = []*TableField{
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
)

func init() {
	err := json.Unmarshal([]byte(coreConfig), &CoreConfig)
	if err != nil {
		panic(err)
	}

	for _, symbol := range CoreSymbols {
		symbol.IsBuiltin = true
	}
}
