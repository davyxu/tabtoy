package cssrc

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/gen/binpak"
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

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
	UsefulFunc["CSType"] = func(tf *model.TypeDefine) string {

		convertedType := model.LanguagePrimitive(tf.FieldType, "cs")

		if tf.IsArray() {
			return fmt.Sprintf("List<%s>", convertedType)
		}

		return convertedType
	}

	UsefulFunc["CSTag"] = func(globals *model.Globals, fieldIndex int, tf *model.TypeDefine) string {

		tag := binpak.MakeTag(globals, tf, fieldIndex)

		return fmt.Sprintf("0x%x", tag)
	}

	UsefulFunc["CSReader"] = func(globals *model.Globals, tf *model.TypeDefine) (ret string) {

		convertedType := model.LanguagePrimitive(tf.FieldType, "cs")

		switch {
		case convertedType == "float":
			ret = "Float"
		case convertedType == "string":
			ret = "String"
		case convertedType == "bool":
			ret = "Bool"
		case globals.Types.IsEnumKind(tf.FieldType):
			ret = "Enum"
		default:
			ret = convertedType
		}

		return
	}

	UsefulFunc["CSDefaultValue"] = func(globals *model.Globals, tf *model.TypeDefine) string {

		convertedType := model.LanguagePrimitive(tf.FieldType, "cs")

		if tf.IsArray() {
			return fmt.Sprintf("new List<%s>()", convertedType)
		} else {
			return wrapSingleValue(globals, tf, "")
		}

		return convertedType
	}

}
