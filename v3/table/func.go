package table

import (
	"github.com/ahmetb/go-linq"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 取类型的默认值
func FetchDefaultValue(tf *TableField) (ret string) {

	linq.From(BuiltinConfig.FieldType).WhereT(func(ft *FieldType) bool {

		return ft.InputFieldName == tf.FieldType
	}).ForEachT(func(ft *FieldType) {

		ret = ft.DefaultValue
	})

	return
}

// 将类型转为对应语言的原始类型
func LanguagePrimitive(tf *TableField, lanType string) string {

	var convertedType string
	linq.From(BuiltinConfig.FieldType).WhereT(func(ft *FieldType) bool {

		return ft.InputFieldName == tf.FieldType
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
		convertedType = tf.FieldType
	}

	return convertedType
}

// 将定义用的类型，转换为不同语言对应的复合类型
func LanguageType(tf *TableField, lanType string) string {

	convertedType := LanguagePrimitive(tf, lanType)

	if tf.IsArray() {
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
	UsefulFunc["LanguageType"] = LanguageType
}
