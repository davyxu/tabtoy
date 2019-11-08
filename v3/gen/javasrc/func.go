package javasrc

import (
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 将定义用的类型，转换为不同语言对应的复合类型

func init() {
	UsefulFunc["JavaType"] = func(tf *model.TypeDefine, requireRef bool) string {

		convertedType := model.LanguagePrimitive(tf.FieldType, "java")

		if requireRef {
			// https://www.geeksforgeeks.org/difference-between-an-integer-and-int-in-java/
			switch convertedType {
			case "int":
				convertedType = "Integer"
			case "short":
				convertedType = "Short"
			case "long":
				convertedType = "Integer"
			case "float":
				convertedType = "Float"
			case "double":
				convertedType = "Double"
			case "boolean":
				convertedType = "Boolean"
			}
		}

		if tf.IsArray() {
			return convertedType + "[]"
		}

		return convertedType
	}

}
