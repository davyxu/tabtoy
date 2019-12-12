package javasrc

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 将定义用的类型，转换为不同语言对应的复合类型

func wrapSingleValue(globals *model.Globals, valueType *model.TypeDefine, value string) string {
	switch {
	case valueType.FieldType == "string": // 字符串
		return util.StringEscape(value)
	case valueType.FieldType == "float32":
		return value
	case globals.Types.IsEnumKind(valueType.FieldType): // 枚举
		t := globals.Types.ResolveEnum(valueType.FieldType, value)
		if t != nil {
			return t.Define.ObjectType + "." + t.Define.FieldName
		}

		return ""
	case valueType.FieldType == "bool":

		switch value {
		case "是", "yes", "YES", "1", "true", "TRUE", "True":
			return "true"
		case "否", "no", "NO", "0", "false", "FALSE", "False":
			return "false"
		}

		return "false"
	}

	if value == "" {
		return model.FetchDefaultValue(valueType.FieldType)
	}

	return value
}

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

	UsefulFunc["JavaDefaultValue"] = func(globals *model.Globals, tf *model.TypeDefine) string {

		convertedType := model.LanguagePrimitive(tf.FieldType, "java")

		if tf.IsArray() {
			return fmt.Sprintf("new %s[]{}", convertedType)
		} else {
			return wrapSingleValue(globals, tf, "")
		}

		return convertedType
	}

}
