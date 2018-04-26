package table

import (
	"github.com/ahmetb/go-linq"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 将定义用的类型，转换为不同语言对应的目标类型
func ConverToLanType(tf *TypeField, lanType string) string {

	var convertedType string
	linq.From(config.FieldType).WhereT(func(ft FieldType) bool {

		return ft.InputFieldName == tf.FieldType
	}).SelectT(func(ft FieldType) string {

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
		convertedType = tf.FieldType
	}

	if tf.IsArray {
		switch lanType {
		case "cs":
			return convertedType + "[]"
		case "go":
			return "[]" + convertedType
		default:
			panic("unknown lan type: " + lanType)
		}
	}

	return convertedType
}

func init() {
	UsefulFunc["ConverToLanType"] = ConverToLanType
}
