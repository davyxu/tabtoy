package gosrc

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func KeyValueTypeNames(globals *model.Globals) (ret []string) {
	linq.From(globals.IndexList).WhereT(func(pragma *model.IndexDefine) bool {
		return pragma.Kind == model.TableKind_KeyValue
	}).SelectT(func(pragma *model.IndexDefine) string {

		return pragma.TableType
	}).Distinct().ToSlice(&ret)

	return
}

// 将定义用的类型，转换为不同语言对应的复合类型

func init() {
	UsefulFunc["GoType"] = func(tf *model.TypeDefine) string {

		convertedType := model.LanguagePrimitive(tf.FieldType, "go")

		if tf.IsArray() {
			return "[]" + convertedType
		}

		return convertedType
	}

	UsefulFunc["GoTabTag"] = func(fieldType *model.TypeDefine) string {

		var sb strings.Builder

		var kv []string

		if fieldType.Name != "" {
			kv = append(kv, fmt.Sprintf("tb_name:\"%s\"", fieldType.Name))
		}

		if len(kv) > 0 {
			sb.WriteString("`")

			for _, s := range kv {
				sb.WriteString(s)
			}

			sb.WriteString("`")
		}

		return sb.String()
	}

	UsefulFunc["HasKeyValueTypes"] = func(globals *model.Globals) bool {
		return len(KeyValueTypeNames(globals)) > 0
	}

	UsefulFunc["GetKeyValueTypeNames"] = KeyValueTypeNames

	UsefulFunc["JsonTabOmit"] = func() string {
		return "`json:\"-\"`"
	}

}
