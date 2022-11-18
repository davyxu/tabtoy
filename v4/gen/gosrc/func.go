package gosrc

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 将定义用的类型，转换为不同语言对应的复合类型

func init() {
	UsefulFunc["GoType"] = func(tf *model.DataField) string {

		convertedType := util.LanguagePrimitive(tf.FieldType, "go")

		if tf.IsArray() {
			return "[]" + convertedType
		}

		return convertedType
	}

	UsefulFunc["GoTabTag"] = func(fieldType *model.DataField) string {

		var sb strings.Builder

		var kv []string

		if len(kv) > 0 {
			sb.WriteString("`")

			for _, s := range kv {
				sb.WriteString(s)
			}

			sb.WriteString("`")
		}

		return sb.String()
	}

	UsefulFunc["JsonTabOmit"] = func() string {
		return "`json:\"-\"`"
	}

}
