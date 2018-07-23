package model

import (
	"github.com/ahmetb/go-linq"
)

type FieldType struct {
	InputFieldName string `tb_name:"输入字段"`
	GoFieldName    string `tb_name:"Go字段"`
	CSFieldName    string `tb_name:"C#字段"`
	DefaultValue   string `tb_name:"默认值"`
}

// 将表中输入的字段类型转换为各种语言类型

var (
	FieldTypes = []*FieldType{
		{"int32", "int32", "Int32", "0"},
		{"int64", "int64", "Int64", "0"},
		{"int", "int32", "Int32", "0"},
		{"uint64", "uint64", "UInt64", "0"},
		{"uint64", "uint64", "UInt64", "0"},
		{"float", "float32", "float", "0"},
		{"double", "float64", "double", "0"},
		{"float32", "float32", "float", "0"},
		{"float64", "float64", "double", "0"},
		{"bool", "bool", "bool", "FALSE"},
		{"string", "string", "string", ""},
	}
)

// 取类型的默认值
func FetchDefaultValue(tf *TypeDefine) (ret string) {

	linq.From(FieldTypes).WhereT(func(ft *FieldType) bool {

		return ft.InputFieldName == tf.FieldType
	}).ForEachT(func(ft *FieldType) {

		ret = ft.DefaultValue
	})

	return
}

// 将类型转为对应语言的原始类型
func LanguagePrimitive(fieldType string, lanType string) string {

	var convertedType string
	linq.From(FieldTypes).WhereT(func(ft *FieldType) bool {

		return ft.InputFieldName == fieldType
	}).SelectT(func(ft *FieldType) string {

		switch lanType {
		case "cs":
			return ft.CSFieldName
		case "go":
			return ft.GoFieldName
		default:
			panic("unknown lan type: " + lanType)
		}
	}).ForEachT(func(typeName string) {

		convertedType = typeName
	})

	if convertedType == "" {
		convertedType = fieldType
	}

	return convertedType
}

// 原始类型是否存在，例如: int32, int64
func PrimitiveExists(fieldType string) bool {

	return linq.From(FieldTypes).WhereT(func(ft *FieldType) bool {

		return ft.InputFieldName == fieldType
	}).Count() > 0
}
