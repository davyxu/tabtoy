package model

import (
	"github.com/pkg/errors"
)

type FieldType struct {
	InputFieldName string `tb_name:"输入字段"` // 表中输入的类型
	GoFieldName    string `tb_name:"Go字段"` //  转换为go的类型
	CSFieldName    string `tb_name:"C#字段"`
	JavaFieldName  string `tb_name:"Java字段"`
	PBFieldName    string `tb_name:"Protobuf字段"`
	DefaultValue   string `tb_name:"默认值"`
}

// 将表中输入的字段类型转换为各种语言类型

var (
	FieldTypes = []*FieldType{
		{"int16", "int16", "Int16", "int", "int32", "0"},
		{"int32", "int32", "Int32", "int", "int32", "0"},
		{"int64", "int64", "Int64", "long", "int64", "0"},
		{"int", "int32", "Int32", "int", "int32", "0"},
		{"uint", "uint32", "UInt32", "int", "uint32", "0"},
		{"uint16", "uint16", "UInt16", "int", "uint32", "0"},
		{"uint32", "uint32", "UInt32", "int", "uint32", "0"},
		{"uint64", "uint64", "UInt64", "long", "uint64", "0"},
		{"float", "float32", "float", "float", "float", "0"},
		{"double", "float64", "double", "double", "double", "0"},
		{"float32", "float32", "float", "float", "float", "0"},
		{"float64", "float64", "double", "double", "double", "0"},
		{"bool", "bool", "bool", "boolean", "bool", "FALSE"},
		{"string", "string", "string", "String", "string", ""},
	}

	FieldTypeByType = map[string]*FieldType{}
)

func init() {

	for _, ft := range FieldTypes {
		FieldTypeByType[ft.InputFieldName] = ft
	}
}

// 取类型的默认值
func FetchDefaultValue(fieldType string) (ret string) {

	if ft, ok := FieldTypeByType[fieldType]; ok {
		return ft.DefaultValue
	}

	return
}

// 将类型转为对应语言的原始类型
func LanguagePrimitive(fieldType string, lanType string) string {

	if ft, ok := FieldTypeByType[fieldType]; ok {
		switch lanType {
		case "cs":
			return ft.CSFieldName
		case "go":
			return ft.GoFieldName
		case "java":
			return ft.JavaFieldName
		case "pb":
			return ft.PBFieldName
		default:
			panic("unknown lan type: " + lanType)
		}
	}

	return fieldType
}

// 原始类型是否存在，例如: int32, int64
func PrimitiveExists(fieldType string) bool {

	if _, ok := FieldTypeByType[fieldType]; ok {
		return true
	}

	return false
}

func ParseBool(s string) (bool, error) {
	switch s {
	case "是", "yes", "YES", "1", "true", "TRUE", "True":
		return true, nil
	case "否", "no", "NO", "0", "false", "FALSE", "False":
		return false, nil
	case "":
		return false, nil
	}

	return false, errors.New("invalid bool value")
}
